package main

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceRegex(t *testing.T) {
	deleteOutput()

	regxArr := []string{"Test", "KeinTest", "Test", "", "NichtLeer"}
	outArr := []string{"Test", "Test", "123", "NichtLeer", ""}
	expectedArr := []string{"REDACTED", "Test", "123", "NichtLeer", ""}
	for v := range expectedArr {
		assert.Equal(t, expectedArr[v], replaceRegex(regxArr[v], outArr[v]), "return value has to be equal to \"expectedArr\" instance")
	}
}

func TestSaveAuditShell(t *testing.T) {
	pathTest := "./output/artefacts/test.txt"

	assert.NoError(t, saveAuditFileFromShell(testText, "test"))
	CheckFileExists(t, pathTest)
	CheckFileContent(t, pathTest, "TestTestTest123!", nil)

	os.Chmod(pathTest, 0000)
	assert.Error(t, saveAuditFileFromShell(testText, "test"))

	deleteOutput()
}

func TestSaveAuditCommand(t *testing.T) {
	fileNameConfig := "configWinSystemdetails.json"
	flags.input = "./output/" + fileNameConfig
	configWinTestCommands := `{
        "commands": [
        ],
        "system": 
        {
			` + getConfigSystemForOS() + `
        }
      }`

	smallAuditCommand := SmallAudit{
		Name:       "TestName",
		Command:    "ls",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditWrongCommand := SmallAudit{
		Name:       "TestName",
		Command:    "Notls",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditFile := SmallAudit{
		Name:       "TestName",
		Command:    "ls",
		Arguments:  []string{},
		Filepath:   "TestFilepath",
		Classified: false,
	}
	fileWriter(configWinTestCommands, fileNameConfig, false)
	ReadConfig()
	deleteOutput()

	assert.EqualError(t, saveAuditFileFromCommand(smallAuditWrongCommand), "Could not find command: Notls")
	assert.Error(t, saveAuditFileFromCommand(smallAuditFile), "exit status")
	assert.NoFileExists(t, pathAudit)
	assert.NoFileExists(t, "./output/artefacts/TestName.txt")

	assert.Nil(t, saveAuditFileFromCommand(smallAuditCommand))
	CheckFileExists(t, pathAudit)
	CheckFileExists(t, "./output/artefacts/TestName.txt")

	deleteOutput()
}

func TestSaveArtefact(t *testing.T) {
	if !strings.EqualFold(runtime.GOOS, "windows") {
		return
	}

	smallAuditFile := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   ".\\input\\artefacts\\TestArtefact",
		Classified: false,
	}
	smallAuditNoFile := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditInvalidFile := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   "WrongFilepath",
		Classified: false,
	}
	os.Mkdir("./input", 0777)
	os.Mkdir("./input/artefacts", 0777)
	os.Create("./input/artefacts/TestArtefact")
	os.Mkdir("./output", 0777)
	os.Mkdir("./output/artefacts", 0777)

	saveArtefact(smallAuditFile)
	CheckFileExists(t, pathAudit)
	assert.NoFileExists(t, pathError, "Expects file "+pathError)
	CheckFileContent(t, pathAudit, "[INFO]", nil)
	deleteOutput()

	saveArtefact(smallAuditNoFile)
	CheckFileExists(t, pathAudit)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathAudit, "[ERROR]: TestName cannot save artefact:", nil)
	CheckFileContent(t, pathError, "[ERROR]: TestName cannot save artefact:", nil)
	deleteOutput()

	saveArtefact(smallAuditInvalidFile)
	CheckFileExists(t, pathAudit)
	CheckFileExists(t, pathError)
	CheckFileContent(t, pathAudit, "[ERROR]", nil)
	CheckFileContent(t, pathError, "[ERROR]", nil)

	deleteFile("./input/artefacts/TestArtefact")
	os.RemoveAll("./input/artefacts")
	deleteOutput()
}

func TestCopy(t *testing.T) {
	if !strings.EqualFold(runtime.GOOS, "windows") {
		return
	}

	os.Mkdir("./output", 0777)
	os.Mkdir("./output/artefacts", 0777)
	os.Create("./output/artefacts/TestArtefact")
	a, b := copyArtefact(".\\output\\artefacts\\TestArtefact", "TestString")
	assert.Zero(t, a)
	assert.Nil(t, b)

	a, b = copyArtefact(".\\output\\artefacts\\TestArtefact", ".\\output\\artefacts\\TestArtefact")
	assert.Zero(t, a)
	assert.Nil(t, b)

	os.Chmod(".\\output\\artefacts\\TestArtefact", 0000)

	a, b = copyArtefact(".\\output\\artefacts\\TestArtefact", ".\\output\\artefacts\\TestArtefact")

	assert.Zero(t, a)
	assert.Error(t, b)

	os.Chmod(".\\output\\artefacts\\TestArtefact", 0777)

	a, b = copyArtefact("WrongPath", "TestString")
	assert.Zero(t, a)
	assert.Error(t, b)

	a, b = copyArtefact("", "")
	assert.Zero(t, a)
	assert.Error(t, b)

	a, b = copyArtefact(".\\output", ".\\output")
	assert.Zero(t, a)
	assert.Error(t, b)

	deleteFile("./output/artefacts/TestArtefact")
	deleteFile("./TestString")
	deleteOutput()
}
