package tokenizer

import "strings"

type JackTokenizer struct {
	Input string
	// position     int
	// readPosition int
	// ch           byte
}

func New(input string) *JackTokenizer {
	return &JackTokenizer{Input: preprocessCode(input)}
}

func preprocessCode(input string) string {
	var result strings.Builder
	i := 0
	for i < len(input) {
		// remove multi-line comments
		if i+1 < len(input) && input[i:i+2] == "/*" {
			i += 2
			for i < len(input) {
				if i+1 < len(input) && input[i:i+2] == "*/" {
					i += 2
					break
				}
				i++
			}
			continue
		}

		// remove single-line comments
		if i+1 < len(input) && input[i:i+2] == "//" {
			i += 2
			for i < len(input) && input[i] != '\n' {
				i++
			}
			continue
		}

		// remove whitespace
		if input[i] != ' ' && input[i] != '\t' && input[i] != '\n' && input[i] != '\r' {
			result.WriteByte(input[i])
		}
		i++
	}
	return result.String()
}
