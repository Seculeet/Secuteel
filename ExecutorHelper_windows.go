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
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"golang.org/x/sys/windows/registry"
)

// function to query registry for value
func getRegistry(key registry.Key, path, value string) (string, error) {
	regKey, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return "", errors.New("registry not found")
	}
	defer regKey.Close()

	regStringValue, _, err := regKey.GetStringValue(value)
	if err != nil {
		regIntValue, _, err := regKey.GetIntegerValue(value)
		if err != nil {
			return "", errors.New("value not found: " + value)
		}
		regStringValue = fmt.Sprint(regIntValue)
	}
	return regStringValue, nil
}

// parse input to usable format
func parseRegistry(fullReg string) (registry.Key, string) {

	splitReg := strings.SplitN(fullReg, "\\", 2)
	key := parseKey(splitReg[0])
	path := splitReg[1]
	return key, path

}

// return known registries
func parseKey(key string) registry.Key {
	switch strings.ToUpper(key) {
	case "HKLM:":
		return registry.LOCAL_MACHINE
	case "HKCU:":
		return registry.CURRENT_USER
	case "HKCR:":
		return registry.CLASSES_ROOT
	case "HKU:":
		return registry.USERS
	case "HKCC:":
		return registry.CURRENT_CONFIG
	case "HKPD:":
		return registry.PERFORMANCE_DATA
	default:
		return registry.CURRENT_USER
	}
}

// available through JS
func RegQuery(registryPath string, value string) error {
	output = ""
	key, path := parseRegistry(registryPath)
	result, err := getRegistry(key, path, value)
	if err != nil {
		return err
	}
	output = result
	copyErr := saveAuditFileFromShell(result, bigAudit.Name)
	if copyErr != nil {
		if debugModeEnabled {
			WriteDebugLog(bigAudit.Name+" cannot save artefact: "+copyErr.Error(), "ERROR")
		}
		WriteErrorLog("cannot save artefact: "+copyErr.Error(), bigAudit.Name)
		return err
	}
	if debugModeEnabled {
		WriteDebugLog(bigAudit.Name+" saved artefact", "INFO")
	}
	return nil
}

// registers Windows commands
func createCommandVM() {
	VmCommand = goja.New()

	err := VmCommand.Set("call", Call)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
	err = VmCommand.Set("printToLog", PrintToLog)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
	err = VmCommand.Set("callCompare", CallCompare)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
	err = VmCommand.Set("callContains", CallContains)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
	err = VmCommand.Set("printToConsole", PrintToConsole)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
	err = VmCommand.Set("shell", Shell)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}

	err = VmCommand.Set("regQuery", RegQuery)
	if err != nil {
		WriteErrorLog(err.Error(), "")
	}
}
