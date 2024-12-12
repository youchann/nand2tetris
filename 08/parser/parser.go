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
	c := token.CommandSymbol(strings.Fields(p.commandStrList[p.currentIndex])[0])
	switch c {
	case token.PUSH:
		return token.C_PUSH
	case token.POP:
		return token.C_POP
	case token.LABEL:
		return token.C_LABEL
	case token.GOTO:
		return token.C_GOTO
	case token.IF_GOTO:
		return token.C_IF
	case token.FUNCTION:
		return token.C_FUNCTION
	case token.RETURN:
		return token.C_RETURN
	case token.CALL:
		return token.C_CALL
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
