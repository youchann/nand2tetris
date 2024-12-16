package token

type Token struct {
	Type    TokenType
	Literal string
}

func (token *Token) Xml() string {
	switch token.Type {
	case KEYWORD:
		return "<keyword> " + token.Literal + " </keyword>"
	case SYMBOL:
		if token.Literal == ">" {
			return "<symbol> &gt; </symbol>"
		} else if token.Literal == "<" {
			return "<symbol> &lt; </symbol>"
		} else if token.Literal == "&" {
			return "<symbol> &amp; </symbol>"
		}
		return "<symbol> " + token.Literal + " </symbol>"
	case IDENTIFIER:
		return "<identifier> " + token.Literal + " </identifier>"
	case INT_CONST:
		return "<integerConstant> " + token.Literal + " </integerConstant>"
	case STRING_CONST:
		return "<stringConstant> " + token.Literal + " </stringConstant>"
	default:
		return ""
	}
}

type TokenType string

const (
	KEYWORD      TokenType = "KEYWORD"
	SYMBOL       TokenType = "SYMBOL"
	IDENTIFIER   TokenType = "IDENTIFIER"
	INT_CONST    TokenType = "INT_CONST"
	STRING_CONST TokenType = "STRING_CONST"
)

type Keyword string

const (
	CLASS       Keyword = "class"
	METHOD      Keyword = "method"
	FUNCTION    Keyword = "function"
	CONSTRUCTOR Keyword = "constructor"
	INT         Keyword = "int"
	BOOLEAN     Keyword = "boolean"
	CHAR        Keyword = "char"
	VOID        Keyword = "void"
	VAR         Keyword = "var"
	STATIC      Keyword = "static"
	FIELD       Keyword = "field"
	LET         Keyword = "let"
	DO          Keyword = "do"
	IF          Keyword = "if"
	ELSE        Keyword = "else"
	WHILE       Keyword = "while"
	RETURN      Keyword = "return"
	TRUE        Keyword = "true"
	FALSE       Keyword = "false"
	NULL        Keyword = "null"
	THIS        Keyword = "this"
)

var KeyWordMap = map[string]Keyword{
	"class":       CLASS,
	"method":      METHOD,
	"function":    FUNCTION,
	"constructor": CONSTRUCTOR,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

type Symbol string

const (
	LEFT_CURLY_BRACKET   Symbol = "{"
	RIGHT_CURLY_BRACKET  Symbol = "}"
	LEFT_ROUND_BRACKET   Symbol = "("
	RIGHT_ROUND_BRACKET  Symbol = ")"
	LEFT_SQUARE_BRACKET  Symbol = "["
	RIGHT_SQUARE_BRACKET Symbol = "]"
	PERIOD               Symbol = "."
	COMMA                Symbol = ","
	SEMICOLON            Symbol = ";"
	PLUS                 Symbol = "+"
	MINUS                Symbol = "-"
	ASTERISK             Symbol = "*"
	SLASH                Symbol = "/"
	AND                  Symbol = "&"
	PIPE                 Symbol = "|"
	LESS_THAN            Symbol = "<"
	GREATER_THAN         Symbol = ">"
	EQUAL                Symbol = "="
	TILDE                Symbol = "~"
)
