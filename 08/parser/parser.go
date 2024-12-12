package parser

import (
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/08/token"
)

type Parser struct {
	commandStrList []string
	currentIndex   int
}

func New(input string) *Parser {
	return &Parser{
		commandStrList: preprocessCode(input),
		currentIndex:   0,
	}
}

func (p *Parser) HasMoreLines() bool {
	return p.currentIndex < len(p.commandStrList)
}

func (p *Parser) Advance() {
	p.currentIndex++
}

func (p *Parser) CommandType() token.CommandType {
	commandStr := p.commandStrList[p.currentIndex]
	switch {
	case strings.HasPrefix(commandStr, "push"):
		return token.C_PUSH
	case strings.HasPrefix(commandStr, "pop"):
		return token.C_POP
	// case strings.HasPrefix(commandStr, "label"):
	// 	return C_LABEL
	// case strings.HasPrefix(commandStr, "goto"):
	// 	return C_GOTO
	// case strings.HasPrefix(commandStr, "if-goto"):
	// 	return C_IF
	// case strings.HasPrefix(commandStr, "function"):
	// 	return C_FUNCTION
	// case strings.HasPrefix(commandStr, "return"):
	// 	return C_RETURN
	// case strings.HasPrefix(commandStr, "call"):
	// 	return C_CALL
	default:
		return token.C_ARITHMETIC
	}
}

func (p *Parser) Arg1() string {
	commandStr := p.commandStrList[p.currentIndex]
	if p.CommandType() == token.C_ARITHMETIC {
		return strings.Fields(commandStr)[0]
	}
	return strings.Fields(commandStr)[1]
}

func (p *Parser) Arg2() int {
	commandStr := p.commandStrList[p.currentIndex]
	if i, err := strconv.Atoi(strings.Fields(commandStr)[2]); err == nil {
		return i
	}
	return 0
}

// TODO: 構文として正しいかどうかのチェックを加える
func preprocessCode(input string) []string {
	lines := strings.Split(input, "\n")
	return trimSpaces(removeEmptyLines(removeComments(lines)))
}

func trimSpaces(lines []string) []string {
	var processedLines []string
	for _, line := range lines {
		processedLines = append(processedLines, strings.TrimSpace(line))
	}
	return processedLines
}

func removeEmptyLines(lines []string) []string {
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return nonEmptyLines
}

func removeComments(lines []string) []string {
	var resultLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		if idx := strings.Index(line, "//"); idx != -1 {
			processedLine := strings.TrimSpace(line[:idx])
			if processedLine != "" {
				resultLines = append(resultLines, processedLine)
			}
		} else {
			resultLines = append(resultLines, trimmed)
		}
	}
	return resultLines
}
