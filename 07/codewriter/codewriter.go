package codewriter

import (
	"os"
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/07/token"
)

type CodeWriter struct {
	filename string
	assembly []string
}

func New(filename string) *CodeWriter {
	return &CodeWriter{
		filename: filename,
		assembly: generateInit(),
	}
}

func (c *CodeWriter) WriteArithmetic(command token.CommandSymbol) {
	switch command {
	case token.ADD:
		c.assembly = append(c.assembly, generateAdd()...)
	case token.EQ:
		c.assembly = append(c.assembly, generateEQ()...)
	}
}

func (c *CodeWriter) WritePushPop(command token.CommandType, segment token.Segment, index int) {
	switch command {
	case token.C_PUSH:
		c.assembly = append(c.assembly, generatePush(segment, index)...)
	case token.C_POP:
	default:
	}
}

func (c *CodeWriter) Close() {
	// infinite loop
	c.assembly = append(c.assembly, "(END)")
	c.assembly = append(c.assembly, "  @END")
	c.assembly = append(c.assembly, "  0;JMP")

	err := os.WriteFile(c.filename, []byte(strings.Join(c.assembly, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func generateInit() []string {
	var result []string
	result = append(result, "@256")
	result = append(result, "D=A")
	result = append(result, "@SP")
	result = append(result, "M=D")
	return result
}

func generatePush(segment token.Segment, index int) []string {
	switch segment {
	case token.SEGMENT_CONSTANT:
		return generatePushConstant(index)
	default:
		return nil
	}
}

func generatePushConstant(index int) []string {
	var result []string
	result = append(result, "@"+strconv.Itoa(index))
	result = append(result, "D=A")
	result = append(result, "@SP")
	result = append(result, "A=M")
	result = append(result, "M=D")
	result = append(result, "@SP")
	result = append(result, "M=M+1")
	return result
}

func generateAdd() []string {
	var result []string
	result = append(result, "@SP")
	result = append(result, "AM=M-1")
	result = append(result, "D=M")
	result = append(result, "A=A-1")
	result = append(result, "M=D+M")
	return result
}

func generateEQ() []string {
	var result []string
	result = append(result, "@SP")
	result = append(result, "AM=M-1")
	result = append(result, "D=M")
	result = append(result, "A=A-1")
	result = append(result, "D=D-M")
	result = append(result, "@EQ_TRUE")
	result = append(result, "D;JEQ")
	result = append(result, "@SP")
	result = append(result, "A=M-1")
	result = append(result, "M=0")
	result = append(result, "@EQ_END")
	result = append(result, "0;JMP")
	result = append(result, "(EQ_TRUE)")
	result = append(result, "@SP")
	result = append(result, "A=M-1")
	result = append(result, "M=-1")
	result = append(result, "(EQ_END)")
	return result
}
