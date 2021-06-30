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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
)

var VmCommand *goja.Runtime

var bigAudit BigAudit
var output string
var expectedTypes []string
var zipLocation string
var dontSaveArtefact bool
var auditResult bool
var sanityErr error
var debugModeEnabled bool
var pw string

// defines allow compare types
func init() {
	expectedTypes = []string{"==", "!=", ">=", ">", "<=", "<", "nil", "contains", "containsReg"}
	zipLocation = ""
	dontSaveArtefact = false
}

// compare specified system in config to actual system, quit if no match
func sanityCheck() error {
	osName := runtime.GOOS
	configSys := GetSystem()
	if !strings.EqualFold(osName, configSys.SystemName) {
		return errors.New("OS is: " + osName + ", expected: " + configSys.SystemName)
	}

	hasAdminPermissions := hasAdminPermissions()
	wantAdminPermissions := configSys.RootPermissions
	if hasAdminPermissions != wantAdminPermissions {
		if hasAdminPermissions {
			return errors.New("you are running the tool as root")
		}
		return errors.New("you are not root")
	}
	return nil
}

// get permissions through access rights
func hasAdminPermissions() bool {
	osName := runtime.GOOS
	// TODO Write debug  runtime.GOOS
	if strings.EqualFold(osName, "windows") {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	}
	currentUser, _ := user.Current()
	id, _ := currentUser.GroupIds()
	return id[0] == "0"
}

