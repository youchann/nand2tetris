package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/youchann/nand2tetris/06/code"
	"github.com/youchann/nand2tetris/06/parser"
)

func getHackFilePath(asmPath string) string {
	dir := filepath.Dir(asmPath)
	fileName := filepath.Base(asmPath)
	baseName := strings.TrimSuffix(fileName, ".asm")
	hackName := baseName + ".hack"
	return filepath.Join(dir, hackName)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [filename]")
		os.Exit(1)
	}

	filename := os.Args[1]
	if filepath.Ext(filename) != ".asm" {
		fmt.Fprintf(os.Stderr, "Error: File must have .asm extension\n")
		os.Exit(1)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	resultFilePath := getHackFilePath(filename)
	hackFile, err := os.Create(resultFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer hackFile.Close()

	p := parser.New(string(content))
	for p.HasMoreLines() {
		result := ""
		switch p.CommandType() {
		case parser.A_INSTRUCTION:
			result = code.Symbol(p.Symbol()) + "\n"
		case parser.C_INSTRUCTION:
			result = "111" + code.Comp(p.Comp()) + code.Dest(p.Dest()) + code.Jump(p.Jump()) + "\n"
		case parser.L_INSTRUCTION:
		}

		if result != "" {
			_, err := hackFile.WriteString(result)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
		}
		p.Advance()
	}
}
