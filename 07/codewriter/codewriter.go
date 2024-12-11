package codewriter

import (
	"os"
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/07/token"
	"golang.org/x/exp/rand"
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
	case token.SUB:
		c.assembly = append(c.assembly, generateSUB()...)
	case token.NEG:
		c.assembly = append(c.assembly, generateNEG()...)
	case token.EQ, token.LT, token.GT:
		c.assembly = append(c.assembly, generateCompare(command)...)
	case token.AND:
		c.assembly = append(c.assembly, generateAND()...)
	case token.OR:
		c.assembly = append(c.assembly, generateOR()...)
	case token.NOT:
		c.assembly = append(c.assembly, generateNOT()...)
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
	c.assembly = append(c.assembly, "@END")
	c.assembly = append(c.assembly, "0;JMP")

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
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "A=A-1", "M=D+M")       // M = RAM[SP-1] + RAM[SP-2]
	return result
}

func generateSUB() []string {
	var result []string
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "A=A-1", "M=M-D")       // M = RAM[SP-2] - RAM[SP-1]
	return result
}

func generateNEG() []string {
	var result []string
	result = append(result, "@SP", "A=M-1", "M=-M") // negate RAM[SP-1]
	return result
}

func generateCompare(command token.CommandSymbol) []string {
	flag := strconv.Itoa(rand.Intn(1000000))
	var result []string
	var jump string
	switch command {
	case token.EQ:
		jump = "JEQ" // D == 0
	case token.LT:
		jump = "JLT" // D < 0
	case token.GT:
		jump = "JGT" // D > 0
	}
	result = append(result, "@SP", "AM=M-1", "D=M")                   // move RAM[SP-1] to D
	result = append(result, "A=A-1", "D=M-D")                         // D = RAM[SP-2] - RAM[SP-1]
	result = append(result, "@TRUE"+flag, "D;"+jump)                  // if <jump>, jump to TRUE
	result = append(result, "@SP", "A=M-1", "M=0")                    // set RAM[SP-2] to 0 (false)
	result = append(result, "@END"+flag, "0;JMP")                     // jump to END
	result = append(result, "(TRUE"+flag+")", "@SP", "A=M-1", "M=-1") // set RAM[SP-2] to -1 (true)
	result = append(result, "(END"+flag+")")
	return result
}

func generateAND() []string {
	var result []string
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "A=A-1", "M=D&M")       // AND RAM[SP-1] with RAM[SP-2]
	return result
}

func generateOR() []string {
	var result []string
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "A=A-1", "M=D|M")       // OR RAM[SP-1] with RAM[SP-2]
	return result
}

func generateNOT() []string {
	var result []string
	result = append(result, "@SP", "A=M-1", "M=!M") // NOT RAM[SP-1]
	return result
}
