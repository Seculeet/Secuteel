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
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFilesInOutput(t *testing.T) {
	if !strings.EqualFold(runtime.GOOS, "windows") {
		return
	}

	auditArr := []string{"output\\audit.log"}
	logArr := []string{"output\\audit.log", "output\\error.log"}
	allArr := []string{"output\\audit.log", "output\\error.log", "output\\result.json"}
	os.Mkdir("./output", 0777)
	os.Create("./output/audit.log")
	assert.Equal(t, GetAllFilesInOutput(), auditArr)

	os.Create("./output/error.log")
	assert.Equal(t, GetAllFilesInOutput(), logArr)

	os.Create("./output/result.json")
	assert.Equal(t, GetAllFilesInOutput(), allArr)

	deleteOutput()
	assert.Panics(t, func() { GetAllFilesInOutput() })
}

func TestZipFiles(t *testing.T) {
	os.Mkdir("./output", 0777)
	os.Create("./output/audit.log")
	assert.Nil(t, ZipFiles("TestString", []string{}))
	assert.Nil(t, ZipFiles("TestString", []string{"TestString"}))
	assert.Nil(t, ZipFiles("./output/audit.log", []string{}))
	assert.Nil(t, ZipFiles("./output/audit.log", []string{"./output/audit.log"}))
	assert.Nil(t, ZipFiles("TestString", []string{"./output/audit.log"}))

	assert.Error(t, ZipFiles("", []string{}))
	assert.Error(t, ZipFiles("", []string{"./output/audit.log"}))
	assert.Error(t, ZipFiles("./output/audit.log", []string{"Test"}))

	assert.FileExists(t, "./TestString")

	deleteFile("./TestString")
	deleteOutput()
}

func TestAddFileToZip(t *testing.T) {
	os.Create("./TestString")
	assert.Error(t, AddFileToZip(nil, ""))
	assert.Panics(t, func() { AddFileToZip(nil, "TestString") })

	deleteFile("./TestString")
	assert.Error(t, AddFileToZip(nil, "TestString"))

	deleteOutput()
}
