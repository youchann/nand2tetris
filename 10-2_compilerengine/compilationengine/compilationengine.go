package compilationengine

import (
	"slices"

	"github.com/youchann/nand2tetris/10-2_compilerengine/token"
	"github.com/youchann/nand2tetris/10-2_compilerengine/tokenizer"
)

type CompilationEngine struct {
	tokenizer *tokenizer.JackTokenizer
	indent    int
	XML       string
}

func New(t *tokenizer.JackTokenizer) *CompilationEngine {
	return &CompilationEngine{
		tokenizer: t,
		indent:    0,
		XML:       "",
	}
}

func (ce *CompilationEngine) CompileClass() {
	ce.print("<class>")
	ce.indent++

	ce.process("class")

	// className
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()

	ce.process("{")
	ce.CompileClassVarDec()
	ce.CompileSubroutine()
	ce.process("}")

	ce.indent--
	ce.print("</class>")
}

func (ce *CompilationEngine) CompileClassVarDec() {
	for ce.tokenizer.CurrentToken().Literal == "static" || ce.tokenizer.CurrentToken().Literal == "field" {
		ce.print("<classVarDec>")
		ce.indent++

		// static or field
		ce.process(ce.tokenizer.CurrentToken().Literal)
		// type
		ce.processType()

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()
			if ce.tokenizer.CurrentToken().Literal == "," {
				ce.process(",")
			}
		}

		ce.process(";")

		ce.indent--
		ce.print("</classVarDec>")
	}
}

func (ce *CompilationEngine) CompileSubroutine() {
	subroutineType := []token.Keyword{token.CONSTRUCTOR, token.FUNCTION, token.METHOD}
	for slices.Contains(subroutineType, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		ce.print("<subroutineDec>")
		ce.indent++

		// constructor, function, or method
		ce.process(ce.tokenizer.CurrentToken().Literal)

		// void or type
		voidOrType := []token.Keyword{token.VOID, token.INT, token.CHAR, token.BOOLEAN}
		if !slices.Contains(voidOrType, token.Keyword(ce.tokenizer.CurrentToken().Literal)) && ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected type or void but got " + ce.tokenizer.CurrentToken().Literal)
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		// subroutineName
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		ce.process("(")
		ce.CompileParameterList()
		ce.process(")")
		ce.CompileSubroutineBody()

		ce.indent--
		ce.print("</subroutineDec>")
	}
}

func (ce *CompilationEngine) CompileParameterList() {
	ce.print("<parameterList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		// type
		ce.processType()

		// varName
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "," {
			ce.process(",")
		}
	}

	ce.indent--
	ce.print("</parameterList>")
}

func (ce *CompilationEngine) CompileSubroutineBody() {
	ce.print("<subroutineBody>")
	ce.indent++

	ce.process("{")
	ce.CompileVarDec()
	ce.CompileStatements()
	ce.process("}")

	ce.indent--
	ce.print("</subroutineBody>")
}

func (ce *CompilationEngine) CompileVarDec() {
	for ce.tokenizer.CurrentToken().Literal == "var" {
		ce.print("<varDec>")
		ce.indent++

		ce.process("var")
		ce.processType()

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()

			if ce.tokenizer.CurrentToken().Literal == "," {
				ce.process(",")
			}
		}

		ce.process(";")

		ce.indent--
		ce.print("</varDec>")
	}
}

func (ce *CompilationEngine) CompileStatements() {
	ce.print("<statements>")
	ce.indent++

	statementPrefix := []token.Keyword{token.LET, token.IF, token.WHILE, token.DO, token.RETURN}
	for slices.Contains(statementPrefix, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		switch token.Keyword(ce.tokenizer.CurrentToken().Literal) {
		case token.LET:
			ce.CompileLet()
		case token.IF:
			ce.CompileIf()
		case token.WHILE:
			ce.CompileWhile()
		case token.DO:
			ce.CompileDo()
		case token.RETURN:
			ce.CompileReturn()
		}
	}

	ce.indent--
	ce.print("</statements>")
}

func (ce *CompilationEngine) CompileLet() {
	ce.print("<letStatement>")
	ce.indent++

	ce.process("let")

	// varName
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()

	if ce.tokenizer.CurrentToken().Literal == "[" {
		ce.process("[")
		ce.CompileExpression()
		ce.process("]")
	}

	ce.process("=")
	ce.CompileExpression()
	ce.process(";")

	ce.indent--
	ce.print("</letStatement>")
}

