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
	"strconv"
	"strings"
)

var banner string
var appName string

func init() {
	appName = "Secuteel"
	banner =
		`	   _____                 __            __
	  / ___/___  _______  __/ /____  ___  / /
	  \__ \/ _ \/ ___/ / / / __/ _ \/ _ \/ / 
	 ___/ /  __/ /__/ /_/ / /_/  __/  __/ /  
	/____/\___/\___/\__,_/\__/\___/\___/_/   							 
` + "\n" + appName + " is an auditing tool, with multi-platform support, written in go.\n"
}

func printBanner() {
	if flags.input == "" {
		fmt.Println(banner + "\nFor help run following command: " + strings.ToLower(appName) + " -h")
	} else {
		fmt.Println(banner)
	}
}

func printHelpText() {
	synopsisText := `
	` + strings.ToLower(appName) + ` -input|--input [-h] [-p] [-v] [-s] [-debug] [-output|--output] [-add|--add]

`
	helpText := banner + "\n"
	helpText += strings.ToUpper("Synopsis") + synopsisText
	helpText += strings.ToUpper("Examples") + "\n\t" + strings.ToLower(appName) + " --input \"configFile.json\"\n\t" + strings.ToLower(appName) + " --input \"configFile.json\" --output \"myOutput.zip\"\n\n"
	helpText += strings.ToUpper("Description") + "\n"
	helpText += "\t-input, --input= 'expects path to config file (.json)'\n"
	helpText += "\t-output, --output= 'specify your output file (.zip)'\n"
	helpText += "\t-add, --add= 'add a list of allowed commands'\n"
	helpText += "\t-v \t'toggle verbose mode for console'\n"
	helpText += "\t-s \t'skip sanity check'\n"
	helpText += "\t-debug\t'activate debug mode for log files'\n"
	helpText += "\t-p\t'set password to encrypt output zip folder'\n"
	helpText += "\t-h\t'help'\n\n"
	helpText += strings.ToUpper("See also") + "\n\tComplete guide: https://github.com/Seculeet/secuteel\n\n"
	helpText += strings.ToUpper("Reporting Bugs") + "\n\thttps://github.com/Seculeet/secuteel/issues\n\n"
	helpText += strings.ToUpper("Copyright") + "\n\tThe " + appName + " app was published under the MIT license.\n"
	helpText += "\thttps://opensource.org/licenses/MIT\n\n"
	fmt.Print(helpText)
}

// displays progress bar while executing
func printProgressBar(allAuditslength int, auditPosition int) {
	auditPercent := getAuditPercent(allAuditslength, auditPosition)
	bar := getProgressBarUnicode(auditPercent)

	printTxt := "\r" + fmt.Sprintf("Progress: [%20s] %3d%s%5s", bar, auditPercent, " %", "")
	if allAuditslength != auditPosition {
		fmt.Print(printTxt)
	} else {
		fmt.Println(printTxt)
	}
}

// helper for progress bar
func getProgressBarUnicode(percent int) string {
	barLength := 20
	steps := 100 / barLength
	var bar string
	for i := steps; i <= 100; i += steps {
		if i <= percent {
			bar += "\u2588"
		} else {
			bar += "\u2591"
		}
	}
	return fmt.Sprintf("%20s", bar)
}

// helper for progress bar
func getAuditPercent(allAuditsLength int, auditPosition int) int {
	auditPercent := float64(auditPosition) / float64(allAuditsLength) * 100
	return int(auditPercent)
}

func printCommandStarted(index int, auditListLength int) {
	fmt.Print("Audit started (" + strconv.Itoa(index) + "/" + strconv.Itoa(auditListLength) + "): ")
}

func printCommandResult(auditResult bool) {
	var printTxt string
	if auditResult {
		printTxt = "SUCCESS"
	} else {
		printTxt = "FAILED"
	}
	fmt.Println(printTxt)
}