func main() {
	// query flags once
	flagsErr := initFlags()
	if flagsErr != nil {
		errTxt := flagsErr.Error()
		WriteErrorLog(errTxt, "")
		os.Exit(0)
	}

	debugModeEnabled = flags.debug

	if debugModeEnabled {
		WriteDebugLog(fmt.Sprint(flags), "INFO")
	}

	// TODO write debug toggelt flags
	if flags.help {
		printHelpText()
		os.Exit(0)
	}

	printBanner()

	FirstAuditEntry = false
	FirstErrorEntry = false
	FirstResultEntry = false

	removeOutputErr := os.RemoveAll("./output/artefacts")

	if removeOutputErr != nil {
		WriteErrorLog(removeOutputErr.Error(), "")
		WriteDebugLog(removeOutputErr.Error(), "ERROR")
	}

	if debugModeEnabled {
		WriteDebugLog("./output folder successfully deleted", "INFO")
	}

	if flags.encryptZip {
		PrintToConsole("Enter password for zip: ")
		pwFirst, pwErr := readPasswordFromCommandLine()
		pw = pwFirst
		if pwErr != nil {
			if debugModeEnabled {
				WriteDebugLog("password could not be read from command line", "ERROR")
			}
			WriteErrorLog(pwErr.Error(), "")
			os.Exit(0)
		}
		if pw == "" {
			PrintToConsole("Password cannot be empty String")
			if debugModeEnabled {
				WriteDebugLog("cannot use empty String as password", "ERROR")
			}
			WriteErrorLog("cannot use emtpy String as password", "")
			os.Exit(0)
		}
		if debugModeEnabled {
			WriteDebugLog("first password could be read from command line", "INFO")
		}
		PrintToConsole("Please repeat the password: ")
		pwAgain, pwErr := readPasswordFromCommandLine()
		if pwErr != nil {
			if debugModeEnabled {
				WriteDebugLog("password could not be read from command line", "ERROR")
			}
			WriteErrorLog(pwErr.Error(), "")
			os.Exit(0)
		}
		if debugModeEnabled {
			WriteDebugLog("second password could be read from command line", "INFO")
		}
		if pw != pwAgain {
			PrintToConsole("Password does not match")
			if debugModeEnabled {
				WriteDebugLog("password does not match", "ERROR")
			}
			WriteErrorLog("Password does not match", "")
			os.Exit(0)
		}
		if debugModeEnabled {
			WriteDebugLog("zipfile successfully encrypted", "INFO")
		}
	}

	ReadConfig()
	getAdditionalCommands()
	createCommandVM()

	if !flags.skipSanity {
		if debugModeEnabled {
			WriteDebugLog("starting sanity check", "INFO")
		}
		sanityErr = sanityCheck()
	} else {
		if debugModeEnabled {
			WriteDebugLog("Skipping sanity check", "INFO")
		}
	}
	if sanityErr != nil {
		errTxt := sanityErr.Error()
		fmt.Println(errTxt)
		WriteErrorLog(errTxt, "")
		if debugModeEnabled {
			WriteDebugLog(errTxt, "ERROR")
		}
		os.Exit(0)
	}
	if debugModeEnabled {
		WriteDebugLog("sanity check passed", "INFO")
	}
	WriteLog("sanity check passed", "INFO")

	bigAuditsValidErr := checkBigAuditsAreValid()
	if bigAuditsValidErr != nil {
		fmt.Println(bigAuditsValidErr)
		WriteErrorLog(bigAuditsValidErr.Error(), "")
		if debugModeEnabled {
			WriteDebugLog(bigAuditsValidErr.Error(), "ERROR")
		}
		os.Exit(0)
	}

	allAuditLength := len(GetBigAudits())

	for i, v := range GetBigAudits() {
		bigAudit = v

		if flags.verbose {
			printCommandStarted(i+1, allAuditLength)
		} else {
			printProgressBar(allAuditLength, i+1)
		}

		dontSaveArtefact = bigAudit.DontSaveArtefact
		executeErr := runCommand()

		// typeExpected = "" gets defaulted to the standard compare
		if bigAudit.TypeExpected == "" {
			bigAudit.TypeExpected = "=="
		}
		if executeErr != nil {
			if executeErr.Error() == "registry not found" || executeErr.Error() == "could not find given value in registry" {
				WriteResultJSON(v, true, false, output, executeErr.Error(), bigAudit.TypeExpected)
			} else {
				errString := executeErr.Error()
				errStringSplit := strings.Split(errString, "fromShell:")

				if len(errStringSplit) > 1 {
					errString = removeSuffix(errStringSplit[1])
					WriteResultJSON(v, false, false, output, errString, bigAudit.TypeExpected)
				} else {
					WriteResultJSON(v, false, false, output, "command not executed", bigAudit.TypeExpected)
				}
			}
			WriteErrorLog(strings.ReplaceAll(executeErr.Error(), "fromShell:", ""), bigAudit.Name+":")
			if debugModeEnabled {
				WriteDebugLog(bigAudit.Name+" command could not be executed: "+executeErr.Error(), "ERROR")

			}
			auditResult = false
		} else {
			if debugModeEnabled {
				WriteDebugLog(bigAudit.Name+" command was executed", "INFO")

			}
			auditSuccess, compareErr := compareOutput(v)
			if compareErr == nil && auditSuccess {
				WriteLog(bigAudit.Name+" finished! (Executed + output == expected)", "INFO")
				if debugModeEnabled {
					WriteDebugLog(bigAudit.Name+" output == expected", "INFO")
				}
				auditResult = true
			} else if compareErr == nil && !auditSuccess {
				WriteLog(bigAudit.Name+" output != expected", "FAIL")
				if debugModeEnabled {
					WriteDebugLog(bigAudit.Name+" output != expected", "ERROR")
				}
				auditResult = false
			} else {
				// Could not compare Output and Expected
				//goland:noinspection GoNilness
				WriteErrorLog(compareErr.Error(), bigAudit.Name)
				auditResult = false
			}
		}
		if flags.verbose {
			printCommandResult(auditResult)
		}

	}

	//create zip for output
	checkZipLocation()
	err := ZipFiles(zipLocation, GetAllFilesInOutput())
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}

}

