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
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	flags.input = "configWin"
}

func TestConvertInputToJsonEqual(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{
			"filename",
			"filename.json",
		},
		{
			"",
			".json",
		},
		{
			"filenameWith.json",
			"filenameWith.json",
		},
		{
			"filenameWith.JSON",
			"filenameWith.JSON",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.expected, convertInputToJson(test.input), "Check \""+test.input+"\" converted to \""+test.expected+"\"")
	}
}

func TestHasJsonEndingTrue(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{
			"filename.json",
		},
		{
			".json",
		},
		{
			"filename.JSON",
		},
	}
	for _, test := range testCases {
		assert.True(t, inputHasJsonEnding(test.input), "Check \""+test.input+"\" has \".json\" ending")
	}
}

func TestHasJsonEndingFalse(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{
			"filenamejson",
		},
		{
			"",
		},
	}
	for _, test := range testCases {
		assert.False(t, inputHasJsonEnding(test.input), "Check \""+test.input+"\" has not \".json\" ending")
	}
}

func TestConvertInputToZipEqual(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{
			"filename",
			"filename.zip",
		},
		{
			"",
			".zip",
		},
		{
			"filenameWith.zip",
			"filenameWith.zip",
		},
		{
			"filenameWith.ZIP",
			"filenameWith.ZIP",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.expected, convertInputToZip(test.input), "Check \""+test.input+"\" converted to \""+test.expected+"\"")
	}
}

func TestHasZipEndingTrue(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{
			"filename.zip",
		},
		{
			".zip",
		},
		{
			"filename.ZIP",
		},
	}
	for _, test := range testCases {
		assert.True(t, inputHasZipEnding(test.input), "Check \""+test.input+"\" has \".zip\" ending")
	}
}

func TestHasZipEndingFalse(t *testing.T) {
	var testCases = []struct {
		input string
	}{
		{
			"filenamezip",
		},
		{
			"",
		},
	}
	for _, test := range testCases {
		assert.False(t, inputHasZipEnding(test.input), "Check \""+test.input+"\" has not \".zip\" ending")
	}
}

func TestConvertBackslashesToSlashEqual(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{
			"\\\\",
			"/",
		},
		{
			"\\",
			"/",
		},
		{
			"",
			"",
		},
		{
			"\\output\\testfile.json",
			"/output/testfile.json",
		},
		{
			"\\\\output\\testfile.json",
			"/output/testfile.json",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.expected, convertBackslashesToSlash(test.input), "Check \""+test.input+"\" converted to \""+test.expected+"\"")
	}
}

func TestConvertInputToPathEqual(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{
			"configWin.json",
			"./configWin.json",
		},
		{
			"configWin",
			"./configWin",
		},
		{
			"",
			"./",
		},
		{
			"/",
			"/",
		},
		{
			"C:\\Users\\Dummy\\Desktop\\CEP\\cep-seculeet\\seculeetApp\\configWin.json",
			"C:/Users/Dummy/Desktop/CEP/cep-seculeet/seculeetApp/configWin.json",
		},
		{
			"C:\\\\Users\\\\Dummy\\\\Desktop\\\\CEP\\\\cep-seculeet\\\\seculeetApp\\\\configWin.json",
			"C:/Users/Dummy/Desktop/CEP/cep-seculeet/seculeetApp/configWin.json",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.expected, convertInputToPath(test.input), "Check \""+test.input+"\" converted to \""+test.expected+"\"")
	}
}

func TestReadInputWithInputFlag(t *testing.T) {
	var testCases = []struct {
		inputFlag         string
		expectedInputFlag string
	}{
		{
			"configFile",
			"./configFile.json",
		},
		{
			"configFile.json",
			"./configFile.json",
		},
		{
			"input/configFile.json",
			"input/configFile.json",
		},
		{
			"/input/configFile.json",
			"/input/configFile.json",
		},
		{
			"./input/configFile.JSON",
			"./input/configFile.JSON",
		},
	}
	for _, test := range testCases {
		flags.input = test.inputFlag
		readInputFlag, _, _ := readInput()
		assert.Equal(t, test.expectedInputFlag, readInputFlag, "Check \""+readInputFlag+"\" is equal to \""+test.expectedInputFlag+"\"")
	}
	deleteOutput()
}

