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
	"flag"
	"strings"
)

type Flags struct {
	input         string
	output        string
	addedCommands []string
	verbose       bool
	debug         bool
	skipSanity    bool
	help          bool
	encryptZip    bool
}

var flags Flags

// get flags once
func initFlags() error {
	if flag.Parsed() {
		return nil
	}

	input := flag.String("input", "", "specify config file")
	output := flag.String("output", "", "specify output file")
	add := flag.String("add", "", "add specific commands")
	verbose := flag.Bool("v", false, "toggle verbose mode")
	debug := flag.Bool("debug", false, "toggle debug mode")
	skipSanity := flag.Bool("s", false, "skip sanity check")
	help := flag.Bool("h", false, "help")
	encryptZip := flag.Bool("p", false, "encrypt Zip")
	flag.Parse()

	var commands []string
	if *add != "" {
		*add = strings.ReplaceAll(*add, " ", "")
		commands = strings.Split(*add, ",")
	}

	err := checkValidFilename(*output)
	if err != nil {
		return errors.New("The character \"" + err.Error() + "\" in output is not allowed")
	}

	flags.input, flags.output, flags.addedCommands = *input, *output, commands
	flags.verbose, flags.skipSanity, flags.help = *verbose, *skipSanity, *help
	flags.debug = *debug
	flags.encryptZip = *encryptZip
	return nil
}
