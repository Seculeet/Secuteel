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

func TestMainWihtoutConfigFlag(t *testing.T) {

	expectedErrorLogContent := "[ERROR]: you need to specify an input file"
	expectedAuditLogContent := []string{"[INFO] : Create output folder", expectedErrorLogContent}

	flags.input = ""

	assert.Panics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")
	CheckFileContent(t, pathError, expectedErrorLogContent, nil)

	assert.NoFileExists(t, pathResult, "Check \""+pathResult+"\" file doesn't exists")

	deleteOutput()
}

func TestMainValidConfigFlagValidJSONValidCommandsValidSystem(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "check_Browser_is_disabled",
				"command": "echo hallo",
				"typeExpected": "==",
				"expected": "hallo",
				"description": "Checking Bluetooth Support Service is disabled"
			},
			{
				"name": "check_lfsvc_disabled",
				"command": "echo hallo",
				"typeExpected": "==",
				"expected": "hall",
				"description": "Checking Geolocation Service is disabled",
				"dontSaveArtefact": true
			}
        ],
        "system":
        {
			` + getConfigSystemForOS() + `
        }
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : sanity check passed",
		"[INFO] : check_Browser_is_disabled separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact check_Browser_is_disabled successfully saved",
		"[INFO] : Name: check_Browser_is_disabled",
		"check_Browser_is_disabled: echo hallo executed",
		"[INFO] : check_Browser_is_disabled finished! (Executed + output == expected)",
		"[INFO] : check_lfsvc_disabled separated in: echo",
		"[INFO] : Name: check_lfsvc_disabled",
		"check_lfsvc_disabled: echo hallo executed",
		"[FAIL] : check_lfsvc_disabled output != expected",
		"[INFO] : Zip folder:",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "check_Browser_is_disabled",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "check_lfsvc_disabled",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": false,
				"Expected Value": "hall",
				"Actual Value": "hallo",
				"Operator": "=="
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" run without panics")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")

	pathArtefacts := "./output/artefacts/"
	firstAuditArtefact := pathArtefacts + "check_Browser_is_disabled.txt"
	secondAuditArtefact := pathArtefacts + "check_lfsvc_disabled.txt"

	assert.FileExists(t, firstAuditArtefact, "Check first audit \""+firstAuditArtefact+"\" artefact exists")
	assert.NoFileExists(t, secondAuditArtefact, "Check second audit \""+secondAuditArtefact+"\" artefact doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainValidConfigFlagValidJSONEmptyCommandsValidSystem(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
        ],
        "system":
        {
			` + getConfigSystemForOS() + `
        }
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : sanity check passed",
		"[ERROR]: you have to specify at least one command in your config file",
	}

	expectedErrorLogContent := "[ERROR]: you have to specify at least one command in your config file"

	blackBoxWriter(configContent, configFileName)

	assert.Panics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")
	CheckFileContent(t, pathError, expectedErrorLogContent, nil)

	assert.NoFileExists(t, pathResult, "Check \""+pathResult+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainValidConfigFlagValidJSONValidCommandsEmptySystem(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "check_Browser_is_disabled",
				"command": "Get-ItemPropertyValue Registry::HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\bthserv -Name Start",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking Bluetooth Support Service is disabled"
			}
        ],
        "system":
        {}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[ERROR]: you have to specify the system in your config file",
	}

	expectedErrorLogContent := "[ERROR]: you have to specify the system in your config file"

	blackBoxWriter(configContent, configFileName)

	assert.Panics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")
	CheckFileContent(t, pathError, expectedErrorLogContent, nil)

	assert.NoFileExists(t, pathResult, "Check \""+pathResult+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainSetHelpFlag(t *testing.T) {

	deleteOutput()

	flags.help = true
	assert.Panics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.NoFileExists(t, pathAudit, "Check \""+pathAudit+"\" file doesn't exists")
	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")
	assert.NoFileExists(t, pathResult, "Check \""+pathResult+"\" file doesn't exists")
	flags.help = false
}

func TestMainValidConfigFlagValidJSONValidCommandsSystemWithRoot(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "check_Browser_is_disabled",
				"command": "Get-ItemPropertyValue Registry::HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\bthserv -Name Start",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking Bluetooth Support Service is disabled"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `,
			"root": true
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[ERROR]: you are not root",
	}

	expectedErrorLogContent := "[ERROR]: you are not root"

	blackBoxWriter(configContent, configFileName)

	assert.Panics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")
	CheckFileContent(t, pathError, expectedErrorLogContent, nil)

	assert.NoFileExists(t, pathResult, "Check \""+pathResult+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainValidConfigFlagValidJSONValidCommandsSystemRegQuery(t *testing.T) {
	if !strings.EqualFold(runtime.GOOS, "windows") {
		return
	}

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "get_bluetooth_status_0",
				"command": "regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_1",
				"command": "regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Starttt')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_2",
				"command": "regQuery('HKCU:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_3",
				"command": "regQuery('HKCR:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_4",
				"command": "regQuery('HKU:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_5",
				"command": "regQuery('HKCC:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_6",
				"command": "regQuery('HKPD:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			},
			{
				"name": "get_bluetooth_status_7",
				"command": "regQuery('HKFF:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"typeExpected": "==",
				"expected": "3",
				"description": "Checking bluetooth enabled"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : Name: get_bluetooth_status_0",
		"get_bluetooth_status_0: regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start') executed",
		"[INFO] : get_bluetooth_status_0 finished! (Executed + output == expected)",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "get_bluetooth_status_0",
				"Command": "regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "get_bluetooth_status_1",
				"Command": "regQuery('HKLM:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Starttt')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_2",
				"Command": "regQuery('HKCU:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_3",
				"Command": "regQuery('HKCR:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_4",
				"Command": "regQuery('HKU:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_5",
				"Command": "regQuery('HKCC:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_6",
				"Command": "regQuery('HKPD:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			},
			{
				"Name": "get_bluetooth_status_7",
				"Command": "regQuery('HKFF:\\SYSTEM\\CurrentControlSet\\Services\\bthserv', 'Start')",
				"Command was executed": false,
				"Error-Message": "command not executed"
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainDifferentStringsInCommand(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	halloWithBacktick := "`hallo`"
	configContent := `{
        "commands": [
			{
				"name": "different_string_1",
				"command": "echo hallo",
				"typeExpected": "==",
				"expected": "hallo"
			},
			{
				"name": "different_string_2",
				"command": "echo \"hallo\"",
				"typeExpected": "==",
				"expected": "hallo"
			},
			{
				"name": "different_string_3",
				"command": "echo ` + halloWithBacktick + `",
				"typeExpected": "==",
				"expected": "hallo"
			},
			{
				"name": "different_string_4",
				"command": "echo 'hallo'",
				"typeExpected": "==",
				"expected": "hallo"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : sanity check passed",
		"[INFO] : different_string_1 separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact different_string_1 successfully saved",
		"[INFO] : Name: different_string_1",
		"different_string_1: echo hallo executed",
		"[INFO] : different_string_1 finished! (Executed + output == expected)",
		"[INFO] : different_string_2 separated in: echo",
		"[INFO] : Name: different_string_2",
		"different_string_2: echo \"hallo\" executed",
		"[INFO] : different_string_2 finished! (Executed + output == expected)",
		"[INFO] : different_string_3 separated in: echo",
		"[INFO] : artefact different_string_3 successfully saved",
		"[INFO] : Name: different_string_3",
		"different_string_3: echo `hallo` executed",
		"[INFO] : different_string_3 finished! (Executed + output == expected)",
		"[ERROR]: different_string_4: SyntaxError: JavaScript",
		"[INFO] : Zip folder:",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "different_string_1",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_string_2",
				"Command": "echo \"hallo\"",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_string_3",
				"Command": "echo ` + halloWithBacktick + `",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_string_4",
				"Command": "echo 'hallo'",
				"Command was executed": false,
				"Error-Message": "command not executed"
			}
		]
	}`

	expectedErrorLogContent := "[ERROR]: different_string_4: SyntaxError: JavaScript"

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.FileExists(t, pathError, "Check \""+pathError+"\" file exists")
	CheckFileContent(t, pathError, expectedErrorLogContent, nil)

	pathArtefacts := "./output/artefacts/"
	different_string_1 := pathArtefacts + "different_string_1.txt"
	different_string_2 := pathArtefacts + "different_string_2.txt"
	different_string_3 := pathArtefacts + "different_string_3.txt"
	different_string_4 := pathArtefacts + "different_string_4.txt"

	assert.FileExists(t, different_string_1, "Check \""+different_string_1+"\" artefact exists")
	assert.FileExists(t, different_string_2, "Check \""+different_string_2+"\" artefact exists")
	assert.FileExists(t, different_string_3, "Check \""+different_string_3+"\" artefact exists")
	assert.NoFileExists(t, different_string_4, "Check \""+different_string_4+"\" artefact doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainDifferentOwnJavaScriptFunctions(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "func_callCompare_1",
				"command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callCompare_2",
				"command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"expected": "halloo"
			},
			{
				"name": "func_callContains_1",
				"command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callContains_2",
				"command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"expected": "halloo"
			},
			{
				"name": "func_shell",
				"command": "shell('echo hallo')",
				"expected": "hallo"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : func_callCompare_1 separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact func_callCompare_1 successfully saved",
		"[INFO] : Name: func_callCompare_1",
		"func_callCompare_1: if(callCompare('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_1 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"[INFO] : artefact func_callCompare_2 successfully saved",
		"[INFO] : Name: func_callCompare_2",
		"func_callCompare_2: if(callCompare('echo halloo', 'hallo')) printToConsole('should not print') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_1 separated in: echo",
		"[INFO] : artefact func_callContains_1 successfully saved",
		"[INFO] : Name: func_callContains_1",
		"func_callContains_1: if(callContains('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callContains_1 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_2 separated in: echo",
		"[INFO] : artefact func_callContains_2 successfully saved",
		"[INFO] : Name: func_callContains_2",
		"func_callContains_2: if(callContains('echo halloo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"artefact func_shell successfully saved",
		"[INFO] : Name: func_shell",
		"func_shell: shell('echo hallo') executed",
		"[INFO] : func_shell finished! (Executed + output == expected)",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "func_callCompare_1",
				"Command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callCompare_2",
				"Command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_1",
				"Command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_2",
				"Command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_shell",
				"Command": "shell('echo hallo')",
				"Command was executed": true,
				"Output is as expected": true
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
}

func TestMainDifferentOwnJavaScriptFunctionsWithVerbose(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName
	flags.verbose = true

	configContent := `{
        "commands": [
			{
				"name": "func_callCompare_1",
				"command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callCompare_2",
				"command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"expected": "halloo"
			},
			{
				"name": "func_callContains_1",
				"command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callContains_2",
				"command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"expected": "halloo"
			},
			{
				"name": "func_shell",
				"command": "shell('echo hallo')",
				"expected": "hallo"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : func_callCompare_1 separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact func_callCompare_1 successfully saved",
		"[INFO] : Name: func_callCompare_1",
		"func_callCompare_1: if(callCompare('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_1 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"[INFO] : artefact func_callCompare_2 successfully saved",
		"[INFO] : Name: func_callCompare_2",
		"func_callCompare_2: if(callCompare('echo halloo', 'hallo')) printToConsole('should not print') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_1 separated in: echo",
		"[INFO] : artefact func_callContains_1 successfully saved",
		"[INFO] : Name: func_callContains_1",
		"func_callContains_1: if(callContains('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callContains_1 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_2 separated in: echo",
		"[INFO] : artefact func_callContains_2 successfully saved",
		"[INFO] : Name: func_callContains_2",
		"func_callContains_2: if(callContains('echo halloo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"artefact func_shell successfully saved",
		"[INFO] : Name: func_shell",
		"func_shell: shell('echo hallo') executed",
		"[INFO] : func_shell finished! (Executed + output == expected)",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "func_callCompare_1",
				"Command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callCompare_2",
				"Command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_1",
				"Command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_2",
				"Command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_shell",
				"Command": "shell('echo hallo')",
				"Command was executed": true,
				"Output is as expected": true
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
	flags.verbose = false
}

func TestMainDifferentOwnJavaScriptFunctionsWithDebug(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName
	flags.debug = true
	flags.skipSanity = true

	pathDebug := "./output/debug.log"

	configContent := `{
        "commands": [
			{
				"name": "func_callCompare_1",
				"command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callCompare_2",
				"command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"expected": "halloo"
			},
			{
				"name": "func_callContains_1",
				"command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"expected": "hallo"
			},
			{
				"name": "func_callContains_2",
				"command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"expected": "halloo"
			},
			{
				"name": "func_shell",
				"command": "shell('echo hallo')",
				"expected": "hallo"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : func_callCompare_1 separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact func_callCompare_1 successfully saved",
		"[INFO] : Name: func_callCompare_1",
		"func_callCompare_1: if(callCompare('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_1 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"[INFO] : artefact func_callCompare_2 successfully saved",
		"[INFO] : Name: func_callCompare_2",
		"func_callCompare_2: if(callCompare('echo halloo', 'hallo')) printToConsole('should not print') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_1 separated in: echo",
		"[INFO] : artefact func_callContains_1 successfully saved",
		"[INFO] : Name: func_callContains_1",
		"func_callContains_1: if(callContains('echo hallo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callContains_1 finished! (Executed + output == expected)",
		"[INFO] : func_callContains_2 separated in: echo",
		"[INFO] : artefact func_callContains_2 successfully saved",
		"[INFO] : Name: func_callContains_2",
		"func_callContains_2: if(callContains('echo halloo', 'hallo')) printToConsole('should be true') executed",
		"[INFO] : func_callCompare_2 finished! (Executed + output == expected)",
		"[INFO] : func_callCompare_2 separated in: echo",
		"artefact func_shell successfully saved",
		"[INFO] : Name: func_shell",
		"func_shell: shell('echo hallo') executed",
		"[INFO] : func_shell finished! (Executed + output == expected)",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "func_callCompare_1",
				"Command": "if(callCompare('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callCompare_2",
				"Command": "if(callCompare('echo halloo', 'hallo')) printToConsole('should not print')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_1",
				"Command": "if(callContains('echo hallo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_callContains_2",
				"Command": "if(callContains('echo halloo', 'hallo')) printToConsole('should be true')",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "func_shell",
				"Command": "shell('echo hallo')",
				"Command was executed": true,
				"Output is as expected": true
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.FileExists(t, pathDebug, "Check \""+pathDebug+"\" file exists")

	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")

	deleteFile(flags.input)
	deleteOutput()
	flags.debug = false
	flags.skipSanity = false
}

func TestMainDifferentTypeExpected(t *testing.T) {

	configFileName := "theConfig.json"
	flags.input = "./tests/" + configFileName

	configContent := `{
        "commands": [
			{
				"name": "different_TypeExpected_1",
				"command": "echo hallo",
				"typeExpected": "==",
				"expected": "hallo"
			},
			{
				"name": "different_TypeExpected_2",
				"command": "echo hallo",
				"expected": "hallo"
			},
			{
				"name": "different_TypeExpected_3",
				"command": "echo hallo",
				"typeExpected": "contains",
				"expected": "hallo"
			},
			{
				"name": "different_TypeExpected_4",
				"command": "echo hallo",
				"typeExpected": "containsReg",
				"expected": "hallo"
			},
			{
				"name": "different_TypeExpected_5",
				"command": "echo hallo",
				"typeExpected": "!=",
				"expected": "halloo"
			},
			{
				"name": "different_TypeExpected_6",
				"command": "echo \"\"",
				"typeExpected": "nil",
				"expected": ""
			},
			{
				"name": "different_TypeExpected_7",
				"command": "echo \"\"",
				"typeExpected": "nil"
			}
        ],
        "system":
		{
			` + getConfigSystemForOS() + `
		}
	}`

	expectedAuditLogContent := []string{
		"[INFO] : Create output folder",
		"[INFO] : input-Path \"./tests/theConfig.json\" is correct",
		"[INFO] : input-File was opened",
		"[INFO] : JSON format correct",
		"[INFO] : different_TypeExpected_1 separated in: echo",
		"[INFO] : Create artefacts folder",
		"[INFO] : artefact different_TypeExpected_1 successfully saved",
		"[INFO] : Name: different_TypeExpected_1",
		"different_TypeExpected_1: echo hallo executed",
		"[INFO] : different_TypeExpected_1 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_2 separated in: echo",
		"[INFO] : Name: different_TypeExpected_2",
		"different_TypeExpected_2: echo hallo executed",
		"[INFO] : different_TypeExpected_2 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_3 separated in: echo",
		"[INFO] : artefact different_TypeExpected_3 successfully saved",
		"[INFO] : Name: different_TypeExpected_3",
		"different_TypeExpected_3: echo hallo executed",
		"[INFO] : different_TypeExpected_3 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_4 separated in: echo",
		"[INFO] : artefact different_TypeExpected_4 successfully saved",
		"[INFO] : Name: different_TypeExpected_4",
		"different_TypeExpected_4: echo hallo executed",
		"[INFO] : different_TypeExpected_4 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_5 separated in: echo",
		"[INFO] : artefact different_TypeExpected_5 successfully saved",
		"[INFO] : Name: different_TypeExpected_5",
		"different_TypeExpected_5: echo hallo executed",
		"[INFO] : different_TypeExpected_5 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_6 separated in: echo",
		"[INFO] : artefact different_TypeExpected_6 successfully saved",
		"[INFO] : Name: different_TypeExpected_6",
		"different_TypeExpected_6: echo \"\" executed",
		"[INFO] : different_TypeExpected_6 finished! (Executed + output == expected)",
		"[INFO] : different_TypeExpected_7 separated in: echo",
		"[INFO] : artefact different_TypeExpected_7 successfully saved",
		"[INFO] : Name: different_TypeExpected_7",
		"different_TypeExpected_7: echo \"\" executed",
		"[INFO] : different_TypeExpected_7 finished! (Executed + output == expected)",
	}

	expectedResultJSONContent := `{
		"./tests/theConfig.json": [
			{
				"Name": "different_TypeExpected_1",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_2",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_3",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_4",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_5",
				"Command": "echo hallo",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_6",
				"Command": "echo \"\"",
				"Command was executed": true,
				"Output is as expected": true
			},
			{
				"Name": "different_TypeExpected_7",
				"Command": "echo \"\"",
				"Command was executed": true,
				"Output is as expected": true
			}
		]
	}`

	blackBoxWriter(configContent, configFileName)

	assert.NotPanics(t, func() { main() }, "Check \"main()\" throw expected panic")

	assert.FileExists(t, pathAudit, "Check \""+pathAudit+"\" file exists")
	CheckFileContent(t, pathAudit, "", expectedAuditLogContent)

	assert.FileExists(t, pathResult, "Check \""+pathResult+"\" file exists")
	CheckFileContent(t, pathResult, expectedResultJSONContent, nil)

	assert.NoFileExists(t, pathError, "Check \""+pathError+"\" file doesn't exists")

	pathArtefacts := "./output/artefacts/"
	different_TypeExpected_1 := pathArtefacts + "different_TypeExpected_1.txt"
	different_TypeExpected_2 := pathArtefacts + "different_TypeExpected_2.txt"
	different_TypeExpected_3 := pathArtefacts + "different_TypeExpected_3.txt"
	different_TypeExpected_4 := pathArtefacts + "different_TypeExpected_4.txt"
	different_TypeExpected_5 := pathArtefacts + "different_TypeExpected_5.txt"
	different_TypeExpected_6 := pathArtefacts + "different_TypeExpected_6.txt"
	different_TypeExpected_7 := pathArtefacts + "different_TypeExpected_7.txt"

	assert.FileExists(t, different_TypeExpected_1, "Check \""+different_TypeExpected_1+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_2, "Check \""+different_TypeExpected_2+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_3, "Check \""+different_TypeExpected_3+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_4, "Check \""+different_TypeExpected_4+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_5, "Check \""+different_TypeExpected_5+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_6, "Check \""+different_TypeExpected_6+"\" artefact exists")
	assert.FileExists(t, different_TypeExpected_7, "Check \""+different_TypeExpected_7+"\" artefact exists")

	deleteFile(flags.input)
	deleteOutput()
}

func getConfigSystemForOS() string {
	if strings.EqualFold(runtime.GOOS, "windows") {
		return `          
		"systemName": "Windows",
		"version": "10.0.19042",
		"shell": "powershell",
		"argument": "/C"`
	}

	return `          
	"systemName": "Linux",
	"version": "5.2",
	"shell": "bash",
	"argument": "-c"`
}

func blackBoxWriter(output string, fileName string) error {
	path := "./tests/"

	createFolderIfNotExist(path, "tests")

	path += fileName
	openFileFlag := os.O_CREATE + os.O_WRONLY

	if checkPathExists(path) {
		deleteFile(path)
	}

	file, _ := os.OpenFile(path, openFileFlag, 0666)
	file.Write([]byte(output))
	defer file.Close()

	return nil
}

func TestBlackBoxWriter(t *testing.T) {
	assert.NoError(t, blackBoxWriter(testText, "audit.log"), "blackBoxWriter should not return error")

	CheckFileExists(t, "./tests/audit.log")
	CheckFileContent(t, "./tests/audit.log", testText, nil)

	os.Remove("./tests/audit.log")
}
