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

func assemble(content string) ([]string, error) {
	var machineCode []string
	p := parser.New(content)
	for p.HasMoreLines() {
		var instruction string
		switch p.CommandType() {
		case parser.A_INSTRUCTION:
			instruction = code.Symbol(p.Symbol())
		case parser.C_INSTRUCTION:
			instruction = "111" + code.Comp(p.Comp()) + code.Dest(p.Dest()) + code.Jump(p.Jump())
		case parser.L_INSTRUCTION:
			// Lコマンドはマシン語を生成しない
			p.Advance()
			continue
		}
		if instruction != "" {
			machineCode = append(machineCode, instruction)
		}
		p.Advance()
	}

	return machineCode, nil
}

func writeToFile(filepath string, instructions []string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	for _, instruction := range instructions {
		_, err := file.WriteString(instruction + "\n")
		if err != nil {
			return fmt.Errorf("failed to write instruction to file: %v", err)
		}
	}

	return nil
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

	machineCode, err := assemble(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error assembling code: %v\n", err)
		os.Exit(1)
	}

	resultFilePath := getHackFilePath(filename)
	err = writeToFile(resultFilePath, machineCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
}
