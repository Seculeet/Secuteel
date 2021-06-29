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
	"strings"
)

type SmallAudit struct {
	Name       string   `json:"name"`
	Command    string   `json:"command"`
	Arguments  []string `json:"arguments"`
	Filepath   string   `json:"filepath"`
	Classified bool     `json:"classified"`
}

// splits the command attribute of struct in smaller commands to execute for better debugging
func separateInSmallAuditsButOnlyForCall(bigAuditName string, auditStep string) []SmallAudit {

	lowerBoarder := 0
	counter := 1
	bigCommand := strings.Fields(auditStep)

	var audits = make([]SmallAudit, 0)
	//count how many small audit steps there are
	for _, v := range bigCommand {
		if v == "|" {
			counter++
		}
	}
	if debugModeEnabled {
		WriteDebugLog(bigAuditName+": "+fmt.Sprint(counter)+" small audits detected", "INFO")
	}
	// make slice with len of small audits
	var stringOfAudits = make([]string, counter)
	// separate small audits in string slice
	for _, v := range bigCommand {
		if v != "|" {
			stringOfAudits[lowerBoarder] += v + " "
		} else {
			lowerBoarder++
		}
	}

	// create audit structs from string slices and add name and expected value from big audit step
	for _, v := range stringOfAudits {
		audit := makeAudit(strings.Fields(v))
		audit.Name = bigAuditName
		audits = append(audits, audit)
	}

	return audits

}

// creates audit type from command
func makeAudit(command []string) SmallAudit {
	audit := SmallAudit{}
	// command ALWAYS has to be the first word!
	audit.Command = command[0]

	if debugModeEnabled {
		WriteDebugLog("CommandType detected: "+audit.Command, "INFO")
	}

	argumentsString := make([]string, 0)
	// if filepath is given it HAS to be the last word in command
	for counter, value := range command {
		// indicates filepath with §file§
		if strings.HasPrefix(value, "§file§") {
			audit.Filepath = value[8:]
			if debugModeEnabled {
				WriteDebugLog(audit.Command+": §file§ detected "+audit.Filepath, "INFO")
			}
			for i := 1; i < counter; i++ {
				argumentsString = append(argumentsString, command[i])
			}
			audit.Arguments = argumentsString
			if debugModeEnabled {
				WriteDebugLog(audit.Command+": argumentes added "+fmt.Sprint(audit.Arguments), "INFO")
			}
			return audit
		}
	}
	// command does not have a filepath -> just separating arguments
	for i := 1; i < len(command); i++ {
		argumentsString = append(argumentsString, command[i])
	}
	audit.Arguments = argumentsString
	if debugModeEnabled {
		WriteDebugLog(audit.Command+": argumentes added "+fmt.Sprint(audit.Arguments), "INFO")
	}
	return audit
}
