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
	"testing"

	"github.com/stretchr/testify/assert"
)

var testArgumentArray []string

func init() {
	flags.input = "configWin"
	testArgumentArray = []string{"TestArgument1", "TestArgument2"}
}

func TestMakeAudit(t *testing.T) {
	smallAuditNoFile := SmallAudit{
		Name:       "",
		Command:    "TestCommand",
		Arguments:  testArgumentArray,
		Filepath:   "",
		Classified: false,
	}
	smallAuditFile := SmallAudit{
		Name:       "",
		Command:    "TestCommand",
		Arguments:  testArgumentArray,
		Filepath:   "TestFile",
		Classified: false,
	}
	smallAuditEmpty := SmallAudit{
		Name:       "",
		Command:    "",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditCommandOnly := SmallAudit{
		Name:       "",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditWrongPrefix := SmallAudit{
		Name:       "",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   "TestFile",
		Classified: false,
	}

	testNoPrefix := makeAudit([]string{"TestCommand", "TestArgument1", "TestArgument2"})
	testPrefix := makeAudit([]string{"TestCommand", "TestArgument1", "TestArgument2", "§file§TestFile"})
	testEmpty := makeAudit([]string{""})
	testCommandOnly := makeAudit([]string{"TestCommand"})
	testWrongPrefix := makeAudit([]string{"TestCommand", "§file§TestFile", "TestArgument1", "TestArgument2"})

	assert.Equal(t, testNoPrefix, smallAuditNoFile)
	assert.Equal(t, testPrefix, smallAuditFile)
	assert.Equal(t, testEmpty, smallAuditEmpty)
	assert.Equal(t, testCommandOnly, smallAuditCommandOnly)
	assert.Equal(t, testWrongPrefix, smallAuditWrongPrefix)
}

func TestSeperate(t *testing.T) {

	smallAuditName := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  []string{},
		Filepath:   "",
		Classified: false,
	}
	smallAuditArguments := SmallAudit{
		Name:       "TestName",
		Command:    "TestCommand",
		Arguments:  testArgumentArray,
		Filepath:   "",
		Classified: false,
	}

	testCommand := separateInSmallAuditsButOnlyForCall("TestName", "TestCommand")
	testArguments := separateInSmallAuditsButOnlyForCall("TestName", "TestCommand TestArgument1  TestArgument2")
	testMultiple := separateInSmallAuditsButOnlyForCall("TestName", "TestCommand TestArgument1 TestArgument2 | TestCommand TestArgument1 TestArgument2")

	assert.Equal(t, testCommand, []SmallAudit{smallAuditName})
	assert.Equal(t, testArguments, []SmallAudit{smallAuditArguments})
	assert.Equal(t, testMultiple, []SmallAudit{smallAuditArguments, smallAuditArguments})
	assert.Panics(t, func() { separateInSmallAuditsButOnlyForCall("", "") })
	assert.Panics(t, func() { separateInSmallAuditsButOnlyForCall("TestName", "") })
	assert.Panics(t, func() { separateInSmallAuditsButOnlyForCall("TestName", "| TestArgument1 TestArgument2") })
	assert.Panics(t, func() { separateInSmallAuditsButOnlyForCall("TestName", "TestCommand TestArgument1 TestArgument2 |") })
}
