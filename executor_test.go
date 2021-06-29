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
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestSanityCheck(t *testing.T) {
	var fileNameConfig string

	fileNameConfig = "configWinTestCommands.json"
	flags.input = "./output/" + fileNameConfig
	configWinTestCommands := `{
        "commands": [
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configWinTestCommands, fileNameConfig, false)
	ReadConfig()
	assert.Nil(t, sanityCheck())
	deleteFile(flags.input)

	fileNameConfig = "configWinTestCommands.json"
	flags.input = "./output/" + fileNameConfig
	configWinTestNoSystemName := `{
        "commands": [
        ],
        "system": 
        {
          "systemName": "NotWindows",
          "version": "10.0.19042"
        }
      }`

	fileWriter(configWinTestNoSystemName, fileNameConfig, false)
	ReadConfig()
	assert.Error(t, sanityCheck())
	deleteFile(flags.input)

	fileNameConfig = "configWinTestRoot.json"
	flags.input = "./output/" + fileNameConfig
	configWinTestRoot := `{
        "commands": [
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `,
			"root": true
        }
      }`

	fileWriter(configWinTestRoot, fileNameConfig, false)
	ReadConfig()
	assert.EqualError(t, sanityCheck(), "you are not root")

	deleteFile(flags.input)
	deleteOutput()
}

func TestHasAdminPermissions(t *testing.T) {
	a := hasAdminPermissions()
	assert.False(t, a)
}

func TestCheckZipLocation(t *testing.T) {

	assert.True(t, checkZipLocation())

	flags.output = "Testzip"
	assert.True(t, checkZipLocation())

	os.Create("ExistingZip.zip")
	flags.output = "ExistingZip"
	assert.Panics(t, func() { checkZipLocation() })
	CheckFileExists(t, pathAudit)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathAudit, "[ERROR]: EOF", nil)
	CheckFileContent(t, pathAudit, "[ERROR]: could not create Zip folder", nil)
	CheckFileContent(t, pathError, "[ERROR]: EOF", nil)
	CheckFileContent(t, pathError, "[ERROR]: could not create Zip folder", nil)

	os.Mkdir("ZipFolder.zip", 0777)
	flags.output = "ZipFolder"
	assert.Panics(t, func() { checkZipLocation() })
	CheckFileContent(t, pathAudit, "[ERROR]: the zip path is a directory", nil)
	CheckFileContent(t, pathError, "[ERROR]: the zip path is a directory", nil)

	deleteFile("ExistingZip.zip")
	os.RemoveAll("./ZipFolder.zip")
	deleteOutput()
}

func TestRemoveSuffix(t *testing.T) {
	var testCases = []struct {
		text     string
		expected string
	}{
		{
			"this is lorem ipsum text\n \n\r\n",
			"this is lorem ipsum text",
		},
		{
			"this is \nlorem ipsum text \n\r \n",
			"this is \nlorem ipsum text",
		},
		{
			"this is lorem ipsum text\r \n\r ",
			"this is lorem ipsum text",
		},
		{
			"",
			"",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, test.expected, removeSuffix(test.text), "Check \""+test.text+"\" is equal to \""+test.expected+"\"")
	}
}

func TestCheckExpectedType(t *testing.T) {
	var testCasesTrue = []struct {
		expectedType string
	}{
		{
			"!=",
		},
		{
			">=",
		},
		{
			">",
		},
		{
			"<=",
		},
		{
			"<",
		},
		{
			"nil",
		},
		{
			"contains",
		},
		{
			"containsReg",
		},
	}
	for _, test := range testCasesTrue {
		assert.True(t, checkExpectedType(test.expectedType))
	}

	var testCasesFalse = []struct {
		expectedType string
	}{
		{
			"",
		},
		{
			"abc",
		},
		{
			"===",
		},
		{
			"=!",
		},
		{
			">>",
		},
		{
			"Contains",
		},
		{
			"=",
		},
	}

	for _, test := range testCasesFalse {
		assert.False(t, checkExpectedType(test.expectedType))
	}
}

func TestValidateOutputAndExpected(t *testing.T) {
	output = "ABC"
	WrongExpectedType := "NotAnExpectedType"

	a, b := validateOutputAndExpected("==", "ABC")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("==", "NotABC")
	assert.False(t, a)

	a, b = validateOutputAndExpected("!=", "ABC")
	assert.False(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("!=", "NotABC")
	assert.True(t, a)

	a, b = validateOutputAndExpected("containsReg", "ABC")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("containsReg", "NotABC")
	assert.False(t, a)

	a, b = validateOutputAndExpected("contains", "ABC")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("contains", "NotABC")
	assert.False(t, a)

	a, b = validateOutputAndExpected(">=", "ABC")
	assert.False(t, a)
	assert.EqualError(t, b, ">= cannot be used on a string")

	_, b = validateOutputAndExpected(">=", "NotABC")
	assert.EqualError(t, b, ">= cannot be used on a string")

	a, b = validateOutputAndExpected(WrongExpectedType, "ABC")
	assert.False(t, a)
	assert.EqualError(t, b, "NotAnExpectedType cannot be used on a string")

	_, b = validateOutputAndExpected(WrongExpectedType, "NotABC")
	assert.EqualError(t, b, "NotAnExpectedType cannot be used on a string")

	output = "123"

	a, b = validateOutputAndExpected("==", "123")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("==", "124")
	assert.False(t, a)

	a, b = validateOutputAndExpected("!=", "123")
	assert.False(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("!=", "124")
	assert.True(t, a)

	a, b = validateOutputAndExpected(">=", "123")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected(">=", "122")
	assert.True(t, a)

	a, _ = validateOutputAndExpected(">=", "124")
	assert.False(t, a)

	a, b = validateOutputAndExpected(">", "123")
	assert.False(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected(">", "122")
	assert.True(t, a)

	a, _ = validateOutputAndExpected(">", "124")
	assert.False(t, a)

	a, b = validateOutputAndExpected("<=", "123")
	assert.True(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("<=", "122")
	assert.False(t, a)

	a, _ = validateOutputAndExpected("<=", "124")
	assert.True(t, a)

	a, b = validateOutputAndExpected("<", "123")
	assert.False(t, a)
	assert.Nil(t, b)

	a, _ = validateOutputAndExpected("<", "122")
	assert.False(t, a)

	a, _ = validateOutputAndExpected("<", "124")
	assert.True(t, a)

	a, b = validateOutputAndExpected("==", "ABC")
	assert.False(t, a)
	assert.EqualError(t, b, "cannot compare String and Int")
}

func TestCompareOutput(t *testing.T) {
	output = "§NOTHING_WAS_RETURNED!§"

	ExecutedArray := []string{"[INFO] : Name: TestName", "TestName: TestCommand executed"}
	FailedArray := []string{"[FAIL] : Name: TestName", "TestName: TestCommand failed"}

	ExecutedResult := `{
		"./output/configCommandExecuted.json": [
			{
				"Name": "TestName",
				"Command": "TestCommand",
				"Command was executed": true,
				"Output is as expected": true
			}
		]
	}`
	WrongTypeExpectedResult := `{
		"./output/configWrongTypeExpected.json": [
			{
				"Name": "TestName",
				"Command": "TestCommand",
				"Command was executed": false,
				"Error-Message": "TestName wrong operator in TypeExpected"
			}
		]
	}`
	NumericTypeExpectedResult := `{
		"./output/configNumericTypeExpected.json": [
			{
				"Name": "TestName",
				"Command": "TestCommand",
				"Command was executed": true,
				"Output is as expected": false,
				"Expected Value": "TestExpected",
				"Actual Value": "TestExpected",
				"Operator": ">="
			}
		]
	}`

	fileNameConfig := "configCommandExecuted.json"
	flags.input = "./output/" + fileNameConfig
	configTypeExpectedNil := `{
        "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": "nil"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configTypeExpectedNil, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]
	a, b := compareOutput(GetBigAudits()[0])
	assert.Empty(t, output)
	assert.True(t, a)
	assert.Nil(t, b)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", ExecutedArray)

	CheckFileExists(t, pathResult)
	CheckFileContent(t, pathResult, ExecutedResult, nil)

	deleteOutput()

	fileNameConfig = "configWrongTypeExpected.json"
	flags.input = "./output/" + fileNameConfig
	configWrongTypeExpected := `{
        "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": "WrongTypeExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configWrongTypeExpected, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]
	a, b = compareOutput(GetBigAudits()[0])
	assert.False(t, a)
	assert.EqualError(t, b, "TestName wrong operator in TypeExpected")

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", FailedArray)

	CheckFileExists(t, pathResult)
	CheckFileContent(t, pathResult, WrongTypeExpectedResult, nil)

	deleteOutput()

	fileNameConfig = "configCommandExecuted.json"
	flags.input = "./output/" + fileNameConfig
	configEmptyTypeExpected := `{
        "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configEmptyTypeExpected, fileNameConfig, false)
	ReadConfig()

	output = "TestExpected"
	bigAudit = GetBigAudits()[0]
	a, b = compareOutput(GetBigAudits()[0])
	assert.True(t, a)
	assert.Nil(t, b)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", ExecutedArray)

	CheckFileExists(t, pathResult)
	CheckFileContent(t, pathResult, ExecutedResult, nil)

	deleteOutput()

	fileNameConfig = "configNumericTypeExpected.json"
	flags.input = "./output/" + fileNameConfig
	configNumericTypeExpected := `{
        "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": ">=",
				"expected": "TestExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configNumericTypeExpected, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]
	a, b = compareOutput(GetBigAudits()[0])
	assert.False(t, a)
	assert.EqualError(t, b, ">= cannot be used on a string")

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", FailedArray)

	CheckFileExists(t, pathResult)
	CheckFileContent(t, pathResult, NumericTypeExpectedResult, nil)

	deleteOutput()
}

func TestRunCommand(t *testing.T) {

	createCommandVM()

	fileNameConfig := "configCommandExecuted.json"
	flags.input = "./output/" + fileNameConfig
	configEmptyCommand := `{
        "commands": [
			{
				"name": "TestName",
				"command": "",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configEmptyCommand, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]

	assert.EqualError(t, runCommand(), "cannot execute command")

	deleteOutput()

	fileNameConfig = "configCommandExecuted.json"
	flags.input = "./output/" + fileNameConfig
	configLSCommand := `{
        "commands": [
			{
				"name": "TestName",
				"command": "ls",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configLSCommand, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]

	assert.Nil(t, runCommand())
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[INFO] : TestName separated in: ls", nil)

	deleteOutput()

	fileNameConfig = "configCommandExecuted.json"
	flags.input = "./output/" + fileNameConfig
	configWrongCommand := `{
        "commands": [
			{
				"name": "TestName",
				"command": "WrongCommand",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	fileWriter(configWrongCommand, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]

	assert.Error(t, runCommand())

	deleteOutput()
}

func TestBetterGojaError(t *testing.T) {

	err1 := errors.New("Unexpected identifier")
	err2 := errors.New("GoError: at main.Call (native)")
	err3 := errors.New("GoError: at main.CallCompare (native)")
	err4 := errors.New("GoError: at main.PrintToConsole (native)")
	err5 := errors.New("GoError: at main.CallContains (native)")
	err6 := errors.New("GoError: at main.RegQuery (native)")
	err7 := errors.New("GoError: at main.Shell (native)")

	assert.Equal(t, betterGojaError(err1), "SyntaxError: JavaScript")
	assert.Empty(t, betterGojaError(err2))
	assert.Empty(t, betterGojaError(err3))
	assert.Empty(t, betterGojaError(err4))
	assert.Empty(t, betterGojaError(err5))
	assert.Empty(t, betterGojaError(err6))
	assert.Empty(t, betterGojaError(err7))
}

func TestShell(t *testing.T) {

	dontSaveArtefact = false

	fileNameConfig := "configTestCommand.json"
	flags.input = "./output/" + fileNameConfig
	configTestCommand := `{
	    "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
	    ],
	    "system":
	    {
			` + getConfigSystemForOS() + `
	    }
	  }`

	fileWriter(configTestCommand, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]
	assert.Nil(t, (Shell("ls")))
	assert.DirExists(t, "./output/artefacts")
	assert.FileExists(t, "./output/artefacts/TestName.txt")
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[INFO] : artefact TestName successfully saved", nil)
	assert.Error(t, (Shell("WrongCommand")))

	deleteOutput()

	bigAudit = GetBigAudits()[0]
	os.Mkdir(pathOutput, 07777)
	os.Mkdir("./output/artefacts", 07777)
	os.Create("./output/artefacts/TestName.txt")
	os.Chmod("./output/artefacts/TestName.txt", 0000)
	assert.Error(t, Shell("ls"))
	assert.DirExists(t, "./output/artefacts")
	assert.FileExists(t, "./output/artefacts/TestName.txt")
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[ERROR]: TestName cannot save artefact for: ls", nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, "[ERROR]: TestName cannot save artefact for: ls", nil)

	deleteOutput()
}

func TestCall(t *testing.T) {

	doubleLSArray := []string{"[INFO] : TestName separated in: ls, ls", "[FAIL] : Name: TestName", "TestName: TestCommand failed"}
	fileNameConfig := "configTestCommand.json"
	flags.input = "./output/" + fileNameConfig
	configEmptyTestCommand := `{
	    "commands": [
			{
				"name": "TestName",
				"command": "TestCommand",
				"typeExpected": "==",
				"expected": "TestExpected"
			}
	    ],
	    "system":
	    {
			` + getConfigSystemForOS() + `
	    }
	  }`

	fileWriter(configEmptyTestCommand, fileNameConfig, false)
	ReadConfig()

	bigAudit = GetBigAudits()[0]
	assert.Nil(t, Call("ls"))
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[INFO] : TestName separated in: ls", nil)

	assert.EqualError(t, Call("ls | ls"), "ls failed")
	CheckFileContent(t, pathAudit, "", doubleLSArray)

	assert.EqualError(t, Call("InvalidCommand"), "Could not find command: InvalidCommand")

	deleteOutput()
}

func TestPrintToLog(t *testing.T) {

	assert.Panics(t, func() { PrintToLog(testText, "WRONGTYPE") })
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[ERROR]: WRONGTYPE is not a valid Log-Type, please use a correct Log-Type", nil)

	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, "[ERROR]: WRONGTYPE is not a valid Log-Type, please use a correct Log-Type", nil)

	assert.True(t, PrintToLog(testText, "INFO"))

	CheckFileContent(t, pathAudit, "[INFO] : TestTestTest123!", nil)

	deleteOutput()
}

func TestCallCompare(t *testing.T) {

	a, b := CallCompare("InvalidCommand", "TestExpected")
	assert.False(t, a)
	assert.EqualError(t, b, "Could not find command: InvalidCommand")

	a, b = CallCompare("echo \"\"", "TestExpected")
	assert.False(t, a)
	assert.Nil(t, b)

	a, b = CallCompare("echo \"\"", "")
	assert.True(t, a)
	assert.Nil(t, b)

	deleteOutput()
}

func TestCallContains(t *testing.T) {

	a, b := CallCompare("InvalidCommand", "TestExpected")
	assert.False(t, a)
	assert.EqualError(t, b, "Could not find command: InvalidCommand")

	a, b = CallCompare("echo \"\"", "TestExpected")
	assert.False(t, a)
	assert.Nil(t, b)

	a, b = CallCompare("echo \"\"", "")
	assert.True(t, a)
	assert.Nil(t, b)

	deleteOutput()
}
