package codewriter

import (
	"os"
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/08/token"
)

type CodeWriter struct {
	filename     string
	assembly     []string
	compareCount int
	callCount    int
}

func New(filename string) *CodeWriter {
	return &CodeWriter{
		filename:     filename,
		assembly:     generateInit(),
		compareCount: 0,
		callCount:    0,
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
		c.assembly = append(c.assembly, generateCompare(command, c.compareCount)...)
		c.compareCount++
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
		c.assembly = append(c.assembly, generatePop(segment, index)...)
	default:
	}
}

func (c *CodeWriter) WriteLabel(label string) {
	c.assembly = append(c.assembly, "("+label+")")
}

func (c *CodeWriter) WriteGoto(label string) {
	c.assembly = append(c.assembly, "@"+label, "0;JMP")
}

func (c *CodeWriter) WriteIf(label string) {
	c.assembly = append(c.assembly, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	c.assembly = append(c.assembly, "@"+label, "D;JNE")     // if D != 0, jump to label
}

func (c *CodeWriter) WriteFunction(functionName string, numLocals int) {
	c.assembly = append(c.assembly, "("+functionName+")")
	for i := 0; i < numLocals; i++ {
		c.assembly = append(c.assembly, "@SP", "A=M", "M=0", "@SP", "M=M+1") // push 0
	}
}

func (c *CodeWriter) WriteReturn() {
	c.assembly = append(c.assembly, "@LCL", "D=M", "@R13", "M=D")                 // R13 = LCL
	c.assembly = append(c.assembly, "@5", "A=D-A", "D=M", "@R14", "M=D")          // R14 = *(LCL-5)
	c.assembly = append(c.assembly, "@SP", "AM=M-1", "D=M", "@ARG", "A=M", "M=D") // *ARG = pop()
	c.assembly = append(c.assembly, "@ARG", "D=M+1", "@SP", "M=D")                // SP = ARG + 1
	c.assembly = append(c.assembly, "@R13", "AM=M-1", "D=M", "@THAT", "M=D")      // THAT = *(LCL-1)
	c.assembly = append(c.assembly, "@R13", "AM=M-1", "D=M", "@THIS", "M=D")      // THIS = *(LCL-2)
	c.assembly = append(c.assembly, "@R13", "AM=M-1", "D=M", "@ARG", "M=D")       // ARG = *(LCL-3)
	c.assembly = append(c.assembly, "@R13", "AM=M-1", "D=M", "@LCL", "M=D")       // LCL = *(LCL-4)
	c.assembly = append(c.assembly, "@R14", "A=M", "0;JMP")                       // goto return address
}

func (c *CodeWriter) WriteCall(functionName string, numArgs int) {
	returnAddress := functionName + "$ret." + strconv.Itoa(c.callCount)
	c.callCount++

	c.assembly = append(c.assembly, "@"+returnAddress, "D=A", "@SP", "A=M", "M=D", "@SP", "M=M+1")                                               // push return address
	c.assembly = append(c.assembly, "@LCL", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")                                                          // push LCL
	c.assembly = append(c.assembly, "@ARG", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")                                                          // push ARG
	c.assembly = append(c.assembly, "@THIS", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")                                                         // push THIS
	c.assembly = append(c.assembly, "@THAT", "D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1")                                                         // push THAT
	c.assembly = append(c.assembly, "@SP", "D=M", "@5", "D=D-A", "@"+strconv.Itoa(numArgs), "D=D-A", "@ARG", "M=D", "@SP", "D=M", "@LCL", "M=D") // ARG = SP - 5 - numArgs, LCL = SP
	c.assembly = append(c.assembly, "@"+functionName, "0;JMP")                                                                                   // goto functionName
	c.assembly = append(c.assembly, "("+returnAddress+")")                                                                                       // (returnAddress)
}

func (c *CodeWriter) Close() {
	// infinite loop
	// NOTE: "END" is not a reserved label in Hack assembly
	c.assembly = append(c.assembly, "(END)", "@END", "0;JMP")
	err := os.WriteFile(c.filename, []byte(strings.Join(c.assembly, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func generateInit() []string {
	var result []string
	return result
}

func generatePush(segment token.Segment, index int) []string {
	switch segment {
	case token.SEGMENT_CONSTANT:
		return generatePushConstant(index)
	case token.SEGMENT_LOCAL, token.SEGMENT_ARGUMENT, token.SEGMENT_THIS, token.SEGMENT_THAT:
		return generatePushMemoryAccess(segment, index)
	case token.SEGMENT_POINTER:
		return generatePushPointer(index)
	case token.SEGMENT_STATIC:
		return generatePushStatic(index)
	case token.SEGMENT_TEMP:
		return generatePushTemp(index)
	default:
		return nil
	}
}

func generatePushConstant(index int) []string {
	var result []string
	result = append(result, "@"+strconv.Itoa(index), "D=A") // D = index
	result = append(result, "@SP", "A=M", "M=D")            // RAM[SP] = D
	result = append(result, "@SP", "M=M+1")                 // SP++
	return result
}

func generatePushMemoryAccess(segment token.Segment, index int) []string {
	var result []string
	var segmentAddr string
	switch segment {
	case token.SEGMENT_LOCAL:
		segmentAddr = "LCL"
	case token.SEGMENT_ARGUMENT:
		segmentAddr = "ARG"
	case token.SEGMENT_THIS:
		segmentAddr = "THIS"
	case token.SEGMENT_THAT:
		segmentAddr = "THAT"
	}
	result = append(result, "@"+strconv.Itoa(index), "D=A")  // D = index
	result = append(result, "@"+segmentAddr, "A=D+M", "D=M") // D = RAM[index + segmentAddr]
	result = append(result, "@SP", "A=M", "M=D")             // RAM[SP] = D
	result = append(result, "@SP", "M=M+1")                  // SP++
	return result
}

func generatePushPointer(index int) []string {
	var result []string
	var pointer string
	if index == 0 {
		pointer = "THIS"
	} else {
		pointer = "THAT"
	}
	result = append(result, "@"+pointer, "D=M")  // D = pointer
	result = append(result, "@SP", "A=M", "M=D") // RAM[SP] = D
	result = append(result, "@SP", "M=M+1")      // SP++
	return result
}

func generatePushStatic(index int) []string {
	var result []string
	result = append(result, "@STATIC"+strconv.Itoa(index), "D=M") // D = RAM[STATIC + index]
	result = append(result, "@SP", "A=M", "M=D")                  // RAM[SP] = D
	result = append(result, "@SP", "M=M+1")                       // SP++
	return result
}

func generatePushTemp(index int) []string {
	var result []string
	result = append(result, "@R5", "D=A")                            // D = 5
	result = append(result, "@"+strconv.Itoa(index), "A=D+A", "D=M") // D = RAM[5 + index]
	result = append(result, "@SP", "A=M", "M=D")                     // RAM[SP] = D
	result = append(result, "@SP", "M=M+1")                          // SP++
	return result
}

func generatePop(segment token.Segment, index int) []string {
	switch segment {
	case token.SEGMENT_LOCAL, token.SEGMENT_ARGUMENT, token.SEGMENT_THIS, token.SEGMENT_THAT:
		return generatePopMemoryAccess(segment, index)
	case token.SEGMENT_POINTER:
		return generatePopPointer(index)
	case token.SEGMENT_STATIC:
		return generatePopStatic(index)
	case token.SEGMENT_TEMP:
		return generatePopTemp(index)
	default:
		return nil
	}
}

func generatePopMemoryAccess(segment token.Segment, index int) []string {
	var result []string
	var segmentAddr string
	switch segment {
	case token.SEGMENT_LOCAL:
		segmentAddr = "LCL"
	case token.SEGMENT_ARGUMENT:
		segmentAddr = "ARG"
	case token.SEGMENT_THIS:
		segmentAddr = "THIS"
	case token.SEGMENT_THAT:
		segmentAddr = "THAT"
	}
	result = append(result, "@"+strconv.Itoa(index), "D=A") // D = index
	result = append(result, "@"+segmentAddr, "D=D+M")       // D = index + segmentAddr
	result = append(result, "@R13", "M=D")                  // R13 = D (temporarily store the address to pop)
	result = append(result, "@SP", "AM=M-1", "D=M")         // move RAM[SP-1] to D
	result = append(result, "@R13", "A=M", "M=D")           // RAM[R13] = D
	return result
}

func generatePopPointer(index int) []string {
	var result []string
	var pointer string
	if index == 0 {
		pointer = "THIS"
	} else {
		pointer = "THAT"
	}
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "@"+pointer, "M=D")     // pointer = D
	return result
}

func generatePopStatic(index int) []string {
	var result []string
	result = append(result, "@SP", "AM=M-1", "D=M") // move RAM[SP-1] to D
	result = append(result, "@STATIC"+strconv.Itoa(index), "M=D")
	return result
}

func generatePopTemp(index int) []string {
	var result []string
	result = append(result, "@R5", "D=A")                     // D = 5
	result = append(result, "@"+strconv.Itoa(index), "D=D+A") // D = 5 + index
	result = append(result, "@R13", "M=D")                    // R13 = D (temporarily store the address to pop)
	result = append(result, "@SP", "AM=M-1", "D=M")           // move RAM[SP-1] to D
	result = append(result, "@R13", "A=M", "M=D")             // RAM[R13] = D
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

func generateCompare(command token.CommandSymbol, compareCount int) []string {
	flag := strconv.Itoa(compareCount)
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
