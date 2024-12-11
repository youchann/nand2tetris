package codewriter

import (
	"os"
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/07/token"
)

type CodeWriter struct {
	filename string
	assembly string
}

func New(filename string) *CodeWriter {
	return &CodeWriter{
		filename: filename,
		assembly: "",
	}
}

func (c *CodeWriter) WriteArithmetic(command token.CommandSymbol) {
	switch command {
	case token.ADD:
		c.assembly += generateAdd() + "\n"
	}
}

func (c *CodeWriter) WritePushPop(command token.CommandType, segment token.Segment, index int) {
	switch command {
	case token.C_PUSH:
		c.assembly += generatePush(segment, index) + "\n"
	case token.C_POP:
	default:
	}
}

func (c *CodeWriter) Close() {
	// infinite loop
	c.assembly += "(END)\n"
	c.assembly += "  @END\n"
	c.assembly += "  0;JMP"

	err := os.WriteFile(c.filename, []byte(c.assembly), 0644)
	if err != nil {
		panic(err)
	}
}

func generatePush(segment token.Segment, index int) string {
	switch segment {
	case token.SEGMENT_CONSTANT:
		return generatePushConstant(index)
	default:
		return ""
	}
}

func generatePushConstant(index int) string {
	var result []string
	result = append(result, "@"+strconv.Itoa(index))
	result = append(result, "D=A")
	result = append(result, "@SP")
	result = append(result, "A=M")
	result = append(result, "M=D")
	result = append(result, "@SP")
	result = append(result, "M=M+1")
	return strings.Join(result, "\n")
}

func generateAdd() string {
	var result []string
	result = append(result, "@SP")
	result = append(result, "AM=M-1")
	result = append(result, "D=M")
	result = append(result, "A=A-1")
	result = append(result, "M=D+M")
	return strings.Join(result, "\n")
}
