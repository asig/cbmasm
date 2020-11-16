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
	"sort"
	"strings"
)

var (
	errorOutput  = log.New(os.Stderr, "", 0)
	statusOutput = log.New(os.Stdout, "", 0)
)

type stringArrayFlag []string

func (f *stringArrayFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringArrayFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type pathListFlag []string

func (f *pathListFlag) String() string {
	return strings.Join(*f, string(os.PathListSeparator))
}

func (f *pathListFlag) Set(value string) error {
	for _, p := range strings.Split(value, string(os.PathListSeparator)) {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			*f = append(*f, p)
		}
	}
	return nil
}

var (
	flagIncludeDirs pathListFlag
	flagDefines     stringArrayFlag
	flagPlain       = flag.Bool("plain", false, "If true, the load address is not added to the generated code.")
	flagDumpLabels  = flag.Bool("dump_labels", true, "If true, the labels will be printed.")
	flagListing     = flag.Bool("listing", false, "If true, a listing is generated.")
)

func usage() {
	errorOutput.Println("Usage: c128asm {-I includedir} {-D define} [-plain] [-dump_labels] [-listing] [inputfile] [outputfile] ")
	os.Exit(1)
}

func printLabels(a *asm.Assembler) {
	labels := a.Labels()
	names := make([]string, 0, len(labels))
	maxLen := 0
	for n := range labels {
		names = append(names, n)
		l := len(n)
		if l > maxLen {
			maxLen = l
		}
	}
	sort.Strings(names)
	statusOutput.Println("Labels:")
	for _, n := range names {
		val := labels[n]
		for len(n) < maxLen {
			n = n + " "
		}
		statusOutput.Printf("%s: $%04x\n", n, val)
	}
	statusOutput.Println()
}

func printListing(a *asm.Assembler) {
	for _, l := range a.ListingLines {
		bytes := []byte{}
		if l.Bytes > 0 {
			start := l.Addr - a.Origin()
			bytes = a.GetBytes()[start : start+l.Bytes]
		}
		var byteStrs []string
		for _, b := range bytes {
			byteStrs = append(byteStrs, fmt.Sprintf("%02x", b))
		}
		for len(byteStrs) < 8 {
			byteStrs = append(byteStrs, "  ")
		}
		statusOutput.Printf("%04x | %s | %s\n", l.Addr, strings.Join(byteStrs, " "), strings.TrimSuffix(string(l.Line.Runes), "\n"))
	}
}

func main() {
	flag.Var(&flagIncludeDirs, "I", "include paths; can be repeated")
	flag.Var(&flagDefines, "D", "defined symbols; can be repeated")
	flag.Parse()
	args := flag.Args()

	if len(flagIncludeDirs) == 0 {
		// default to "." if no include dirs are set.
		flagIncludeDirs = pathListFlag{"."}
	}

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

	assembler := asm.New(flagIncludeDirs)
	assembler.AddDefines(flagDefines)
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

	if *flagDumpLabels {
		printLabels(assembler)
	}
	if *flagListing {
		printListing(assembler)
	}

	statusOutput.Printf("%d bytes written to %q.", len(bytes), outputFilename)
}
