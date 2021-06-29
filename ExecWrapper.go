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
	"os/exec"
	"runtime"
	"strings"
)

var SupportedCommands []string
var commandInSteps string

// allowed commands through call and call compare
func init() {
	SupportedCommands = []string{"Get-MpComputerStatus", "Get-ItemPropertyValue", "Get-NetFirewallProfile", "Test-Path",
		"type", "echo", "reg", "findstr", "dir", "ls", "cat", "grep", "find", "useradd", "stat", "mount", "systemctl",
		"egrep", "test", "call", "Select-String", "%", "modprobe", "df", "rpm", "zypper", "crontab", "stat", "sysctl",
		"journalctl", "apparmor_status", "timedatectl", "ss", "lsof", "iw", "ip", "lsmod", "firewall-cmd", "nmcli",
		"iptables", "ip6tables", "nft", "auditctl", "sshd", "useradd", "rmmod", "awk", "xargs", "subscription-manager",
		"dnf", "sestatus", "ps", "authselect"}

}

// add new commands via -add
func getAdditionalCommands() {
	additionalCommands := flags.addedCommands
	SupportedCommands = appendUnique(SupportedCommands, additionalCommands)
}

func appendUnique(slice []string, newElements []string) []string {
	stringMap := make(map[string]bool)
	allElements := append(slice, newElements...)

	for _, element := range allElements {
		element = strings.ReplaceAll(element, " ", "")
		if element != "" {
			stringMap[strings.ToLower(element)] = true
		}
	}

	allUniqueElements := make([]string, 0)
	for element := range stringMap {
		allUniqueElements = append(allUniqueElements, element)
	}
	return allUniqueElements
}

// checks if command is allowed and executes in small steps
func WrapperForAll(command string, args ...string) *exec.Cmd {

	if len(args) == 0 {
		commandInSteps += command
	} else if commandInSteps == "" {
		commandInSteps += command + " " + strings.Join(args, " ")
	} else {
		commandInSteps += " | " + command + " " + strings.Join(args, " ")
	}

	if runtime.GOOS == "windows" {
		cmd := []string{GetSystem().Argument, commandInSteps}
		return exec.Command(GetSystem().Shell, cmd...)
	} else {
		return exec.Command(GetSystem().Shell, GetSystem().Argument, commandInSteps)
	}
}

// gets fail position and execute errors
func AuditWrapper(audit ...SmallAudit) ([]*exec.Cmd, int, error) {
	var ok bool
	var finished []*exec.Cmd
	commandInSteps = ""

	for i, v := range audit {
		ok = false
		com := strings.ToLower(v.Command)
		for _, supportedCommand := range SupportedCommands {
			if strings.EqualFold(com, supportedCommand) {
				ok = true
				if v.Filepath != "" {
					v.Arguments = append(v.Arguments, v.Filepath)
				}
				foundCommand := []*exec.Cmd{WrapperForAll(com, v.Arguments...)}
				finished = append(finished, foundCommand...)
				break
			}
		}
		if !ok {
			errStr := "Could not find command: " + v.Command
			failPosition := i
			return finished, failPosition, errors.New(errStr)
		}
	}
	return finished, -1, nil
}