func (ce *CompilationEngine) CompileIf() {
	ce.print("<ifStatement>")
	ce.indent++

	ce.process("if")
	ce.process("(")
	ce.CompileExpression()
	ce.process(")")
	ce.process("{")
	ce.CompileStatements()
	ce.process("}")
	if ce.tokenizer.CurrentToken().Literal == "else" {
		ce.process("else")
		ce.process("{")
		ce.CompileStatements()
		ce.process("}")
	}

	ce.indent--
	ce.print("</ifStatement>")
}

func (ce *CompilationEngine) CompileWhile() {
	ce.print("<whileStatement>")
	ce.indent++

	ce.process("while")
	ce.process("(")
	ce.CompileExpression()
	ce.process(")")
	ce.process("{")
	ce.CompileStatements()
	ce.process("}")

	ce.indent--
	ce.print("</whileStatement>")
}

func (ce *CompilationEngine) CompileDo() {
	ce.print("<doStatement>")
	ce.indent++

	ce.process("do")

	// subroutineCall
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
	if ce.tokenizer.CurrentToken().Literal == "." {
		ce.process(".")
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	}

	ce.process("(")
	ce.CompileExpressionList()
	ce.process(")")
	ce.process(";")

	ce.indent--
	ce.print("</doStatement>")
}

func (ce *CompilationEngine) CompileReturn() {
	ce.print("<returnStatement>")
	ce.indent++

	ce.process("return")
	if ce.tokenizer.CurrentToken().Literal != ";" {
		ce.CompileExpression()
	}
	ce.process(";")

	ce.indent--
	ce.print("</returnStatement>")
}

func (ce *CompilationEngine) CompileExpression() {
	ce.print("<expression>")
	ce.indent++

	ce.CompileTerm()

	operand := []token.Symbol{token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.AND, token.PIPE, token.LESS_THAN, token.GREATER_THAN, token.EQUAL}
	for slices.Contains(operand, token.Symbol(ce.tokenizer.CurrentToken().Literal)) {
		ce.process(ce.tokenizer.CurrentToken().Literal)
		ce.CompileTerm()
	}

	ce.indent--
	ce.print("</expression>")
}

func (ce *CompilationEngine) CompileTerm() {
	ce.print("<term>")
	ce.indent++

	constants := []token.TokenType{token.INT_CONST, token.STRING_CONST}
	keywordConstants := []token.Keyword{token.TRUE, token.FALSE, token.NULL, token.THIS}
	if slices.Contains(constants, ce.tokenizer.CurrentToken().Type) || slices.Contains(keywordConstants, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	} else if ce.tokenizer.CurrentToken().Literal == "(" {
		ce.process("(")
		ce.CompileExpression()
		ce.process(")")
	} else if ce.tokenizer.CurrentToken().Literal == "-" || ce.tokenizer.CurrentToken().Literal == "~" {
		ce.process(ce.tokenizer.CurrentToken().Literal)
		ce.CompileTerm()
	} else {
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "[" {
			ce.process("[")
			ce.CompileExpression()
			ce.process("]")
		} else if ce.tokenizer.CurrentToken().Literal == "." {
			ce.process(".")
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + ce.tokenizer.CurrentToken().Literal)
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()

			ce.process("(")
			ce.CompileExpressionList()
			ce.process(")")
		} else if ce.tokenizer.CurrentToken().Literal == "(" {
			ce.process("(")
			ce.CompileExpressionList()
			ce.process(")")
		}
	}

	ce.indent--
	ce.print("</term>")
}

func (ce *CompilationEngine) CompileExpressionList() int {
	count := 0
	ce.print("<expressionList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		ce.CompileExpression()
		count++
		if ce.tokenizer.CurrentToken().Literal == "," {
			ce.process(",")
		}
	}

	ce.indent--
	ce.print("</expressionList>")
	return count
}

func (ce *CompilationEngine) process(str string) {
	if ce.tokenizer.CurrentToken().Literal != str {
		panic("expected " + str + " but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
}

func (ce *CompilationEngine) processType() {
	types := []token.Keyword{token.INT, token.CHAR, token.BOOLEAN}
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && !slices.Contains(types, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		panic("expected type but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
}

func (ce *CompilationEngine) print(str string) {
	indentation := ""
	for i := 0; i < ce.indent; i++ {
		indentation += "  "
	}
	ce.XML += indentation + str + "\n"
}
