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
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAdditionalCommands(t *testing.T) {
	newCmd := "new-command"
	secondCmd := "second"
	emptyCmd := " "
	flags.addedCommands = []string{newCmd, strings.ToUpper(newCmd), secondCmd, strings.ToLower(newCmd), emptyCmd}
	getAdditionalCommands()

	var isNewCmdInSupportedCommands bool
	var isSecondCmdInSupportedCommands bool
	var isEmptyCmdInSupportedCommands bool

	countNewCmd := 0
	expectedCountNewCmd := 1

	countSecondCmd := 0
	expectedCountSecondCmd := 1

	for _, supportedCommand := range SupportedCommands {
		if strings.EqualFold(newCmd, supportedCommand) {
			isNewCmdInSupportedCommands = true
			countNewCmd++
		}

		if strings.EqualFold(secondCmd, supportedCommand) {
			isSecondCmdInSupportedCommands = true
			countSecondCmd++
		}

		if strings.EqualFold(emptyCmd, supportedCommand) {
			isEmptyCmdInSupportedCommands = true
		}
	}
	assert.True(t, isNewCmdInSupportedCommands, "Check if \""+newCmd+"\" was added to SupportedCommands")
	assert.Equal(t, expectedCountNewCmd, countNewCmd, "Check if \""+newCmd+"\" is included only once")

	assert.True(t, isSecondCmdInSupportedCommands, "Check if \""+secondCmd+"\" was added to SupportedCommands")
	assert.Equal(t, expectedCountSecondCmd, countSecondCmd, "Check if \""+secondCmd+"\" is included only once")

	assert.False(t, isEmptyCmdInSupportedCommands, "Check if \""+emptyCmd+"\" was not added to SupportedCommands")
}

func TestAppendUnique(t *testing.T) {
	var testCases = []struct {
		sliceOne      []string
		sliceTwo      []string
		expectedSlice []string
	}{
		{
			[]string{"A", "B"},
			[]string{"C", "D"},
			[]string{"a", "b", "c", "d"},
		},
		{
			[]string{"a", "", "b"},
			[]string{"c", "d", "     "},
			[]string{"a", "b", "c", "d"},
		},
		{
			[]string{"a", "a"},
			[]string{"b", "b"},
			[]string{"a", "b"},
		},
		{
			[]string{"a", "b"},
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			[]string{"a", "a"},
			[]string{"a", "a"},
			[]string{"a"},
		},
		{
			[]string{},
			[]string{"a"},
			[]string{"a"},
		},
		{
			[]string{"a"},
			[]string{},
			[]string{"a"},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
	}
	for _, test := range testCases {
		newSlice := appendUnique(test.sliceOne, test.sliceTwo)
		sort.Strings(newSlice)
		sort.Strings(test.expectedSlice)
		assert.Equal(t, test.expectedSlice, newSlice, "Check \""+strings.Join(newSlice, " ")+"\" is equal to \""+strings.Join(test.expectedSlice, " ")+"\"")
	}
}

func TestWrapperForAll(t *testing.T) {
	// auditStart := exec.Command("powershell", []string{"/C", "get-mpcomputerstatus"}...)
	// auditEnd := exec.Command("powershell", []string{"/C", "findstr", "AntispywareEnabled"}...)
	// audits := []*exec.Cmd{auditStart, auditEnd}
	// WrapperForAll(command string, args ...string)
}