func checkZipLocation() bool {
	var userInput string
	zipLocation = flags.output
	if zipLocation == "" {
		zipLocation = time.Now().Format("02.01.06_15-04-05" + ".zip")
		WriteLog("Zip folder: "+zipLocation+" created", "INFO")
		return true
	}
	if !inputHasZipEnding(zipLocation) {
		zipLocation = convertInputToZip(zipLocation)
	}

	fileinfo, err := os.Stat(zipLocation)

	if err == nil {
		if fileinfo.IsDir() {
			notDir := errors.New("the zip path is a directory")
			WriteErrorLog(notDir.Error(), "")
			fmt.Println(notDir)
			os.Exit(0)
		}
		fmt.Println("\nDo you want to overwrite(y/N):  " + fileinfo.Name())

		reader := bufio.NewReader(os.Stdin)
		userInputRune, _, err := reader.ReadRune()
		userInput = string(userInputRune)

		if err != nil {
			WriteErrorLog(err.Error(), "")
		}

		// only overwrite with permission
		switch strings.ToLower(userInput) {
		case "y":
			err := os.Remove(zipLocation)
			if err != nil {
				WriteErrorLog(err.Error(), "")
			} else {
				WriteLog(fileinfo.Name()+" overwritten", "INFO")
			}
		default:
			defaultErr := errors.New("could not create Zip folder")
			WriteErrorLog(defaultErr.Error(), "")
			os.Exit(0)
		}
	}
	return true
}

// compare given typeExpected to allow types
func checkExpectedType(expectedType string) bool {
	for _, v := range expectedTypes {
		if v == expectedType {
			return true
		}
	}
	return false
}

func validateOutputAndExpected(expectedType string, expected string) (bool, error) {
	outputInt, err1 := strconv.ParseInt(output, 10, 64)
	expectedInt, err2 := strconv.ParseInt(expected, 10, 64)
	// both are string
	if err1 != nil && err2 != nil {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" comparing String with String", "INFO")
		}
		if expectedType == "==" || expectedType == "!=" || expectedType == "contains" || expectedType == "containsReg" {

			if debugModeEnabled {
				WriteDebugLog(bigAudit.Name+" Strings compared with "+expectedType+" ok", "INFO")
			}

			if expectedType == "==" {
				if output == expected {
					return true, nil
				} else {
					return false, nil
				}
			} else if expectedType == "!=" {
				if output != expected {
					return true, nil
				} else {
					return false, nil
				}
			} else if expectedType == "containsReg" {
				re := regexp.MustCompile(expected)
				if re.MatchString(output) {
					return true, nil
				} else {
					return false, nil
				}
			} else {
				if strings.Contains(output, expected) {
					return true, nil
				} else {
					return false, nil
				}
			}
		} else {
			err := errors.New(expectedType + " cannot be used on a string")
			if debugModeEnabled {
				WriteDebugLog(bigAudit.Name+" cannot compare Strings with "+expectedType, "ERROR")
			}

			return false, err
		}
	} else if err1 == nil && err2 == nil {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" comparing Int with Int", "INFO")
		}
		// Has to be an int, compare without limitations
		switch expectedType {
		case "==":
			return outputInt == expectedInt, nil
		case "!=":
			return outputInt != expectedInt, nil
		case ">=":
			return outputInt >= expectedInt, nil
		case ">":
			return outputInt > expectedInt, nil
		case "<=":
			return outputInt <= expectedInt, nil
		case "<":
			return outputInt < expectedInt, nil
		}
	}
	if debugModeEnabled {
		WriteDebugLog(bigAudit.Name+" cannot compare String and Int", "ERROR")
	}
	return false, errors.New("cannot compare String and Int")
}

