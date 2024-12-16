package tokenizer

import (
	"strconv"
	"strings"

	"github.com/youchann/nand2tetris/10-1_tokenizer/token"
)

type JackTokenizer struct {
	input        string
	currentToken *token.Token
	nextPosition int
}

func New(input string) *JackTokenizer {
	t := &JackTokenizer{input: preprocessCode(input)}
	t.Advance()
	return t
}

// TODO: いらないかも
func (t *JackTokenizer) CurrentToken() *token.Token {
	return t.currentToken
}
func (t *JackTokenizer) Input() string {
	return t.input
}

func (t *JackTokenizer) HasMoreTokens() bool {
	return t.nextPosition < len(t.input)
}

func (t *JackTokenizer) Advance() {
	// skip whitespace
	for t.nextPosition < len(t.input) && (t.input[t.nextPosition] == ' ' || t.input[t.nextPosition] == '\n' || t.input[t.nextPosition] == '\r' || t.input[t.nextPosition] == '\t') {
		t.nextPosition++
		if t.nextPosition == len(t.input) {
			return
		}
	}

	// keyword or identifier
	if isLetter(t.input[t.nextPosition]) {
		start := t.nextPosition
		for t.nextPosition < len(t.input) && (isLetter(t.input[t.nextPosition]) || isDigit(t.input[t.nextPosition])) {
			t.nextPosition++
		}
		if _, exists := token.KeyWordMap[t.input[start:t.nextPosition]]; exists {
			t.currentToken = &token.Token{
				Type:    token.KEYWORD,
				Literal: t.input[start:t.nextPosition],
			}
		} else {
			t.currentToken = &token.Token{
				Type:    token.IDENTIFIER,
				Literal: t.input[start:t.nextPosition],
			}
		}
		return
	}

	// symbol
	if isSymbol(t.input[t.nextPosition]) {
		t.currentToken = &token.Token{
			Type:    token.SYMBOL,
			Literal: string(t.input[t.nextPosition]),
		}
		t.nextPosition++
		return
	}

	// integer constant
	if isDigit(t.input[t.nextPosition]) {
		start := t.nextPosition
		for t.nextPosition < len(t.input) && isDigit(t.input[t.nextPosition]) {
			t.nextPosition++
		}
		t.currentToken = &token.Token{
			Type:    token.INT_CONST,
			Literal: t.input[start:t.nextPosition],
		}
		return
	}

	// string constant
	if t.input[t.nextPosition] == '"' {
		start := t.nextPosition + 1
		t.nextPosition++
		for t.nextPosition < len(t.input) && t.input[t.nextPosition] != '"' {
			t.nextPosition++
		}
		t.currentToken = &token.Token{
			Type:    token.STRING_CONST,
			Literal: t.input[start:t.nextPosition],
		}
		t.nextPosition++
		return
	}
}

func (t *JackTokenizer) TokenType() token.TokenType {
	return t.currentToken.Type
}

func (t *JackTokenizer) Keyword() token.Keyword {
	return token.Keyword(t.currentToken.Literal)
}

func (t *JackTokenizer) Symbol() token.Symbol {
	return token.Symbol(t.currentToken.Literal)
}

func (t *JackTokenizer) Identifier() string {
	return t.currentToken.Literal
}

func (t *JackTokenizer) IntVal() int {
	v, _ := strconv.Atoi(t.currentToken.Literal) // TODO: handle error
	return v
}

func (t *JackTokenizer) StringVal() string {
	return t.currentToken.Literal
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(ch byte) bool {
	return strings.Contains("{}()[].,;+-*/&|<>=~", string(ch))
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
		result.WriteByte(input[i])
		i++
	}
	return result.String()
}
