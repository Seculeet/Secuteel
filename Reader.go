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
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type System struct {
	System Systemdetails `json:"system"`
}

type Systemdetails struct {
	SystemName      string `json:"systemname"`
	Version         string `json:"version"`
	Shell           string `json:"shell"`
	Argument        string `json:"argument"`
	RootPermissions bool   `json:"root"`
}

type BigAudit struct {
	Name             string `json:"name"`
	Command          string `json:"command"`
	DontSaveArtefact bool   `json:"dontSaveArtefact"`
	BlackenContent   string `json:"blackenContent"`
	TypeExpected     string `json:"typeExpected"`
	Expected         string `json:"expected"`
	Desc             string `json:"description"`
}

type BigAudits struct {
	BigAuditArray []BigAudit `json:"commands"`
}

var commands BigAudits
var system System
var ConfigName string

//	Reads config file and adds values in structs
//	uses path that is being made in readInput() method
func ReadConfig() {
	var jsonFile *os.File

	path, _, _ := readInput()
	ConfigName = path

	jsonFile, err := os.Open(path)
	if err != nil {
		WriteErrorLog(err.Error(), "")
		if debugModeEnabled {
			WriteDebugLog(err.Error(), "ERROR")
		}
		os.Exit(0)
	} else {
		WriteLog("input-File was opened", "INFO")
		if debugModeEnabled {
			WriteDebugLog("input-File was opened", "INFO")
		}
	}
	defer jsonFile.Close()
	byteValue, readAllErr := ioutil.ReadAll(jsonFile)

	if readAllErr != nil {
		WriteErrorLog(readAllErr.Error(), "")
		if debugModeEnabled {
			WriteDebugLog(readAllErr.Error(), "ERROR")
		}
	}

	commands = BigAudits{}
	system = System{}

	if checkJsonFormat(byteValue, &commands) && checkJsonFormat(byteValue, &system) {
		WriteLog("JSON format correct", "INFO")
		if debugModeEnabled {
			WriteDebugLog("JSON format correct", "INFO")
		}
	} else {
		errTxt := "JSON format incorrect"
		WriteErrorLog(errTxt, "")
		if debugModeEnabled {
			WriteDebugLog("JSON format incorrect", "ERROR")
		}
		os.Exit(0)
	}

	err = checkSystemIsValid()
	if err != nil {
		fmt.Println(err)
		WriteErrorLog(err.Error(), "")
		os.Exit(0)
	}
}

func GetBigAudits() []BigAudit {
	return commands.BigAuditArray
}

func checkSystemIsValid() error {
	if system == (System{}) {
		if debugModeEnabled {
			WriteDebugLog("you have to specify the system in your config file", "ERROR")
		}
		return errors.New("you have to specify the system in your config file")
	}
	if debugModeEnabled {
		WriteDebugLog("system details read successfully", "INFO")
	}

	sys := GetSystem()
	if len(sys.SystemName) == 0 {
		if debugModeEnabled {
			WriteDebugLog("you have to specify the systemName in your config file", "ERROR")
		}
		return errors.New("you have to specify the systemName in your config file")
	}

	if debugModeEnabled {
		WriteDebugLog("systemName in your config is correct", "INFO")
	}

	if len(sys.Version) == 0 {
		if debugModeEnabled {
			WriteDebugLog("you have to specify the version in your config file", "ERROR")
		}
		return errors.New("you have to specify the version in your config file")
	}

	if debugModeEnabled {
		WriteDebugLog("version in your config is correct", "INFO")
	}
	if debugModeEnabled {
		WriteDebugLog("detected operating system: "+runtime.GOOS, "INFO")
	}

	if strings.EqualFold(runtime.GOOS, "windows") {
		if len(sys.Shell) == 0 {
			system.System.Shell = "powershell"
			if debugModeEnabled {
				WriteDebugLog("using shell: powershell", "INFO")
			}
		}

		if len(sys.Argument) == 0 {
			system.System.Argument = "/C"
			if debugModeEnabled {
				WriteDebugLog("using systemArgument: /C", "INFO")
			}
		}
	} else {
		if len(sys.Shell) == 0 {
			system.System.Shell = "bash"
			if debugModeEnabled {
				WriteDebugLog("using shell: bash", "INFO")
			}
		}

		if len(sys.Argument) == 0 {
			system.System.Argument = "-c"
			if debugModeEnabled {
				WriteDebugLog("using systemArgument: -c", "INFO")
			}
		}
	}
	return nil
}