func TestReadInputNoInputFlag(t *testing.T) {
	flags.input = ""
	assert.Panics(t, func() { readInput() }, "Check \"readInput()\" throws panic if \"inputFlag\" is empty")
	expectedErrorTxt := "you need to specify an input file"
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expectedErrorTxt, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expectedErrorTxt, nil)

	deleteOutput()
}

func TestReadInputNotExistingInputFile(t *testing.T) {
	flags.input = "NotExistingInputFile"
	assert.NotPanics(t, func() { readInput() }, "Check \"readInput()\" run without panics")

	deleteOutput()
}

func TestInputHasEndingTrue(t *testing.T) {
	var testCases = []struct {
		filename string
		ending   string
	}{
		{
			"file.json",
			".json",
		},
		{
			".json",
			".json",
		},
		{
			".json",
			".JSON",
		},
		{
			".zip",
			"zip",
		},
		{
			".bla",
			"bLa",
		},
		{
			"file.json",
			"",
		},
		{
			"",
			"",
		},
	}
	for _, test := range testCases {
		assert.True(t, inputHasEnding(test.filename, test.ending), "Check \""+test.filename+"\" has \""+test.ending+"\" ending")
	}
}

func TestInputHasEndingFalse(t *testing.T) {
	var testCases = []struct {
		filename string
		ending   string
	}{
		{
			".jsonFile",
			"json",
		},
		{
			"jsonFile",
			"json",
		},
		{
			"file.jayson",
			"json",
		},
	}
	for _, test := range testCases {
		assert.False(t, inputHasEnding(test.filename, test.ending), "Check \""+test.filename+"\" has not \""+test.ending+"\" ending")
	}
}

func TestReadConfigWithInputFlag(t *testing.T) {
	configFileName := "testConf.json"
	flags.input = "./output/" + configFileName

	auditName := "test_audit"
	auditCommand := "test command"
	auditDontSaveArtefact := true
	auditTypeExpected := "=="
	auditExpected := "true"
	auditDesc := "test description"

	sysName := "Windows"
	sysVersion := "10.0.19042"
	sysShell := "powershell"
	sysArgument := "/C"

	configFileContent := `{
		"commands": [
		  {
			"name": "` + auditName + `",
			"command": "` + auditCommand + `",
			"dontSaveArtefact": ` + strconv.FormatBool(auditDontSaveArtefact) + `,
			"typeExpected": "` + auditTypeExpected + `",
			"expected": "` + auditExpected + `",
			"description": "` + auditDesc + `"
		  }
		],
		"system":
		  {
			"systemName": "` + sysName + `",
			"version": "` + sysVersion + `",
			"shell": "` + sysShell + `",
			"argument": "` + sysArgument + `"
		  }
	  }`
	fileWriter(configFileContent, configFileName, false)
	assert.NotPanics(t, func() { ReadConfig() }, "Check \"ReadConfig()\" run without panics")

	firstAudit := GetBigAudits()[0]
	assert.Equal(t, auditName, firstAudit.Name, "Check \""+firstAudit.Name+"\" is equal to \""+auditName+"\"")
	assert.Equal(t, auditCommand, firstAudit.Command, "Check \""+firstAudit.Command+"\" is equal to \""+auditCommand+"\"")
	assert.Equal(t, auditDontSaveArtefact, firstAudit.DontSaveArtefact, "Check \""+strconv.FormatBool(firstAudit.DontSaveArtefact)+"\" is equal to \""+strconv.FormatBool(auditDontSaveArtefact)+"\"")
	assert.Equal(t, auditTypeExpected, firstAudit.TypeExpected, "Check \""+firstAudit.TypeExpected+"\" is equal to \""+auditTypeExpected+"\"")
	assert.Equal(t, auditExpected, firstAudit.Expected, "Check \""+firstAudit.Expected+"\" is equal to \""+auditExpected+"\"")
	assert.Equal(t, auditDesc, firstAudit.Desc, "Check \""+firstAudit.Desc+"\" is equal to \""+auditDesc+"\"")

	systemInfo := GetSystem()
	assert.Equal(t, sysName, systemInfo.SystemName, "Check \""+systemInfo.SystemName+"\" is equal to \""+sysName+"\"")
	assert.Equal(t, sysVersion, systemInfo.Version, "Check \""+systemInfo.Version+"\" is equal to \""+sysVersion+"\"")
	assert.Equal(t, sysShell, systemInfo.Shell, "Check \""+systemInfo.Shell+"\" is equal to \""+sysShell+"\"")
	assert.Equal(t, sysArgument, systemInfo.Argument, "Check \""+systemInfo.Argument+"\" is equal to \""+sysArgument+"\"")

	deleteFile(flags.input)
	deleteOutput()
}