// way to determine if input was empty because it failed or actually empty
func compareOutput(audit BigAudit) (bool, error) {
	if output == "§NOTHING_WAS_RETURNED!§" {
		output = ""
	}

	if !checkExpectedType(bigAudit.TypeExpected) {
		err := errors.New(audit.Name + " wrong operator in TypeExpected")
		WriteCommandFailedLog(audit, err)
		WriteResultJSON(audit, false, false, output, err.Error(), "")
		if debugModeEnabled {
			WriteDebugLog(audit.Name+" wrong operator in TypeExpected", "ERROR")

		}
		return false, err
	}

	if debugModeEnabled {
		WriteDebugLog(audit.Name+" expectedType is correct", "INFO")
	}

	if audit.TypeExpected == "nil" {
		WriteCommandSuccessLog(audit)
		WriteResultJSON(audit, true, true, "", "", "")
		if debugModeEnabled {
			WriteDebugLog(audit.Name+" TypeExpected == nil, so output and expected validation is necessary", "INFO")
		}
		return true, nil
	}

	auditSuccessful, compareErr := validateOutputAndExpected(bigAudit.TypeExpected, audit.Expected)
	if compareErr != nil {
		WriteCommandFailedLog(audit, compareErr)
		WriteResultJSON(audit, true, auditSuccessful, output, compareErr.Error(), bigAudit.TypeExpected)
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+compareErr.Error(), "ERROR")
		}
	} else {
		WriteCommandSuccessLog(audit)
		WriteResultJSON(audit, true, auditSuccessful, output, "", bigAudit.TypeExpected)
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" syntax for comparison successful", "INFO")
		}
	}
	return auditSuccessful, compareErr
}

// run command in JavaScript
func runCommand() error {

	cmd := bigAudit.Command

	cmd = strings.ReplaceAll(cmd, "\\", "\\\\")
	cmd = strings.ReplaceAll(cmd, "`", "\\'")
	cmd = removeWhitespacePrefix(cmd)

	if len(cmd) == 0 {
		return errors.New("cannot execute command")
	}

	commandType := strings.Fields(cmd)[0]

	for _, v := range SupportedCommands {
		if strings.EqualFold(v, commandType) {
			cmd = "call('" + cmd + "')"
			if debugModeEnabled {
				WriteDebugLog(bigAudit.Name+".Command surrounded with "+cmd, "INFO")
			}
			break
		}
	}

	output = ""
	javaScriptOutput, err := VmCommand.RunString(cmd)

	if output == "" && javaScriptOutput != nil {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" output is empty", "INFO")
		}
		if javaScriptOutput.String() == "undefined" {
			output = "§NOTHING_WAS_RETURNED!§"
		} else {
			output = javaScriptOutput.String()
		}
	}

	if err != nil {

		err = errors.New(betterGojaError(err))
		return err
	}

	return nil
}

// function returned error or JS Syntax error
func betterGojaError(err error) string {
	errString := err.Error()
	if strings.Contains(errString, "Unexpected identifier") {
		return "SyntaxError: JavaScript"
	}

	if strings.Contains(errString, "GoError:") {
		errString = strings.ReplaceAll(errString, "GoError:", "")
		errString = strings.ReplaceAll(errString, " at main.Call (native)", "")
		errString = strings.ReplaceAll(errString, " at main.CallCompare (native)", "")
		errString = strings.ReplaceAll(errString, " at main.PrintToConsole (native)", "")
		errString = strings.ReplaceAll(errString, " at main.CallContains (native)", "")
		errString = strings.ReplaceAll(errString, " at main.RegQuery (native)", "")
		errString = strings.ReplaceAll(errString, " at main.Shell (native)", "")
		errString = removeWhitespacePrefix(errString)
	}
	return errString
}

// directly execute command in shell
// useful for commands that return errors as output -> see .log files
func Shell(cmd string) error {
	var auditErr error
	var out []byte
	var err error
	var errOut []byte

	if strings.EqualFold(runtime.GOOS, "windows") {
		cmdSlice := strings.Fields(cmd)
		helperSlice := []string{GetSystem().Argument}
		helperSlice = append(helperSlice, cmdSlice...)

		out, err = exec.Command(GetSystem().Shell, helperSlice...).Output()
	} else {
		out, err = exec.Command(GetSystem().Shell, GetSystem().Argument, cmd).Output()
	}

	execCmd := exec.Command(GetSystem().Shell, GetSystem().Argument, cmd)
	errOut, _ = execCmd.CombinedOutput()
	errOut = []byte(removeSuffix(string(errOut)))

	if err != nil {
		if len(errOut) > 0 {
			return errors.New("fromShell:" + string(errOut))
		}
		return err
	}

	output = removeSuffix(string(out))
	if !dontSaveArtefact {
		auditErr = saveAuditFileFromShell(output, bigAudit.Name)
	}

	if auditErr != nil {
		WriteErrorLog("cannot save artefact for: "+cmd, bigAudit.Name)
		return auditErr
	}
	WriteLog("artefact "+bigAudit.Name+" successfully saved", "INFO")

	return nil
}

