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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// copies input to artefact file
func copyArtefact(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		err := src + " is not a regular file"
		return 0, fmt.Errorf(err)
	}

	source, err := ioutil.ReadFile(src)
	if err != nil {
		return 0, err
	}

	if bigAudit.BlackenContent != "" {
		source = []byte(replaceRegex(bigAudit.BlackenContent, string(source)))
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	r := strings.NewReader(string(source))
	nBytes, err := io.Copy(destination, r)
	return nBytes, err
}

// replaces regex in string
func replaceRegex(regex string, output string) string {
	if regex != "" {
		re := regexp.MustCompile(regex)
		s := re.ReplaceAllString(output, "REDACTED")
		if debugModeEnabled {
			WriteDebugLog("regex detected: "+re.String(), "INFO")
		}

		return s
	} else {
		if debugModeEnabled {
			WriteDebugLog("no regex detected", "INFO")
		}
		return output
	}
}

// copies shell to get artefact from direct execute
func saveAuditFileFromShell(shellOutput string, name string) error {

	createFolderIfNotExist("./output/", "output")
	createFolderIfNotExist("./output/artefacts/", "artefacts")

	fileName := "artefacts/" + name + ".txt"

	return fileWriter(shellOutput, fileName, true)
}

// in case there is no filepath given the first command gets executed and saved as a txt
func saveAuditFileFromCommand(small SmallAudit) error {

	execCommand, _, err := AuditWrapper(small)
	if err != nil {
		return err //errors.New("can't save Artefact")
	}

	artefact, err := execCommand[0].Output()
	if err != nil {
		return err //errors.New("exit status 1, empty bytes")
	}

	createFolderIfNotExist("./output/", "output")
	createFolderIfNotExist("./output/artefacts/", "artefacts")

	fileName := "artefacts/" + small.Name + ".txt"

	if bigAudit.BlackenContent != "" {
		artefact = []byte(replaceRegex(bigAudit.BlackenContent, string(artefact)))
	}
	return fileWriter(string(artefact), fileName, true)
}

// can be called from Windows or Linux to save artefact, must be called with the first smallAudit
func saveArtefact(audit SmallAudit) {
	if audit.Filepath != "" {
		s := strings.Split(audit.Filepath, "\\")
		fileEnd := s[len(s)-1]
		_, err := copyArtefact(audit.Filepath, "./output/artefacts/"+fileEnd)

		if err != nil {
			if debugModeEnabled {
				WriteDebugLog("cannot save artefact "+err.Error(), "ERROR")
			}
			WriteErrorLog(err.Error(), "")
		} else {
			if debugModeEnabled {
				WriteDebugLog(audit.Name+" artefact successfully saved", "INFO")
			}
			WriteLog("artefact "+audit.Name+" successfully saved", "INFO")
		}
	} else {
		err := saveAuditFileFromCommand(audit)

		if err != nil {
			WriteErrorLog("cannot save artefact", audit.Name)
			if debugModeEnabled {
				WriteDebugLog(audit.Name+" cannot save artefact", "ERROR")
			}
		} else {
			WriteLog("artefact "+audit.Name+" successfully saved", "INFO")
		}
	}
}