func TestReadConfigWithInputFlagAndValidJSONNonValidSystem(t *testing.T) {
	configFileName := "testConfNoSystem.json"
	flags.input = "./output/" + configFileName

	configFileContent := `{
		"commands": [
		  {
			"name": "test_audit",
			"command": "test command",
			"typeExpected": "==",
			"expected": "true",
			"description": "test description"
		  }
		],
		"system": {
		}
	  }`
	fileWriter(configFileContent, configFileName, false)
	var tmpSystem11 System
	json.Unmarshal([]byte(configFileContent), &tmpSystem11)
	system = tmpSystem11
	assert.Panics(t, func() { ReadConfig() }, "Check \"ReadConfig()\" throws panic if system not specified in config file")
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, "[ERROR]: you have to specify the system in your config file", nil)

	deleteFile(flags.input)
	deleteOutput()
}

func TestReadConfigWithInputFlagNonValidJSON(t *testing.T) {
	nonValidConfigFile := "configWinNonValid.json"
	flags.input = "./output/" + nonValidConfigFile
	nonValidJSON := "}{"
	fileWriter(nonValidJSON, nonValidConfigFile, false)
	assert.Panics(t, func() { ReadConfig() }, "Check \"ReadConfig()\" throws panic if config file JSON format is not valid")
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, "[ERROR]: JSON format incorrect", nil)

	deleteFile(flags.input)
	deleteOutput()
}

func TestReadConfigNoInputFlag(t *testing.T) {
	flags.input = ""
	assert.Panics(t, func() { ReadConfig() }, "Check \"ReadConfig()\" throws panic if \"inputFlag\" is empty")
	expectedErrorTxt := "you need to specify an input file"
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expectedErrorTxt, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expectedErrorTxt, nil)

	deleteOutput()
}

func TestReadConfigNotExistingInputFile(t *testing.T) {
	flags.input = "./input/NotExistingInputFile.json"
	assert.Panics(t, func() { ReadConfig() }, "Check \"ReadConfig()\" throws panic if the specified input file not exist")
	expectedErrorTxt := "open " + flags.input
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expectedErrorTxt, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expectedErrorTxt, nil)

	deleteOutput()
}

func TestCheckJsonFormatTrue(t *testing.T) {
	var testCases = []struct {
		json string
	}{
		{
			`{
				"commands":[
				   {
					  "name":"check_antivir_enabled",
					  "command":"callContains('Get-MpComputerStatus | findstr AntivirusEnabled', 'True')",
					  "typeExpected":"contains",
					  "expected":"True",
					  "description":"Checking Antivir is enabled"
				   }
				]
			 }`,
		},
		{
			`{
				"commands": [
					{
						"name": "check_firewall_profile_private",
						"command": "callContains('Get-NetFirewallProfile -Name Private | findstr Enabled', 'True')",
						"typeExpected": "contains",
						"expected": "True",
						"description": "Checking Firewall enabled on Private"
					},
					{
						"name": "check_BTAGService_is_disabled",
						"command": "Get-ItemPropertyValue Registry::HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\BTAGService -Name Start",
						"typeExpected": "==",
						"expected": "4",
						"description": "Checking Bluetooth Audio Gateway Service is disabled",
						"dontSaveArtefact": true
					}
				]
			}`,
		},
	}
	var testStruct BigAudits
	for _, test := range testCases {
		assert.True(t, checkJsonFormat([]byte(test.json), &testStruct), "Check \"json\" format is valid")
	}
}

