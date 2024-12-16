package compilationengine

import (
	"errors"

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

func (ce *CompilationEngine) CompileClass() error {
	if ce.tokenizer.CurrentToken().Literal != "class" {
		return errors.New("expected 'class' keyword")
	}
	ce.print("<class>")
	ce.indent++

	if err := ce.process("class"); err != nil {
		return err
	}

	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		return errors.New("expected identifier")
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()

	if err := ce.process("{"); err != nil {
		return err
	}

	if err := ce.CompileClassVarDec(); err != nil {
		return err
	}
	if err := ce.CompileSubroutine(); err != nil {
		return err
	}

	if err := ce.process("}"); err != nil {
		return err
	}

	ce.indent--
	ce.print("</class>")
	return nil
}

func (ce *CompilationEngine) CompileClassVarDec() error {
	for ce.tokenizer.CurrentToken().Literal == "static" || ce.tokenizer.CurrentToken().Literal == "field" {
		ce.print("<classVarDec>")
		ce.indent++

		// static or field
		if err := ce.process(ce.tokenizer.CurrentToken().Literal); err != nil {
			return err
		}

		// type
		if err := ce.processType(); err != nil {
			return err
		}

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				return errors.New("expected identifier")
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()

			if ce.tokenizer.CurrentToken().Literal == "," {
				if err := ce.process(","); err != nil {
					return err
				}
			}
		}

		if err := ce.process(";"); err != nil {
			return err
		}

		ce.indent--
		ce.print("</classVarDec>")
	}
	return nil
}

func (ce *CompilationEngine) CompileSubroutine() error {
	for ce.tokenizer.CurrentToken().Literal == "constructor" || ce.tokenizer.CurrentToken().Literal == "function" || ce.tokenizer.CurrentToken().Literal == "method" {
		ce.print("<subroutineDec>")
		ce.indent++

		// constructor, function, or method
		if err := ce.process(ce.tokenizer.CurrentToken().Literal); err != nil {
			return err
		}

		// void or type
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && ce.tokenizer.CurrentToken().Literal != "void" && ce.tokenizer.CurrentToken().Literal != "int" && ce.tokenizer.CurrentToken().Literal != "char" && ce.tokenizer.CurrentToken().Literal != "boolean" {
			return errors.New("expected type")
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		// subroutineName
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			return errors.New("expected identifier")
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		if err := ce.process("("); err != nil {
			return err
		}

		if err := ce.CompileParameterList(); err != nil {
			return err
		}

		if err := ce.process(")"); err != nil {
			return err
		}

		if err := ce.CompileSubroutineBody(); err != nil {
			return err
		}

		ce.indent--
		ce.print("</subroutineDec>")
	}
	return nil
}

func (ce *CompilationEngine) CompileParameterList() error {
	ce.print("<parameterList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		// type
		if err := ce.processType(); err != nil {
			return err
		}

		// varName
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			return errors.New("expected identifier")
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "," {
			if err := ce.process(","); err != nil {
				return err
			}
		}
	}

	ce.indent--
	ce.print("</parameterList>")
	return nil
}

func (ce *CompilationEngine) CompileSubroutineBody() error {
	ce.print("<subroutineBody>")
	ce.indent++
	if err := ce.process("{"); err != nil {
		return err
	}

	if err := ce.CompileVarDec(); err != nil {
		return err
	}

	if err := ce.CompileStatements(); err != nil {
		return err
	}

	if err := ce.process("}"); err != nil {
		return err
	}
	ce.indent--
	ce.print("</subroutineBody>")
	return nil
}

func (ce *CompilationEngine) CompileVarDec() error {
	for ce.tokenizer.CurrentToken().Literal == "var" {
		ce.print("<varDec>")
		ce.indent++

		if err := ce.process("var"); err != nil {
			return err
		}

		// type
		if err := ce.processType(); err != nil {
			return err
		}

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				return errors.New("expected identifier")
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()

			if ce.tokenizer.CurrentToken().Literal == "," {
				if err := ce.process(","); err != nil {
					return err
				}
			}
		}

		if err := ce.process(";"); err != nil {
			return err
		}

		ce.indent--
		ce.print("</varDec>")
	}
	return nil
}

