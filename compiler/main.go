package main

import (
	"fmt"
	"os"

	"github.com/aireya/aireyac/lexer"
	"github.com/aireya/aireyac/parser"
	"github.com/aireya/aireyac/generator/gen_cpp"
	"github.com/aireya/aireyac/generator/gen_go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.aya>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Lexing
	l := lexer.New(f)

	// Parsing
	p := parser.New(l)
	astFile := p.ParseFile()

	if len(p.Errors()) != 0 {
		fmt.Fprintf(os.Stderr, "Syntax Errors:\n")
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "  %s\n", msg)
		}
		os.Exit(1)
	}

	// Generating C++
	cppGen := gen_cpp.New(astFile)
	out, err := cppGen.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating C++: %v\n", err)
		os.Exit(1)
	}

	outFileName := inputFile + ".h"
	err = os.WriteFile(outFileName, []byte(out), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully generated %s\n", outFileName)

	// Generating Go
	goGen := gen_go.New(astFile)
	outGo, err := goGen.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Go: %v\n", err)
		os.Exit(1)
	}
	outGoName := inputFile + ".go"
	err = os.WriteFile(outGoName, []byte(outGo), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully generated %s\n", outGoName)
}