func TestCheckJsonFormatFalse(t *testing.T) {
	var testCases = []struct {
		json string
	}{
		{
			`{
				"commands":[
				   {
					  "name":"check_antivir_enabled",
					  "command":"callContains('Get-MpComputerStatus | findstr AntivirusEnabled', 'True')",
				   }
				]
			 }`,
		},
		{
			``,
		},
		{
			`{
				"commands":[
				   {
					  "name":"check_antivir_enabled",
					  "command":"callContains('Get-MpComputerStatus | findstr AntivirusEnabled', 'True')"
				   }
			 }`,
		},
	}
	var testStruct BigAudits
	for _, test := range testCases {
		assert.False(t, checkJsonFormat([]byte(test.json), &testStruct), "Check \"json\" format is not valid")
	}
}

func TestGetFileContentWithValidPath(t *testing.T) {
	filename := "validFile.json"
	filePath := "./output/" + filename
	content := "{}"
	fileWriter(content, filename, false)
	assert.NotPanics(t, func() { getFileContent(filePath) }, "Check \"getFileContent()\" run without panics")
	fileContent := getFileContent(filePath)
	assert.Equal(t, content, fileContent, "Check the file content is equal to expected content")

	deleteOutput()
}

func TestGetFileContentNonValidPath(t *testing.T) {
	path := "./output/nonValidPath.file"
	assert.Panics(t, func() { getFileContent(path) }, "Check \"getFileContent()\" throws panic if the specified file does not exist")
	expectedErrorTxt := path + " does not exist."
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expectedErrorTxt, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expectedErrorTxt, nil)

	deleteOutput()
}

