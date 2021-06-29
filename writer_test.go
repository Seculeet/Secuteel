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

var pathOutput string
var pathAudit string
var pathError string
var pathResult string
var testText string
var infoLogType string

var testBigAuditFull BigAudit
var testBigAuditNoName BigAudit
var testBigAuditNoCommand BigAudit
var testBigAuditEmpty BigAudit

func init() {
	flags.input = "configWin"
	pathOutput = "./output/"
	pathAudit = "./output/audit.log"
	pathError = "./output/error.log"
	pathResult = "./output/result.json"
	testText = "TestTestTest123!"
	infoLogType = "INFO"

	testBigAuditFull = BigAudit{
		Name:             "TestName",
		Command:          "TestCommand",
		DontSaveArtefact: true,
		TypeExpected:     "TestTypeExpected",
		Expected:         "TestExpected",
		Desc:             "TestDesc",
	}
	testBigAuditNoName = BigAudit{
		Name:             "",
		Command:          "TestCommand",
		DontSaveArtefact: true,
		TypeExpected:     "TestTypeExpected",
		Expected:         "TestExpected",
		Desc:             "TestDesc",
	}
	testBigAuditNoCommand = BigAudit{
		Name:             "TestName",
		Command:          "",
		DontSaveArtefact: true,
		TypeExpected:     "TestTypeExpected",
		Expected:         "TestExpected",
		Desc:             "TestDesc",
	}
	testBigAuditEmpty = BigAudit{
		Name:             "",
		Command:          "",
		DontSaveArtefact: true,
		TypeExpected:     "TestTypeExpected",
		Expected:         "TestExpected",
		Desc:             "TestDesc",
	}
}

func TestCreateFolderIfNotExist(t *testing.T) {
	createFolderIfNotExist(pathOutput, "output")
	assert.DirExists(t, pathOutput, "Create output folder")
	assert.NoFileExists(t, pathAudit)

	createFolderIfNotExist("./output/artefacts", "artefacts")

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "[INFO] : Create artefacts folder", nil)

	assert.Panics(t, func() { createFolderIfNotExist("./invalid/invalid", "invalid") }, "Create \"./invalid/invalid\" folder has to return error")
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathAudit, "Can't create invalid folder.", nil)
	CheckFileContent(t, pathError, "Can't create invalid folder.", nil)

	deleteOutput()
}

func TestCheckPathExists(t *testing.T) {
	os.Mkdir(pathOutput, 0777)
	assert.True(t, checkPathExists(pathOutput), "Output path has to return true")
	assert.False(t, checkPathExists("InvalidPath"), "Invalid paths have to return false")

	deleteOutput()
}

func TestFileWriter(t *testing.T) {
	FirstAuditEntry = false
	testArray := []string{"[INFO] : Create output folder", testText}

	assert.NoError(t, fileWriter(testText, "audit.log", true), "fileWriter should not return error")

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", testArray)

	os.Chmod("./output/audit.log", 0000)
	assert.Error(t, fileWriter(testText, "audit.log", true))

	deleteOutput()
}

func TestFileWriterAppend(t *testing.T) {
	appendedText := ", appendedText\n"
	expected := testText + appendedText

	fileWriter(testText, "audit.log", false)
	fileWriter(appendedText, "audit.log", true)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expected, nil)

	deleteOutput()
}

func TestIsValidLogTypeDefaultCase(t *testing.T) {

	expected := "[ERROR]: WRONGTYPE is not a valid Log-Type, please use a correct Log-Type"

	var testCases = []struct {
		logType string
	}{
		{
			"INFO",
		},
		{
			"DEBUG",
		},
		{
			"WARN",
		},
		{
			"FAIL",
		},
		{
			"ERROR",
		},
	}
	for _, test := range testCases {
		assert.Equal(t, isValidLogType(test.logType), test.logType, "\""+test.logType+"\""+" needs to return "+"\""+test.logType+"\"")
	}

	assert.Panics(t, func() { isValidLogType("WRONGTYPE") }, "Expected panic")
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expected, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expected, nil)

	deleteOutput()
}

func TestWriteLog(t *testing.T) {
	expected := "[INFO] : TestTestTest123!"

	WriteLog(testText, infoLogType)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expected, nil)

	deleteOutput()
}

func TestWriteErrorLog(t *testing.T) {
	expected := "[ERROR]: TestTestTest123!"

	WriteErrorLog(testText, "")

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expected, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expected, nil)

	deleteOutput()
}

func TestWriteErrorLogAudit(t *testing.T) {
	testAuditName := "TEST321"

	expected := "[ERROR]: TEST321 TestTestTest123!"

	WriteErrorLog(testText, testAuditName)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, expected, nil)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathError, expected, nil)

	deleteOutput()
}

