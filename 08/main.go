package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/youchann/nand2tetris/08/codewriter"
	"github.com/youchann/nand2tetris/08/parser"
	"github.com/youchann/nand2tetris/08/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [filename.vm or directory]")
		os.Exit(1)
	}

	path := os.Args[1]
	var vmFiles []string

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
			os.Exit(1)
		}
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".vm" {
				vmFiles = append(vmFiles, filepath.Join(path, entry.Name()))
			}
		}
		if len(vmFiles) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No .vm files found in directory\n")
			os.Exit(1)
		}
	} else {
		if filepath.Ext(path) != ".vm" {
			fmt.Fprintf(os.Stderr, "Error: File must have .vm extension\n")
			os.Exit(1)
		}
		vmFiles = append(vmFiles, path)
	}

	outputPath := filepath.Join(
		filepath.Dir(path),
		strings.TrimSuffix(fileInfo.Name(), ".vm")+".asm",
	)

	c := codewriter.New()
	defer c.Close(outputPath)
	for _, filename := range vmFiles {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
			os.Exit(1)
		}
		c.Setfilename(strings.TrimSuffix(filepath.Base(filename), ".vm"))
		p := parser.New(string(content))
		for p.HasMoreLines() {
			switch p.CommandType() {
			case token.C_ARITHMETIC:
				c.WriteArithmetic(token.CommandSymbol(p.Arg1()))
			case token.C_PUSH, token.C_POP:
				c.WritePushPop(p.CommandType(), token.Segment(p.Arg1()), p.Arg2())
			case token.C_LABEL:
				c.WriteLabel(p.Arg1())
			case token.C_GOTO:
				c.WriteGoto(p.Arg1())
			case token.C_IF:
				c.WriteIf(p.Arg1())
			case token.C_FUNCTION:
				c.WriteFunction(p.Arg1(), p.Arg2())
			case token.C_RETURN:
				c.WriteReturn()
			case token.C_CALL:
				c.WriteCall(p.Arg1(), p.Arg2())
			}
			p.Advance()
		}
	}
}