func TestCheckSystemIsValid(t *testing.T) {
	configJSONNoSystem := `{}`
	var tmpSystem1 System
	json.Unmarshal([]byte(configJSONNoSystem), &tmpSystem1)
	system = tmpSystem1
	err := checkSystemIsValid()
	assert.EqualError(t, err, "you have to specify the system in your config file", "Check \"checkSystemIsValid()\" throws expected error")

	configJSONNoSystemName := `{
		"system":
		{
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C"
		}
	}`
	var tmpSystem2 System
	json.Unmarshal([]byte(configJSONNoSystemName), &tmpSystem2)
	system = tmpSystem2
	err = checkSystemIsValid()
	assert.EqualError(t, err, "you have to specify the systemName in your config file", "Check \"checkSystemIsValid()\" throws expected error")

	configJSONNoSystemNameValue := `{
		"system":
		{
		  "systemName": "",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C"
		}
	}`
	var tmpSystem3 System
	json.Unmarshal([]byte(configJSONNoSystemNameValue), &tmpSystem3)
	system = tmpSystem3
	err = checkSystemIsValid()
	assert.EqualError(t, err, "you have to specify the systemName in your config file", "Check \"checkSystemIsValid()\" throws expected error")

	configJSONNoVersion := `{
		"system":
		{
		  "systemName": "Windows",
		  "shell": "powershell",
		  "argument": "/C"
		}
	}`
	var tmpSystem4 System
	json.Unmarshal([]byte(configJSONNoVersion), &tmpSystem4)
	system = tmpSystem4
	err = checkSystemIsValid()
	assert.EqualError(t, err, "you have to specify the version in your config file", "Check \"checkSystemIsValid()\" throws expected error")

	configJSONNoVersionValue := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "",
		  "shell": "powershell",
		  "argument": "/C"
		}
	}`
	var tmpSystem5 System
	json.Unmarshal([]byte(configJSONNoVersionValue), &tmpSystem5)
	system = tmpSystem5
	err = checkSystemIsValid()
	assert.EqualError(t, err, "you have to specify the version in your config file", "Check \"checkSystemIsValid()\" throws expected error")

	configJSONNoShell := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "argument": "/C"
		}
	}`
	var tmpSystem6 System
	json.Unmarshal([]byte(configJSONNoShell), &tmpSystem6)
	system = tmpSystem6
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	if strings.EqualFold(runtime.GOOS, "windows") {
		assert.Equal(t, system.System.Shell, "powershell", "Checks if the shell is not specified, whether the default value was set")
	} else {
		assert.Equal(t, system.System.Shell, "bash", "Checks if the shell is not specified, whether the default value was set")
	}

	configJSONNoShellValue := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "",
		  "argument": "/C"
		}
	}`
	var tmpSystem7 System
	json.Unmarshal([]byte(configJSONNoShellValue), &tmpSystem7)
	system = tmpSystem7
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	if strings.EqualFold(runtime.GOOS, "windows") {
		assert.Equal(t, system.System.Shell, "powershell", "Checks if the shell is not specified, whether the default value was set")
	} else {
		assert.Equal(t, system.System.Shell, "bash", "Checks if the shell is not specified, whether the default value was set")
	}

	configJSONNoArgument := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell"
		}
	}`
	var tmpSystem8 System
	json.Unmarshal([]byte(configJSONNoArgument), &tmpSystem8)
	system = tmpSystem8
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	if strings.EqualFold(runtime.GOOS, "windows") {
		assert.Equal(t, system.System.Argument, "/C", "Checks if the argument is not specified, whether the default value was set")
	} else {
		assert.Equal(t, system.System.Argument, "-c", "Checks if the argument is not specified, whether the default value was set")
	}

	configJSONNoArgumentValue := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": ""
		}
	}`
	var tmpSystem9 System
	json.Unmarshal([]byte(configJSONNoArgumentValue), &tmpSystem9)
	system = tmpSystem9
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	if strings.EqualFold(runtime.GOOS, "windows") {
		assert.Equal(t, system.System.Argument, "/C", "Checks if the argument is not specified, whether the default value was set")
	} else {
		assert.Equal(t, system.System.Argument, "-c", "Checks if the argument is not specified, whether the default value was set")
	}

	configJSONValid := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C"
		}
	}`
	var tmpSystem10 System
	json.Unmarshal([]byte(configJSONValid), &tmpSystem10)
	system = tmpSystem10
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	assert.False(t, system.System.RootPermissions, "Checks if root is not specified, whether the default value was set")

	configJSONNoRootValue := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C",
		  "root": ""
		}
	}`
	var tmpSystem11 System
	json.Unmarshal([]byte(configJSONNoRootValue), &tmpSystem11)
	system = tmpSystem11
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	assert.False(t, system.System.RootPermissions, "Checks if root is not specified, whether the default value was set")

	configJSONRootValueAsString := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C",
		  "root": "true"
		}
	}`
	var tmpSystem12 System
	json.Unmarshal([]byte(configJSONRootValueAsString), &tmpSystem12)
	system = tmpSystem12
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	assert.False(t, system.System.RootPermissions, "Checks if root is not specified, whether the default value was set")

	configJSONValidWithRoot := `{
		"system":
		{
		  "systemName": "Windows",
		  "version": "10.0.19042",
		  "shell": "powershell",
		  "argument": "/C",
		  "root": true
		}
	}`
	var tmpSystem13 System
	json.Unmarshal([]byte(configJSONValidWithRoot), &tmpSystem13)
	system = tmpSystem13
	err = checkSystemIsValid()
	assert.NoError(t, err, "Check \"checkSystemIsValid()\" run without error")
	assert.True(t, system.System.RootPermissions, "Check if \"RootPermissions\" is set correctly")
}

func TestGetOrdinalNum(t *testing.T) {
	var testCases = []struct {
		inputNum int
		expected string
		errorTxt string
	}{
		{
			-1,
			"",
			"The entered number must be greater than 0",
		},
		{
			0,
			"",
			"The entered number must be greater than 0",
		},
		{
			1,
			"1st",
			"",
		},
		{
			2,
			"2nd",
			"",
		},
		{
			3,
			"3rd",
			"",
		},
		{
			4,
			"4th",
			"",
		},
		{
			10,
			"10th",
			"",
		},
		{
			21,
			"21st",
			"",
		},
		{
			22,
			"22nd",
			"",
		},
		{
			23,
			"23rd",
			"",
		},
		{
			100,
			"100th",
			"",
		},
		{
			101,
			"101st",
			"",
		},
	}
	for _, test := range testCases {
		if test.errorTxt == "" {
			assert.NotPanics(t, func() { getOrdinalNum(test.inputNum) }, getOrdinalNum, "Check \"getOrdinalNum()\" run without panics")
			ordinalNum := getOrdinalNum(test.inputNum)
			assert.Equal(t, test.expected, ordinalNum, "Check \""+ordinalNum+"\" is equal to \""+test.expected+"\"")
		} else {
			assert.Panics(t, func() { getOrdinalNum(test.inputNum) }, "Check \"getOrdinalNum()\" throws expected error")
		}
	}
	deleteOutput()
}

