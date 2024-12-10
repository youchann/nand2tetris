package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/06/code"
	"github.com/youchann/nand2tetris/06/parser"
	"github.com/youchann/nand2tetris/06/symboltable"
)

func getHackFilePath(asmPath string) string {
	dir := filepath.Dir(asmPath)
	fileName := filepath.Base(asmPath)
	baseName := strings.TrimSuffix(fileName, ".asm")
	hackName := baseName + ".hack"
	return filepath.Join(dir, hackName)
}

func firstPassAssemble(content string) *symboltable.Table {
	st := symboltable.New()
	p := parser.New(content)
	romAddress := 0
	for p.HasMoreLines() {
		switch p.CommandType() {
		case parser.A_INSTRUCTION, parser.C_INSTRUCTION:
			romAddress++
		case parser.L_INSTRUCTION:
			st.AddEntry(p.Symbol(), romAddress)
		}
		p.Advance()
	}
	return st
}

func secondPassAssemble(content string, symbolTable *symboltable.Table) ([]string, error) {
	var machineCode []string
	p := parser.New(content)
	currentRAMAddress := 16
	for p.HasMoreLines() {
		var instruction string
		switch p.CommandType() {
		case parser.A_INSTRUCTION:
			s := p.Symbol()
			if _, err := strconv.Atoi(s); err != nil {
				if !symbolTable.Contains(s) {
					symbolTable.AddEntry(s, currentRAMAddress)
					currentRAMAddress++
				}
				s = strconv.Itoa(symbolTable.GetAddress(s))
			}
			instruction = code.Symbol(s)
		case parser.C_INSTRUCTION:
			instruction = "111" + code.Comp(p.Comp()) + code.Dest(p.Dest()) + code.Jump(p.Jump())
		case parser.L_INSTRUCTION: // first pass already handled this
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

	st := firstPassAssemble(string(content))
	machineCode, err := secondPassAssemble(string(content), st)
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