// Execute command through our wrapper with good debugging
func Call(cmd string) error {
	var stderr bytes.Buffer

	var out []byte
	audits := separateInSmallAuditsButOnlyForCall(bigAudit.Name, cmd)

	var pipesTxt string
	for i := range audits {
		if i != len(audits)-1 {
			pipesTxt += audits[i].Command + ", "
		} else {
			pipesTxt += audits[i].Command
		}
	}
	WriteLog(bigAudit.Name+" separated in: "+pipesTxt, "INFO")

	wrappedAudits, failposition, err := AuditWrapper(audits...)

	if err != nil {
		WriteAuditFailedLog(audits, failposition+1)
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" wrapper failed on position "+fmt.Sprint(failposition+1), "ERROR")
		}
		return err
	}

	if debugModeEnabled {
		WriteDebugLog(bigAudit.Name+" wrapper successfull", "INFO")
	}

	if !dontSaveArtefact {
		saveArtefact(audits[0])
	}
	for i, v := range wrappedAudits {
		v.Stderr = &stderr
		out, _ = v.Output()

		if stderr.String() != "" {
			pipelineError := audits[i].Command + " failed"
			output = removeSuffix(string(out))
			WriteCommandFailedLog(bigAudit, errors.New(pipelineError))
			if debugModeEnabled {
				WriteDebugLog(audits[i].Command+" command cant be executed: "+pipelineError, "ERROR")
			}
			return errors.New(pipelineError)
		}

		if debugModeEnabled {
			WriteDebugLog(audits[i].Command+" command successfully executed", "INFO")
		}
	}
	output = removeSuffix(string(out))
	return nil
}

// trim trailing special characters
func removeSuffix(text string) string {
	for strings.HasSuffix(text, "\n") || strings.HasSuffix(text, "\r") || strings.HasSuffix(text, " ") {
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")
		text = strings.TrimSuffix(text, " ")
	}
	return text
}

func removeWhitespacePrefix(text string) string {
	for strings.HasPrefix(text, " ") {
		text = strings.TrimPrefix(text, " ")
	}
	return text
}

func PrintToConsole(text string) bool {
	fmt.Println(text)
	return true
}

// Directly write to log, can be used with JS
func PrintToLog(text string, logType string) bool {
	WriteLog(text, logType)
	if debugModeEnabled {
		WriteDebugLog(text, logType)
	}
	return true
}

//Should be called in a JS if Condition
//Extens Call Function, by comparing output string
func CallCompare(cmd string, expected string) (bool, error) {
	output = "§CALL_COMPARE_DOES_NOT_EXIST§"
	err := Call(cmd)
	if err != nil {
		return false, err
	}
	if output == expected {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" callCompare successful "+output+" == "+expected, "INFO")
		}
		return true, nil
	} else {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" callCompare not successful "+output+" != "+expected, "ERROR")
		}
		return false, nil
	}

}

//Should be called in a JS if Condition
//Extends Call Function, by checking if expected string is containt in the call Output string
func CallContains(cmd string, expected string) (bool, error) {
	output = "§CALL_CONTAIN_DOES_NOT_EXIST§"
	err := Call(cmd)
	if err != nil {
		return false, err
	}
	if strings.Contains(output, expected) {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" callContains successful "+output+" contains "+expected, "INFO")
		}
		return true, nil
	} else {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" callContains not successful "+output+" does not contain "+expected, "INFO")
		}
		return false, nil
	}
}
