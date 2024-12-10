package parser

import (
	"strings"
)

type instructionType string

const (
	A_INSTRUCTION instructionType = "A_INSTRUCTION" // @Xxx
	C_INSTRUCTION instructionType = "C_INSTRUCTION" // dest=comp;jump
	L_INSTRUCTION instructionType = "L_INSTRUCTION" // (Xxx
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

func (p *Parser) CommandType() instructionType {
	commandStr := p.commandStrList[p.currentIndex]
	switch commandStr[0] {
	case '@':
		return A_INSTRUCTION
	case '0', '1', 'D', 'A', '!', '-', 'M':
		return C_INSTRUCTION
	case '(':
		return L_INSTRUCTION
	default:
		return ""
	}
}

func (p *Parser) Symbol() string {
	commandStr := p.commandStrList[p.currentIndex]
	switch p.CommandType() {
	case A_INSTRUCTION:
		return strings.TrimLeft(commandStr, "@")
	case L_INSTRUCTION:
		return strings.TrimRight(strings.TrimLeft(commandStr, "("), ")")
	default:
		return ""
	}
}

func (p *Parser) Dest() string {
	commandStr := p.commandStrList[p.currentIndex]
	if strings.Contains(commandStr, "=") {
		return strings.Split(commandStr, "=")[0]
	}
	return ""
}

func (p *Parser) Comp() string {
	commandStr := p.commandStrList[p.currentIndex]
	if strings.Contains(commandStr, "=") {
		return strings.Split(strings.Split(commandStr, "=")[1], ";")[0]
	}
	return strings.Split(commandStr, ";")[0]
}

func (p *Parser) Jump() string {
	commandStr := p.commandStrList[p.currentIndex]
	if strings.Contains(commandStr, ";") {
		return strings.Split(commandStr, ";")[1]
	}
	return ""
}

// TODO: 構文として正しいかどうかのチェックを加える
func preprocessCode(input string) []string {
	lines := strings.Split(input, "\n")
	return removeSpaces(removeEmptyLines(removeComments(lines)))
}

func removeSpaces(lines []string) []string {
	var processedLines []string
	for _, line := range lines {
		processedLine := strings.ReplaceAll(line, " ", "")
		processedLines = append(processedLines, processedLine)
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