func TestCheckBigAuditsAreValid(t *testing.T) {
	configJSONNoCommands := `{}`
	var tmpCmd1 BigAudits
	json.Unmarshal([]byte(configJSONNoCommands), &tmpCmd1)
	commands = tmpCmd1
	err := checkBigAuditsAreValid()
	expectedErrorTxt := "you have to specify at least one command in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONNoAudits := `{
		"commands": [
	  ]
	}`
	var tmpCmd2 BigAudits
	json.Unmarshal([]byte(configJSONNoAudits), &tmpCmd2)
	commands = tmpCmd2
	err = checkBigAuditsAreValid()
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONEmptyAudit := `{
		"commands": [
			{}
	  ]
	}`
	var tmpCmd3 BigAudits
	json.Unmarshal([]byte(configJSONEmptyAudit), &tmpCmd3)
	commands = tmpCmd3
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 1st audit. You have to specify the name in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONAuditNoName := `{
		"commands": [
			{
				"command": "echo test",
				"typeExpected": "nil",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd4 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoName), &tmpCmd4)
	commands = tmpCmd4
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 1st audit. You have to specify the name in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONAuditNoNameValue := `{
		"commands": [
			{
				"name": "",
				"command": "echo test",
				"typeExpected": "nil",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd5 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoNameValue), &tmpCmd5)
	commands = tmpCmd5
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 1st audit. You have to specify the name in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONAuditNoCommand := `{
		"commands": [
			{
				"name": "ein_test",
				"typeExpected": "nil",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd6 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoCommand), &tmpCmd6)
	commands = tmpCmd6
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 1st audit. You have to specify the command in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONAuditNoCommandValue := `{
		"commands": [
			{
				"name": "ein_test",
				"command": "",
				"typeExpected": "nil",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd7 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoCommandValue), &tmpCmd7)
	commands = tmpCmd7
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 1st audit. You have to specify the command in your config file"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")

	configJSONAuditNoTypeExpected := `{
		"commands": [
			{
				"name": "ein_test",
				"command": "echo test",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd8 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoTypeExpected), &tmpCmd8)
	commands = tmpCmd8
	err = checkBigAuditsAreValid()
	assert.NoError(t, err, "Check \"checkBigAuditsAreValid()\" run without error")

	configJSONAuditNoTypeExpectedValue := `{
		"commands": [
			{
				"name": "ein_test",
				"command": "echo test",
				"typeExpected": "",
				"expected": "3",
				"description": "The description"
			}
	  ]
	}`
	var tmpCmd9 BigAudits
	json.Unmarshal([]byte(configJSONAuditNoTypeExpectedValue), &tmpCmd9)
	commands = tmpCmd9
	err = checkBigAuditsAreValid()
	assert.NoError(t, err, "Check \"checkBigAuditsAreValid()\" run without error")

	configJSONAlreadyUsedAuditName := `{
		"commands": [
			{
				"name": "ein_test",
				"command": "echo test",
				"typeExpected": "nil",
				"expected": "3",
				"description": "The description"
			},
			{
				"name": "ein_test",
				"command": "echo test123",
				"typeExpected": "nil",
				"expected": "4",
				"description": "The description 2"
			}
	  ]
	}`
	var tmpCmd10 BigAudits
	json.Unmarshal([]byte(configJSONAlreadyUsedAuditName), &tmpCmd10)
	commands = tmpCmd10
	err = checkBigAuditsAreValid()
	expectedErrorTxt = "Issue at the 2nd audit. The audit name \"ein_test\" was already used in the 1st audit"
	assert.EqualError(t, err, expectedErrorTxt, "Check \"checkBigAuditsAreValid()\" throws expected error")
}
