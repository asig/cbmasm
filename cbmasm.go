/*
 * Copyright (c) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of cbmasm.
 *
 * cbmasm is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * cbmasm is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cbmasm.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"flag"
	"fmt"
	"github.com/asig/cbmasm/pkg/asm"
	"github.com/asig/cbmasm/pkg/text"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	errorOutput  = log.New(os.Stderr, "", 0)
	statusOutput = log.New(os.Stdout, "", 0)
)

var (
	flagIncludeDirs = flag.String("I", ".", "include paths")
	flagPlain       = flag.Bool("plain", false, "If true, the load address is not added to the generated code.")
)

func usage() {
	errorOutput.Println("Usage: c128asm [inputfile] [outputfile] [-plain]")
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()

	inputFilename := "<stdin>"
	outputFilename := "<stdout>"
	var err error
	inputFile := os.Stdin
	outputFile := os.Stdout
	if len(args) > 0 {
		inputFilename = args[0]
		inputFile, err = os.Open(inputFilename)
		if err != nil {
			log.Fatalf("Can't open input file %q.", inputFilename)
		}
		defer inputFile.Close()
	}
	if len(args) > 1 {
		outputFilename = args[1]
		outputFile, err = os.Create(outputFilename)
		if err != nil {
			log.Fatalf("Can't open output file %q.", outputFilename)
		}
		defer outputFile.Close()
	} else {
		// Oytput is written to stdout, don't use it for status updates
		statusOutput = errorOutput
	}
	if len(args) > 2 {
		usage()
	}

	raw, err := ioutil.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	t := text.Process(inputFilename, string(raw))

	assembler := asm.New(strings.Split(*flagIncludeDirs, ":"))
	assembler.Assemble(t)
	errors := assembler.Errors()
	if len(errors) > 0 {
		errorOutput.Printf("%d errors occurred:\n", len(errors))
		for _, e := range errors {
			errorOutput.Printf("%s\n", e)
		}
	}
	warnings := assembler.Warnings()
	if len(warnings) > 0 {
		fmt.Printf("%d warnings occurred:\n", len(warnings))
		for _, e := range warnings {
			fmt.Printf("%s\n", e)
		}
	}
	if len(errors) != 0 {
		return
	}

	if !*flagPlain {
		o := assembler.Origin()
		outputFile.Write([]byte{byte(o & 0xff), byte((o >> 8) & 0xff)})
	}
	bytes := assembler.GetBytes()
	outputFile.Write(bytes)

	statusOutput.Printf("%d bytes written to %q.", len(bytes), outputFilename)
}