func (ce *CompilationEngine) CompileStatements() error {
	ce.print("<statements>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal == string(token.LET) || ce.tokenizer.CurrentToken().Literal == string(token.IF) || ce.tokenizer.CurrentToken().Literal == string(token.WHILE) || ce.tokenizer.CurrentToken().Literal == string(token.DO) || ce.tokenizer.CurrentToken().Literal == string(token.RETURN) {
		var err error
		switch ce.tokenizer.CurrentToken().Literal {
		case string(token.LET):
			err = ce.CompileLet()
		case string(token.IF):
			err = ce.CompileIf()
		case string(token.WHILE):
			err = ce.CompileWhile()
		case string(token.DO):
			err = ce.CompileDo()
		case string(token.RETURN):
			err = ce.CompileReturn()
		}
		if err != nil {
			return err
		}
	}

	ce.indent--
	ce.print("</statements>")
	return nil
}

func (ce *CompilationEngine) CompileLet() error {
	ce.print("<letStatement>")
	ce.indent++

	if err := ce.process("let"); err != nil {
		return err
	}

	// varName
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		return errors.New("expected identifier")
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()

	if ce.tokenizer.CurrentToken().Literal == "[" {
		if err := ce.process("["); err != nil {
			return err
		}

		if err := ce.CompileExpression(); err != nil {
			return err
		}

		if err := ce.process("]"); err != nil {
			return err
		}
	}

	if err := ce.process("="); err != nil {
		return err
	}
	if err := ce.CompileExpression(); err != nil {
		return err
	}
	if err := ce.process(";"); err != nil {
		return err
	}

	ce.indent--
	ce.print("</letStatement>")
	return nil
}

func (ce *CompilationEngine) CompileIf() error {
	ce.print("<ifStatement>")
	ce.indent++

	if err := ce.process("if"); err != nil {
		return err
	}
	if err := ce.process("("); err != nil {
		return err
	}
	if err := ce.CompileExpression(); err != nil {
		return err
	}
	if err := ce.process(")"); err != nil {
		return err
	}
	if err := ce.process("{"); err != nil {
		return err
	}
	if err := ce.CompileStatements(); err != nil {
		return err
	}
	if err := ce.process("}"); err != nil {
		return err
	}
	if ce.tokenizer.CurrentToken().Literal == "else" {
		if err := ce.process("else"); err != nil {
			return err
		}
		if err := ce.process("{"); err != nil {
			return err
		}
		if err := ce.CompileStatements(); err != nil {
			return err
		}
		if err := ce.process("}"); err != nil {
			return err
		}
	}

	ce.indent--
	ce.print("</ifStatement>")
	return nil
}

func (ce *CompilationEngine) CompileWhile() error {
	ce.print("<whileStatement>")
	ce.indent++

	if err := ce.process("while"); err != nil {
		return err
	}
	if err := ce.process("("); err != nil {
		return err
	}
	if err := ce.CompileExpression(); err != nil {
		return err
	}
	if err := ce.process(")"); err != nil {
		return err
	}
	if err := ce.process("{"); err != nil {
		return err
	}
	if err := ce.CompileStatements(); err != nil {
		return err
	}
	if err := ce.process("}"); err != nil {
		return err
	}

	ce.indent--
	ce.print("</whileStatement>")
	return nil
}

func (ce *CompilationEngine) CompileDo() error {
	ce.print("<doStatement>")
	ce.indent++

	if err := ce.process("do"); err != nil {
		return err
	}

	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		return errors.New("expected identifier")
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()

	if ce.tokenizer.CurrentToken().Literal == "." {
		if err := ce.process("."); err != nil {
			return err
		}

		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			return errors.New("expected identifier")
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	}

	if err := ce.process("("); err != nil {
		return err
	}
	if _, err := ce.CompileExpressionList(); err != nil {
		return err
	}
	if err := ce.process(")"); err != nil {
		return err
	}
	if err := ce.process(";"); err != nil {
		return err
	}

	ce.indent--
	ce.print("</doStatement>")
	return nil
}

func (ce *CompilationEngine) CompileReturn() error {
	ce.print("<returnStatement>")
	ce.indent++

	if err := ce.process("return"); err != nil {
		return err
	}

	if ce.tokenizer.CurrentToken().Literal != ";" {
		if err := ce.CompileExpression(); err != nil {
			return err
		}
	}

	if err := ce.process(";"); err != nil {
		return err
	}

	ce.indent--
	ce.print("</returnStatement>")
	return nil
}