func TestWriteResultJSON(t *testing.T) {

	FirstResultEntry = false
	os.Mkdir("./output", 0777)

	testOutput := "TestOutput"
	testError := "TestError"
	testOperator := "TestOperator"

	expected := `{
	"input/configTest": [
		{
			"Name": "TestName",
			"Command": "TestCommand",
			"Command was executed": false,
			"Error-Message": "TestError"
		},
		{
			"Name": "TestName",
			"Command": "TestCommand",
			"Command was executed": true,
			"Output is as expected": false,
			"Expected Value": "TestExpected",
			"Actual Value": "TestOutput",
			"Operator": "TestOperator"
		},
		{
			"Name": "TestName",
			"Command": "TestCommand",
			"Command was executed": true,
			"Output is as expected": true
		}
	]
}`

	var testResultCases = []struct {
		commandB bool
		auditB   bool
	}{
		{false, false},
		{true, false},
		{true, true},
	}
	ConfigName = "input/configTest"

	for _, test := range testResultCases {
		WriteResultJSON(testBigAuditFull, test.commandB, test.auditB, testOutput, testError, testOperator)
	}

	CheckFileExists(t, pathResult)
	CheckFileContent(t, pathResult, expected, nil)

	deleteOutput()
}

func TestWriteCommandSuccessLog(t *testing.T) {
	expected := []string{"[INFO] : Name: TestName",
		"TestName: TestCommand executed",
		"[INFO] : Name: TestName",
		"[INFO] : "}

	var testBigAuditArray = []struct {
		testBigAudit BigAudit
	}{
		{testBigAuditFull},
		{testBigAuditNoName},
		{testBigAuditNoCommand},
		{testBigAuditEmpty},
	}

	for _, test := range testBigAuditArray {
		WriteCommandSuccessLog(test.testBigAudit)
	}
	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", expected)

	deleteOutput()
}

func TestWriteCommandFailedLogError(t *testing.T) {
	expected := []string{"[FAIL] : Name: TestName",
		"TestName: TestCommand failed"}

	WriteCommandFailedLog(testBigAuditFull, errors.New("TestError"))

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", expected)

	deleteOutput()
}

func TestIsValidLogTypeUpperCase(t *testing.T) {
	smallLogType := "info"
	assert.Equal(t, isValidLogType(smallLogType), infoLogType, "\"info\" has to convert to upper-case")
}

func TestWriteAuditFailedLog(t *testing.T) {
	expected := []string{"[FAIL] : Name: TestName",
		"[INFO] : TestName: TestCommand testArgument1 testArgument2",
		"[FAIL] : TestName: TestCommand testArgument1 testArgument2",
		"[WARN] : TestName: TestCommand testArgument1 testArgument2"}

	testFailPosition := 1
	testArgumentArray := []string{"testArgument1", "testArgument2"}

	testSmallAudit := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  testArgumentArray,
		Filepath:   "TestFilepath",
		Classified: true,
	}

	testSmallAuditArray := []SmallAudit{testSmallAudit, testSmallAudit, testSmallAudit}

	WriteAuditFailedLog(testSmallAuditArray, testFailPosition)

	CheckFileExists(t, pathAudit)
	CheckFileContent(t, pathAudit, "", expected)

	deleteOutput()
}

func TestUnescapeUnicodeCharactersInJSON(t *testing.T) {
	var testCases = []struct {
		input    []byte
		expected []byte
	}{
		{
			[]byte("Hello World"),
			[]byte("Hello World"),
		},
		{
			[]byte(""),
			[]byte(""),
		},
		{
			[]byte("\u003chtml\u003e"),
			[]byte("<html>"),
		},
		{
			[]byte("\\\\u"),
			[]byte("\\\\u"),
		},
	}
	for _, test := range testCases {
		unescapedContent, err := UnescapeUnicodeCharactersInJSON(test.input)
		inputAsString := string(unescapedContent)
		expectedAsString := string(test.expected)
		assert.Equal(t, expectedAsString, inputAsString, "Check \""+inputAsString+"\" is equal to \""+expectedAsString+"\"")
		assert.Nil(t, err, "Check \"UnescapeUnicodeCharactersInJSON()\" throws no error")
	}
	deleteFile("./TestString")
	deleteOutput()
}

func TestDeleteFile(t *testing.T) {
	os.Mkdir(pathOutput, 0777)
	os.Create("./output/TestFile")
	deleteFile("./output/TestFile")
	assert.NoFileExists(t, "./output/TestFile")
	deleteFile("./output/NotTestFile")

	deleteOutput()
}

func TestDeleteOutput(t *testing.T) {
	os.Mkdir(pathOutput, 0777)
	deleteOutput()
	assert.NoDirExists(t, "./output")
}

func TestReplaceAllWhitespaces(t *testing.T) {
	input1 := "\n\t\t\tTestText"
	input2 := "\n\t\t\tTestText\n\t\t\t"
	expected := "TestText"
	assert.Equal(t, replaceAllWhitespace(input1), expected)
	assert.Equal(t, replaceAllWhitespace(input2), expected)
}

func CheckFileExists(t assert.TestingT, path string) {
	assert.FileExists(t, path, "Expects file "+path)
}

func CheckFileContent(t assert.TestingT, path string, expectedStr string, expectedArr []string) {
	ContentAsBytes, _ := os.ReadFile(path)
	if len(expectedArr) == 0 {
		assert.Contains(t, replaceAllWhitespace(string(ContentAsBytes)), replaceAllWhitespace(expectedStr), "File has to contain expected String")
	} else {
		for _, v := range expectedArr {
			assert.Contains(t, replaceAllWhitespace(string(ContentAsBytes)), replaceAllWhitespace(v), "File has to contain all expected Strings in Array")
		}
	}
}
