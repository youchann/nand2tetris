package compilationengine

import (
	"github.com/youchann/nand2tetris/09/token"
	"github.com/youchann/nand2tetris/09/tokenizer"
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
	for ce.tokenizer.CurrentToken().Literal == "constructor" || ce.tokenizer.CurrentToken().Literal == "function" || ce.tokenizer.CurrentToken().Literal == "method" {
		ce.print("<subroutineDec>")
		ce.indent++

		// constructor, function, or method
		ce.process(ce.tokenizer.CurrentToken().Literal)

		// void or type
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && ce.tokenizer.CurrentToken().Literal != "void" && ce.tokenizer.CurrentToken().Literal != "int" && ce.tokenizer.CurrentToken().Literal != "char" && ce.tokenizer.CurrentToken().Literal != "boolean" {
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

	for ce.tokenizer.CurrentToken().Literal == string(token.LET) || ce.tokenizer.CurrentToken().Literal == string(token.IF) || ce.tokenizer.CurrentToken().Literal == string(token.WHILE) || ce.tokenizer.CurrentToken().Literal == string(token.DO) || ce.tokenizer.CurrentToken().Literal == string(token.RETURN) {
		switch ce.tokenizer.CurrentToken().Literal {
		case string(token.LET):
			ce.CompileLet()
		case string(token.IF):
			ce.CompileIf()
		case string(token.WHILE):
			ce.CompileWhile()
		case string(token.DO):
			ce.CompileDo()
		case string(token.RETURN):
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

	for ce.tokenizer.CurrentToken().Literal == "+" || ce.tokenizer.CurrentToken().Literal == "-" || ce.tokenizer.CurrentToken().Literal == "*" || ce.tokenizer.CurrentToken().Literal == "/" || ce.tokenizer.CurrentToken().Literal == "&" || ce.tokenizer.CurrentToken().Literal == "|" || ce.tokenizer.CurrentToken().Literal == "<" || ce.tokenizer.CurrentToken().Literal == ">" || ce.tokenizer.CurrentToken().Literal == "=" {
		ce.process(ce.tokenizer.CurrentToken().Literal)
		ce.CompileTerm()
	}

	ce.indent--
	ce.print("</expression>")
}

func (ce *CompilationEngine) CompileTerm() {
	ce.print("<term>")
	ce.indent++

	if ce.tokenizer.CurrentToken().Type == token.INT_CONST || ce.tokenizer.CurrentToken().Type == token.STRING_CONST || ce.tokenizer.CurrentToken().Literal == string(token.TRUE) || ce.tokenizer.CurrentToken().Literal == string(token.FALSE) || ce.tokenizer.CurrentToken().Literal == string(token.NULL) || ce.tokenizer.CurrentToken().Literal == string(token.THIS) {
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
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && ce.tokenizer.CurrentToken().Literal != "int" && ce.tokenizer.CurrentToken().Literal != "char" && ce.tokenizer.CurrentToken().Literal != "boolean" {
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