func (ce *CompilationEngine) CompileExpression() error {
	ce.print("<expression>")
	ce.indent++

	if err := ce.CompileTerm(); err != nil {
		return err
	}

	for ce.tokenizer.CurrentToken().Literal == "+" || ce.tokenizer.CurrentToken().Literal == "-" || ce.tokenizer.CurrentToken().Literal == "*" || ce.tokenizer.CurrentToken().Literal == "/" || ce.tokenizer.CurrentToken().Literal == "&" || ce.tokenizer.CurrentToken().Literal == "|" || ce.tokenizer.CurrentToken().Literal == "<" || ce.tokenizer.CurrentToken().Literal == ">" || ce.tokenizer.CurrentToken().Literal == "=" {
		if err := ce.process(ce.tokenizer.CurrentToken().Literal); err != nil {
			return err
		}

		if err := ce.CompileTerm(); err != nil {
			return err
		}
	}

	ce.indent--
	ce.print("</expression>")
	return nil
}

func (ce *CompilationEngine) CompileTerm() error {
	ce.print("<term>")
	ce.indent++

	if ce.tokenizer.CurrentToken().Type == token.INT_CONST || ce.tokenizer.CurrentToken().Type == token.STRING_CONST || ce.tokenizer.CurrentToken().Literal == string(token.TRUE) || ce.tokenizer.CurrentToken().Literal == string(token.FALSE) || ce.tokenizer.CurrentToken().Literal == string(token.NULL) || ce.tokenizer.CurrentToken().Literal == string(token.THIS) {
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	} else if ce.tokenizer.CurrentToken().Literal == "(" {
		if err := ce.process("("); err != nil {
			return err
		}

		if err := ce.CompileExpression(); err != nil {
			return err
		}

		if err := ce.process(")"); err != nil {
			return err
		}
	} else if ce.tokenizer.CurrentToken().Literal == "-" || ce.tokenizer.CurrentToken().Literal == "~" {
		if err := ce.process(ce.tokenizer.CurrentToken().Literal); err != nil {
			return err
		}

		if err := ce.CompileTerm(); err != nil {
			return err
		}
	} else {
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			return errors.New("expected identifier")
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "[" {
			if err := ce.process("["); err != nil {
				return err
			}

			if err := ce.CompileExpression(); err != nil {
				return err
			}

			if err := ce.process("]"); err != nil {
				return err
			}
		} else if ce.tokenizer.CurrentToken().Literal == "." {
			if err := ce.process("."); err != nil {
				return err
			}

			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				return errors.New("expected identifier")
			}
			ce.print(ce.tokenizer.CurrentToken().Xml())
			ce.tokenizer.Advance()

			if err := ce.process("("); err != nil {
				return err
			}

			if _, err := ce.CompileExpressionList(); err != nil {
				return err
			}

			if err := ce.process(")"); err != nil {
				return err
			}
		} else if ce.tokenizer.CurrentToken().Literal == "(" {
			if err := ce.process("("); err != nil {
				return err
			}

			if _, err := ce.CompileExpressionList(); err != nil {
				return err
			}

			if err := ce.process(")"); err != nil {
				return err
			}
		}
	}

	ce.indent--
	ce.print("</term>")
	return nil
}

func (ce *CompilationEngine) CompileExpressionList() (int, error) {
	count := 0
	ce.print("<expressionList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		if err := ce.CompileExpression(); err != nil {
			return 0, err
		}

		count++

		if ce.tokenizer.CurrentToken().Literal == "," {
			if err := ce.process(","); err != nil {
				return 0, err
			}
		}
	}

	ce.indent--
	ce.print("</expressionList>")
	return count, nil
}

func (ce *CompilationEngine) process(str string) error {
	if ce.tokenizer.CurrentToken().Literal == str {
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	} else {
		return errors.New("expected " + str + " but got " + ce.tokenizer.CurrentToken().Literal)
	}
	return nil
}

func (ce *CompilationEngine) processType() error {
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && ce.tokenizer.CurrentToken().Literal != "int" && ce.tokenizer.CurrentToken().Literal != "char" && ce.tokenizer.CurrentToken().Literal != "boolean" {
		return errors.New("expected type")
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
	return nil
}

func (ce *CompilationEngine) print(str string) {
	indentation := ""
	for i := 0; i < ce.indent; i++ {
		indentation += "  "
	}
	ce.XML += indentation + str + "\n"
}
