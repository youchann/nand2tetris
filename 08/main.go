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

func getVMFilePath(asmPath string) string {
	dir := filepath.Dir(asmPath)
	fileName := filepath.Base(asmPath)
	baseName := strings.TrimSuffix(fileName, ".vm")
	hackName := baseName + ".asm"
	return filepath.Join(dir, hackName)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [filename]")
		os.Exit(1)
	}

	filename := os.Args[1]
	if filepath.Ext(filename) != ".vm" {
		fmt.Fprintf(os.Stderr, "Error: File must have .asm extension\n")
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	p := parser.New(string(content))
	c := codewriter.New(getVMFilePath(filename))
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
			// case token.C_FUNCTION:
			// 	c.WriteFunction(p.Arg1(), p.Arg2())
			// case token.C_RETURN:
			// 	c.WriteReturn()
			// case token.C_CALL:
			// 	c.WriteCall(p.Arg1(), p.Arg2())
		}
		p.Advance()
	}
	c.Close()
}
