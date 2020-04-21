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

var stderr = log.New(os.Stderr, "", 0)

var (
	flagIncludeDirs = flag.String("I", ".", "include paths")
)

func usage() {
	stderr.Println("Usage: c128asm [inputfile] [outputfile]")
	os.Exit(1)
}

func main() {
	flag.Parse()

	inputFilename := "<stdin>"
	var err error
	inputFile := os.Stdin
	outputFile := os.Stdout
	if len(os.Args) > 1 {
		inputFilename = os.Args[1]
		inputFile, err = os.Open(inputFilename)
		defer inputFile.Close()
		if err != nil {
			log.Fatalf("Can't open input file %q.", inputFilename)
		}
	}
	if len(os.Args) > 2 {
		outputFile, err = os.Create(os.Args[2])
		if err != nil {
			log.Fatalf("Can't open output file %q.", os.Args[2])
		}
		defer outputFile.Close()
	}
	if len(os.Args) > 3 {
		usage()
	}

	raw, err := ioutil.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	t := text.Process(inputFilename, string(raw))

	assembler := asm.New(t, strings.Split(*flagIncludeDirs, ":"))
	assembler.Assemble()
	errors := assembler.Errors()
	if len(errors) > 0 {
		fmt.Printf("%d errors occurred:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("%s\n", e)
		}
	}
	warnings := assembler.Warnings()
	if len(warnings) > 0 {
		fmt.Printf("%d warnings occurred:\n", len(warnings))
		for _, e := range warnings {
			fmt.Printf("%s\n", e)
		}
	}
	if len(errors) == 0 {
		o := assembler.Origin()
		outputFile.Write([]byte{byte(o & 0xff), byte((o >> 8) & 0xff)})
		outputFile.Write(assembler.GetBytes())
	}
}