func checkBigAuditsAreValid() error {
	bigAudits := GetBigAudits()
	if len(bigAudits) == 0 {
		if debugModeEnabled {
			WriteDebugLog("you have to specify at least one command in your config file", "ERROR")
		}
		return errors.New("you have to specify at least one command in your config file")
	}
	if debugModeEnabled {
		WriteDebugLog("audits detected: "+fmt.Sprint(len(bigAudits)), "INFO")
	}

	auditNameMap := make(map[string]int)
	for i, audit := range bigAudits {
		i++
		err := "Issue at the " + getOrdinalNum(i) + " audit."

		if len(audit.Name) == 0 {
			return errors.New(err + " You have to specify the name in your config file")
		}

		err1 := checkValidFilename(audit.Name)
		if err1 != nil {
			return errors.New(err + " The character \"" + err1.Error() + "\" in name is not allowed")
		}

		if len(audit.Command) == 0 {
			return errors.New(err + " You have to specify the command in your config file")
		}

		auditName := strings.ToLower(audit.Name)
		auditPosition := auditNameMap[auditName]
		if auditPosition > 0 {
			return errors.New(err + " The audit name \"" + audit.Name + "\" was already used in the " + getOrdinalNum(auditPosition) + " audit")
		}
		auditNameMap[auditName] = i
	}
	if debugModeEnabled {
		WriteDebugLog("all audits have the required structure", "INFO")
	}
	return nil
}

func checkValidFilename(filename string) error {
	invalidSymbolsWindows := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	invalidCharsLinux := []string{"/"}

	var invalidSymbols []string
	if strings.EqualFold(runtime.GOOS, "windows") {
		invalidSymbols = invalidSymbolsWindows
	} else {
		invalidSymbols = invalidCharsLinux
	}

	for _, invalidSymbol := range invalidSymbols {
		if strings.Contains(filename, invalidSymbol) {
			return errors.New(invalidSymbol)
		}
	}
	return nil
}

func getOrdinalNum(num int) string {
	if num <= 0 {
		err := "The entered number must be greater than 0"
		WriteErrorLog(err, "")
		os.Exit(0)
	}
	var numTxt string

	temp := num % 10

	switch temp {
	case 1:
		numTxt = "st"
	case 2:
		numTxt = "nd"
	case 3:
		numTxt = "rd"
	default:
		numTxt = "th"
	}
	return strconv.Itoa(num) + numTxt
}

func GetSystem() Systemdetails {
	return system.System
}

/*
	Reads the users input in StdIn
	Checks whether or not the path can be found
	Repeats until a correct input is made
*/
func readInput() (string, string, []string) {
	var inPath string     //path to input file
	var outPath string    //path to output file
	var commands []string // commands to add

	var input string

	in := flags.input

	if in != "" {
		input = convertInputToJson(in)
		in = ""
	} else {
		err := errors.New("you need to specify an input file")
		WriteErrorLog(err.Error(), "")
		if debugModeEnabled {
			WriteDebugLog(err.Error(), "ERROR")
		}
		os.Exit(0)

	}

	inPath = convertInputToPath(input)

	WriteLog("input-Path \""+inPath+"\" is correct", "INFO")
	if debugModeEnabled {
		WriteDebugLog("input-Path \""+inPath+"\" is correct", "INFO")
	}
	return inPath, outPath, commands
}

//	Checks Json Format you dont say so!?
//	true if format is viable
func checkJsonFormat(jsonData []byte, toStruct interface{}) bool {
	err := json.Unmarshal(jsonData, toStruct)
	return err == nil
}

//	Checks if input already has ".json" as ending with inputHasJsonEnding() method
func convertInputToJson(input string) string {
	if !inputHasJsonEnding(input) {
		input += ".json"
	}
	return input
}

//	Calls inputHasEnding() method with ".json" as expected ending
func inputHasJsonEnding(input string) bool {
	jsonEnding := ".json"
	return inputHasEnding(input, jsonEnding)
}

//	Checks if the input already has the expected ending
func inputHasEnding(input string, expectedEnding string) bool {
	if len(input) >= len(expectedEnding) {
		inputEnding := input[len(input)-len(expectedEnding):]
		if strings.EqualFold(inputEnding, expectedEnding) {
			return true
		}
	}
	return false
}

//	Checks whether or not input is a path
//	if yes, returns the path
//	if not, input is the name of the config file
//	generates a relative path to the file and returns it
func convertInputToPath(input string) string {
	input = convertBackslashesToSlash(input)
	if strings.Contains(input, "/") {
		return input
	}
	return "./" + input
}

//	Changes single and double backslashes into normal slashes
//	important for linux and other os
func convertBackslashesToSlash(input string) string {
	input = strings.ReplaceAll(input, "\\\\", "/")
	input = strings.ReplaceAll(input, "\\", "/")
	return input
}

func readPasswordFromCommandLine() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		PrintToConsole("Cannot get password")
		return "", err
	}
	return string(bytePassword), nil

}

func getResultJSONContent() string {
	resultJSONPath := "./output/result.json"
	return getFileContent(resultJSONPath)
}

func getFileContent(path string) string {
	if checkPathExists(path) {
		content, _ := os.ReadFile(path)
		return string(content)
	}

	errTxt := path + " does not exist."
	WriteErrorLog(errTxt, "")
	os.Exit(0)
	return ""
}

func inputHasZipEnding(input string) bool {
	zipEnding := ".zip"
	return inputHasEnding(input, zipEnding)
}

func convertInputToZip(input string) string {
	if !inputHasZipEnding(input) {
		input += ".zip"
	}
	return input
}
