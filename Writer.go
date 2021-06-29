/*
Copyright (c) 2021 Seculeet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type CommandFailedResult struct {
	NameOutput        string `json:"Name"`
	CommandOutput     string `json:"Command"`
	CommandSuccessful bool   `json:"Command was executed"`
	ErrorMessage      string `json:"Error-Message"`
}

type AuditSuccessfulResult struct {
	NameOutput        string `json:"Name"`
	CommandOutput     string `json:"Command"`
	CommandSuccessful bool   `json:"Command was executed"`
	AuditSuccessful   bool   `json:"Output is as expected"`
}

type AuditFailedResult struct {
	NameOutput        string `json:"Name"`
	CommandOutput     string `json:"Command"`
	CommandSuccessful bool   `json:"Command was executed"`
	AuditSuccessful   bool   `json:"Output is as expected"`
	Expected          string `json:"Expected Value"`
	Out               string `json:"Actual Value"`
	Operator          string `json:"Operator"`
}

var FirstAuditEntry bool
var FirstErrorEntry bool
var FirstResultEntry bool

func WriteResultJSON(audit BigAudit, isCommandSuccessful bool, isAuditSuccessful bool, output string, err string, operator string) {

	var fileText string
	var auditResult interface{}
	var auditAsByteArr []byte

	if !isCommandSuccessful {
		auditResult = CommandFailedResult{
			NameOutput:        audit.Name,
			CommandOutput:     audit.Command,
			CommandSuccessful: isCommandSuccessful,
			ErrorMessage:      err,
		}
	} else {
		if isAuditSuccessful {
			auditResult = AuditSuccessfulResult{
				NameOutput:        audit.Name,
				CommandOutput:     audit.Command,
				CommandSuccessful: isCommandSuccessful,
				AuditSuccessful:   isAuditSuccessful,
			}
		} else {
			auditResult = AuditFailedResult{
				NameOutput:        audit.Name,
				CommandOutput:     audit.Command,
				CommandSuccessful: isCommandSuccessful,
				AuditSuccessful:   isAuditSuccessful,
				Expected:          audit.Expected,
				Out:               output,
				Operator:          operator,
			}
		}
	}
	auditAsByteArr, _ = json.MarshalIndent(auditResult, "\t\t", "\t")
	auditAsByteArr, _ = UnescapeUnicodeCharactersInJSON(auditAsByteArr)

	if !FirstResultEntry {
		os.Remove("./output/result.json")
		FirstResultEntry = true
	}

	if !checkPathExists("./output/result.json") {
		_, createFileErr := os.Create("./output/result.json")
		if createFileErr != nil {
			WriteErrorLog("cannot create file ./output/result.json: "+createFileErr.Error(), "")
			if debugModeEnabled {
				WriteDebugLog("cannot create file ./output/result.json: "+createFileErr.Error(), "ERROR")
			}
			os.Exit(0)
		}
		if debugModeEnabled {
			WriteDebugLog("created file ./output/result.json: ", "INFO")
		}
		fileText = "{\n\t" + `"` + ConfigName + `":` + " [\n\t\t"
		fileText += string(auditAsByteArr)
		fileText += "\n\t]\n}"
	} else {
		fileText = getResultJSONContent()
		fileText = strings.Replace(fileText, "\n\t]\n}", ",\n\t\t", -1)
		fileText += string(auditAsByteArr)
		fileText += "\n\t]\n}"
	}
	fileWriter(fileText, "result.json", false)
}

//let auditName empty if u don't want to log an audit
func WriteErrorLog(err string, auditName string) {
	logType := "ERROR"
	logTextErr := getLogText(err, logType)
	if len(auditName) > 0 {
		logTextErr = getLogText(auditName+" "+err, logType)
	}
	fileWriter(logTextErr, "error.log", true)
	fileWriter(logTextErr, "audit.log", true)
}

func WriteAuditFailedLog(audits []SmallAudit, failPosition int) {
	for i, audit := range audits {
		var auditText string
		if len(audit.Name) > 0 && i == 0 {
			auditText += "Name: " + audit.Name
			WriteLog(auditText, "FAIL")
		}
		if len(audit.Command) > 0 {
			auditText = " " + audit.Name + ": " + audit.Command
		}
		if len(audit.Arguments) > 0 {
			for i, arg := range audit.Arguments {
				if len(arg) > 0 {
					auditText += " " + audit.Arguments[i]
				}
			}
		}

		space := "\t\t\t\t"
		if i == failPosition {
			//failed commands
			auditText = space + " [FAIL] :" + auditText + "\n"
			fileWriter(auditText, "audit.log", true)
		} else if i < failPosition {
			//success executed commands
			auditText = space + " [INFO] :" + auditText + "\n"
			fileWriter(auditText, "audit.log", true)
		} else {
			//not executed commands
			auditText = space + " [WARN] :" + auditText + "\n"
			fileWriter(auditText, "audit.log", true)
		}
	}
}

func WriteLog(text string, logType string) {
	text = getLogText(text, logType)
	fileWriter(text, "audit.log", true)
}

func WriteDebugLog(text string, logType string) {
	text = getLogText(text, logType)
	fileWriter(text, "debug.log", true)
}

func isValidLogType(logType string) string {
	logType = strings.ToUpper(logType)
	switch logType {
	case "INFO":
		return logType
	case "DEBUG":
		return logType
	case "WARN":
		return logType
	case "FAIL":
		return logType
	case "ERROR":
		return logType
	default:
		err := logType + " is not a valid Log-Type, please use a correct Log-Type"
		WriteErrorLog(err, "")
		os.Exit(0)
	}
	return ""
}

func fileWriter(output string, fileName string, append bool) error {
	path := "./output/"
	createFolderIfNotExist(path, "output")

	path += fileName
	openFileFlag := os.O_CREATE + os.O_WRONLY

	switch fileName {
	case "audit.log":
		if !FirstAuditEntry {
			os.Remove(path)
			FirstAuditEntry = true
			if checkPathExists("./output") {
				WriteLog("Create output folder", "INFO")
			}
		}
	case "error.log":
		if !FirstErrorEntry {
			os.Remove(path)
			FirstErrorEntry = true
		}
	}

	if append {
		openFileFlag += os.O_APPEND
	} else {
		if checkPathExists(path) {
			deleteFile(path)
		}
	}

	file, err := os.OpenFile(path, openFileFlag, 0666)
	if err != nil {
		errTxt := "Can't create " + path + " file. " + err.Error()
		PrintToConsole(errTxt)
		return errors.New(errTxt)
	}
	_, fileWriterErr := file.Write([]byte(output))
	if fileWriterErr != nil {
		errTxt := "golang fileWriter not working: " + fileWriterErr.Error()
		PrintToConsole(errTxt)
		return errors.New(errTxt)
	}
	defer file.Close()

	return nil
}

func CurrentDateTime() string {
	currentTime := time.Now().Format("02.01.2006 15:04")
	return currentTime
}

func getLogText(logText string, logType string) string {
	logType = isValidLogType(logType)
	if logType == "ERROR" {
		return CurrentDateTime() + " [" + logType + "]: " + logText + "\n"
	}
	return CurrentDateTime() + " [" + logType + "] : " + logText + "\n"
}

func checkPathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func createFolderIfNotExist(path string, name string) {
	if !checkPathExists(path) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			err := "Can't create " + name + " folder. " + err.Error()
			WriteErrorLog(err, "")
			os.Exit(0)
		} else if name == "artefacts" {
			WriteLog("Create artefacts folder", "INFO")
		}
	}
}

func UnescapeUnicodeCharactersInJSON(jsonRaw json.RawMessage) (json.RawMessage, error) {
	str := strconv.Quote(string(jsonRaw))
	str = strings.ReplaceAll(str, `\\u`, `\u`)
	str = strings.ReplaceAll(str, `\\\u`, `\\\\u`)
	str, err := strconv.Unquote(str)
	if err != nil {
		WriteErrorLog(err.Error(), "")
		return nil, err
	}
	return []byte(str), nil
}

func WriteCommandSuccessLog(audit BigAudit) {
	space := "\n\t\t\t\t\t\t  "
	auditText := "Name: " + audit.Name + space + audit.Name + ": " + audit.Command + " executed"
	WriteLog(auditText, "INFO")
}

func WriteCommandFailedLog(audit BigAudit, err error) {
	space := "\n\t\t\t\t\t\t  "
	auditText := "Name: " + audit.Name + space + audit.Name + ": " + audit.Command + " failed"
	if err != nil {
		WriteLog(auditText, "FAIL")
	}
}

func deleteFile(path string) {
	if checkPathExists(path) {
		err := os.Remove(path)
		for err != nil {
			err = os.Remove(path)
		}
	}
}

func deleteOutput() {
	path := "./output"
	if checkPathExists(path) {
		err := os.RemoveAll(path)
		for err != nil {
			err = os.RemoveAll(path)
		}
	}
}

func replaceAllWhitespace(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
