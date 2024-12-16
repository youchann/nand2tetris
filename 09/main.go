package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/youchann/nand2tetris/09/compilationengine"
	"github.com/youchann/nand2tetris/09/tokenizer"
)

func getXMLPath(jackFilePath string) string {
	dir := filepath.Dir(jackFilePath)
	baseFile := filepath.Base(jackFilePath)
	xmlFileName := strings.TrimSuffix(baseFile, ".jack") + "_.xml"
	return filepath.Join(dir, xmlFileName)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [filename.jack or directory]")
		os.Exit(1)
	}

	path := os.Args[1]
	var jackFiles []string

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
			if filepath.Ext(entry.Name()) == ".jack" {
				jackFiles = append(jackFiles, filepath.Join(path, entry.Name()))
			}
		}
		if len(jackFiles) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No .jack files found in directory\n")
			os.Exit(1)
		}
	} else {
		if filepath.Ext(path) != ".jack" {
			fmt.Fprintf(os.Stderr, "Error: File must have .jack extension\n")
			os.Exit(1)
		}
		jackFiles = append(jackFiles, path)
	}

	for _, jackFile := range jackFiles {
		content, err := os.ReadFile(jackFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", jackFile, err)
			os.Exit(1)
		}

		xmlPath := getXMLPath(jackFile)
		xmlFile, err := os.Create(xmlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", xmlPath, err)
			os.Exit(1)
		}
		defer xmlFile.Close()

		t := tokenizer.New(string(content))
		ce := compilationengine.New(t)
		ce.CompileClass()
		xmlFile.WriteString(ce.XML)
	}
}
