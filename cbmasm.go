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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asig/cbmasm/pkg/asm"
	"github.com/asig/cbmasm/pkg/errors"
	"github.com/asig/cbmasm/pkg/text"
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
	flagOutput      = flag.String("output", "prg", fmt.Sprintf("Which output format should be generated. Supported values are: %s", strings.Join(asm.SupportedOutputs, ", ")))
	flagEncoding    = flag.String("encoding", "petscii", fmt.Sprintf("Which encoding should be used. Supported values are: %s", strings.Join(asm.SupportedEncodings, ", ")))
	flagDumpLabels  = flag.Bool("dump_labels", false, "If true, the labels will be printed to stdout.")
	flagLabels      = flag.String("labels", "", "If set, a VICE-compatible 'labels' file is generated.")
	flagListing     = flag.Bool("listing", false, "If true, a listing is generated.")
	flagCPU         = flag.String("cpu", "6502", fmt.Sprintf("CPU to assemble code for. Supported values are: %s", strings.Join(asm.SupportedCPUs, ", ")))
	flagPlatform    = flag.String("platform", "c128", fmt.Sprintf("Target platform. Supported values are: %s", strings.Join(asm.SupportedPlatforms, ", ")))
)

func usage() {
	errorOutput.Printf("Usage: %s [flags] [inputfile] [outputfile]\n", filepath.Base(os.Args[0]))
	errorOutput.Println("Flags:")
	flag.PrintDefaults()
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

func saveViceLabels(a *asm.Assembler, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		log.Printf("Can't open output file %q.", filename)
		return
	}
	labels := a.Labels()
	var symtab []string
	for n, addr := range labels {
		if !strings.HasPrefix(n, ".") {
			n = "." + n
		}
		symtab = append(symtab, fmt.Sprintf("al C:%04x %s\n", addr, n))
	}
	sort.Strings(symtab)
	for _, l := range symtab {
		out.WriteString(l)
	}
	out.Close()
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
		for len(byteStrs) < 5 {
			byteStrs = append(byteStrs, "  ")
		}
		statusOutput.Printf("%04x | %s | %s\n", l.Addr, strings.Join(byteStrs, " "), strings.TrimSuffix(string(l.Line.Runes), "\n"))
	}
}

func init() {
	flag.Usage = usage
	flag.Var(&flagIncludeDirs, "I", "include paths; can be repeated")
	flag.Var(&flagDefines, "D", "defined symbols; can be repeated")
	flag.Parse()

	if !asm.IsSupportedPlatform(*flagPlatform) {
		errorOutput.Printf("Unsupported platform %q. Valid platforms are: %s.", *flagPlatform, strings.Join(asm.SupportedPlatforms, ", "))
		usage()
		os.Exit(1)
	}

	if !asm.IsSupportedCPU(*flagCPU) {
		errorOutput.Printf("Unsupported CPU %q. Valid CPUs are: %s.", *flagCPU, strings.Join(asm.SupportedCPUs, ", "))
		usage()
		os.Exit(1)
	}

	if !asm.IsValidPlatformCPUCombo(*flagPlatform, *flagCPU) {
		errorOutput.Printf("Platform %q is not supported for CPU %q.", *flagPlatform, *flagCPU)
		usage()
		os.Exit(1)
	}

	if !asm.IsSupportedOutput(*flagOutput) {
		errorOutput.Printf("Unsupported output %q. Valid outputs are: %s.", *flagOutput, strings.Join(asm.SupportedOutputs, ", "))
		usage()
		os.Exit(1)
	}

	if !asm.IsSupportedEncoding(*flagEncoding) {
		errorOutput.Printf("Unsupported encoding %q. Valid encodings are: %s.", *flagEncoding, strings.Join(asm.SupportedEncodings, ", "))
		usage()
		os.Exit(1)
	}

	if len(flagIncludeDirs) == 0 {
		// default to "." if no include dirs are set.
		flagIncludeDirs = pathListFlag{"."}
	}
}

func main() {
	args := flag.Args()

	inputFilename := "<stdin>"
	outputFilename := "<stdout>"

	var errs []errors.Error
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
		defer func() {
			outputFile.Close()
			if len(errs) > 0 {
				os.Remove(outputFilename)
			}
		}()
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

	assembler := asm.New(flagIncludeDirs, *flagCPU, *flagPlatform, *flagOutput, *flagEncoding, flagDefines)
	assembler.Assemble(t)
	errs = assembler.Errors()
	if len(errs) > 0 {
		errorOutput.Printf("%d errors occurred:\n", len(errs))
		for _, e := range errs {
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
	if len(errs) != 0 {
		return
	}

	output := assembler.CurrentOutput()
	if output == "prg" {
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
	if *flagLabels != "" {
		saveViceLabels(assembler, *flagLabels)
		statusOutput.Printf("Symbols written to %q.", *flagLabels)
	}

	statusOutput.Printf("%d bytes written to %q.", len(bytes), outputFilename)
}
