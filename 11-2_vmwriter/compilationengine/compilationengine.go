package compilationengine

import (
	"slices"
	"strconv"

	"github.com/youchann/nand2tetris/11-2_vmwriter/symboltable"
	"github.com/youchann/nand2tetris/11-2_vmwriter/token"
	"github.com/youchann/nand2tetris/11-2_vmwriter/tokenizer"
	"github.com/youchann/nand2tetris/11-2_vmwriter/vmwriter"
)

var kindSegmentMap = map[symboltable.Kind]vmwriter.Segment{
	symboltable.STATIC:    vmwriter.STATIC,
	symboltable.FIELD:     vmwriter.THIS, // TODO: need to check
	symboltable.ARGUMENT:  vmwriter.ARGUMENT,
	symboltable.VAR_LOCAL: vmwriter.LOCAL,
}

type CompilationEngine struct {
	tokenizer    *tokenizer.JackTokenizer
	vmwriter     *vmwriter.VMWriter
	indent       int
	XML          string
	className    string
	classST      *symboltable.SymbolTable
	subroutineST *symboltable.SymbolTable
}

func New(n string, t *tokenizer.JackTokenizer, w *vmwriter.VMWriter) *CompilationEngine {
	return &CompilationEngine{
		tokenizer:    t,
		vmwriter:     w,
		indent:       0,
		XML:          "",
		className:    n,
		classST:      symboltable.New(),
		subroutineST: symboltable.New(),
	}
}

func (ce *CompilationEngine) CompileClass() {
	ce.print("<class>")
	ce.indent++

	ce.process("class")

	// className
	name := ce.tokenizer.CurrentToken().Literal
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + name)
	} else if name != ce.className {
		panic("class name does not match file name")
	}
	ce.print("<identifier> " + "name: " + name + ", category: class, index: -1, usage: definition" + " </identifier>")
	ce.tokenizer.Advance()

	ce.process("{")
	ce.compileClassVarDec()
	ce.compileSubroutine()
	ce.process("}")

	ce.indent--
	ce.print("</class>")
}

func (ce *CompilationEngine) compileClassVarDec() {
	for ce.tokenizer.CurrentToken().Literal == "static" || ce.tokenizer.CurrentToken().Literal == "field" {
		ce.print("<classVarDec>")
		ce.indent++

		// static or field
		kind := ce.process(ce.tokenizer.CurrentToken().Literal)
		// type
		typ := ce.processType()

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			name := ce.tokenizer.CurrentToken().Literal
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + name)
			}
			ce.classST.Define(name, typ, symboltable.KindMap[kind])
			ce.print("<identifier> " + "name: " + name + ", category: " + kind + ", index: " + strconv.Itoa(ce.classST.IndexOf(name)) + ", usage: definition" + " </identifier>")
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

func (ce *CompilationEngine) compileSubroutine() {
	subroutineType := []token.Keyword{token.CONSTRUCTOR, token.FUNCTION, token.METHOD}
	for slices.Contains(subroutineType, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		ce.subroutineST.Reset()

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
		name := ce.tokenizer.CurrentToken().Literal
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + name)
		}
		ce.tokenizer.Advance()

		ce.process("(")
		ce.compileParameterList()
		ce.process(")")

		ce.process("{")
		n := ce.compileVarDec()
		ce.vmwriter.WriteFunction(name, n)
		ce.compileStatements()
		ce.process("}")
	}
}

func (ce *CompilationEngine) compileParameterList() {
	ce.print("<parameterList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		// type
		typ := ce.processType()

		// varName
		name := ce.tokenizer.CurrentToken().Literal
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + name)
		}
		ce.subroutineST.Define(name, typ, symboltable.ARGUMENT)
		ce.print("<identifier> " + "name: " + name + ", category: argument, index: " + strconv.Itoa(ce.subroutineST.IndexOf(name)) + ", usage: definition" + " </identifier>")
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "," {
			ce.process(",")
		}
	}

	ce.indent--
	ce.print("</parameterList>")
}

func (ce *CompilationEngine) compileVarDec() int {
	count := 0
	for ce.tokenizer.CurrentToken().Literal == "var" {
		ce.process("var")
		typ := ce.processType()

		// varName
		for ce.tokenizer.CurrentToken().Literal != ";" {
			name := ce.tokenizer.CurrentToken().Literal
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + name)
			}
			ce.subroutineST.Define(name, typ, symboltable.VAR_LOCAL)
			count++
			ce.tokenizer.Advance()

			if ce.tokenizer.CurrentToken().Literal == "," {
				ce.process(",")
			}
		}

		ce.process(";")
	}
	return count
}

func (ce *CompilationEngine) compileStatements() {
	ce.print("<statements>")
	ce.indent++

	statementPrefix := []token.Keyword{token.LET, token.IF, token.WHILE, token.DO, token.RETURN}
	for slices.Contains(statementPrefix, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		switch token.Keyword(ce.tokenizer.CurrentToken().Literal) {
		case token.LET:
			ce.compileLet()
		case token.IF:
			ce.compileIf()
		case token.WHILE:
			ce.compileWhile()
		case token.DO:
			ce.compileDo()
		case token.RETURN:
			ce.compileReturn()
		}
	}

	ce.indent--
	ce.print("</statements>")
}

func (ce *CompilationEngine) compileLet() {
	ce.print("<letStatement>")
	ce.indent++

	ce.process("let")

	// varName
	name := ce.tokenizer.CurrentToken().Literal
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + name)
	}
	if ce.subroutineST.IndexOf(name) != -1 {
		ce.print("<identifier> " + "name: " + name + ", category: var, index: " + strconv.Itoa(ce.subroutineST.IndexOf(name)) + ", usage: using" + " </identifier>")
	} else if ce.classST.IndexOf(name) != -1 {
		ce.print("<identifier> " + "name: " + name + ", category: class, index: " + strconv.Itoa(ce.classST.IndexOf(name)) + ", usage: using" + " </identifier>")
	} else {
		panic("undefined variable " + name)
	}
	ce.tokenizer.Advance()

	if ce.tokenizer.CurrentToken().Literal == "[" {
		ce.process("[")
		ce.compileExpression()
		ce.process("]")
	}

	ce.process("=")
	ce.compileExpression()
	ce.process(";")

	ce.vmwriter.WritePop(kindSegmentMap[ce.subroutineST.KindOf(name)], ce.subroutineST.IndexOf(name))

	ce.indent--
	ce.print("</letStatement>")
}

func (ce *CompilationEngine) compileIf() {
	ce.print("<ifStatement>")
	ce.indent++

	ce.process("if")
	ce.process("(")
	ce.compileExpression()
	ce.process(")")
	ce.process("{")
	ce.compileStatements()
	ce.process("}")
	if ce.tokenizer.CurrentToken().Literal == "else" {
		ce.process("else")
		ce.process("{")
		ce.compileStatements()
		ce.process("}")
	}

	ce.indent--
	ce.print("</ifStatement>")
}

func (ce *CompilationEngine) compileWhile() {
	ce.print("<whileStatement>")
	ce.indent++

	ce.process("while")
	ce.process("(")
	ce.compileExpression()
	ce.process(")")
	ce.process("{")
	ce.compileStatements()
	ce.process("}")

	ce.indent--
	ce.print("</whileStatement>")
}

func (ce *CompilationEngine) compileDo() {
	ce.print("<doStatement>")
	ce.indent++

	ce.process("do")

	// subroutineCall
	name := ce.tokenizer.CurrentToken().Literal
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
		panic("expected identifier but got " + name)
	}
	if ce.subroutineST.IndexOf(name) != -1 {
		ce.print("<identifier> " + "name: " + name + ", category: " + string(ce.subroutineST.KindOf(name)) + ", index: " + strconv.Itoa(ce.subroutineST.IndexOf(name)) + ", usage: using" + " </identifier>")
	} else if ce.classST.IndexOf(name) != -1 {
		ce.print("<identifier> " + "name: " + name + ", category: " + string(ce.classST.KindOf(name)) + ", index: " + strconv.Itoa(ce.classST.IndexOf(name)) + ", usage: using" + " </identifier>")
	} else {
		ce.print("<identifier> " + "name: " + name + ", category: class, index: -1, usage: using" + " </identifier>")
		// TODO: check VM API & class method
		// panic("undefined subroutine " + name)
	}
	ce.tokenizer.Advance()
	if ce.tokenizer.CurrentToken().Literal == "." {
		ce.process(".")
		n := ce.tokenizer.CurrentToken().Literal
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + n)
		}
		ce.print("<identifier> " + "name: " + n + ", category: subroutine, index: -1, usage: using" + " </identifier>")
		name += "." + n
		ce.tokenizer.Advance()
	}

	ce.process("(")
	c := ce.compileExpressionList()
	ce.process(")")
	ce.process(";")

	ce.vmwriter.WriteCall(name, c)
	ce.vmwriter.WritePop(vmwriter.TEMP, 0)

	ce.indent--
	ce.print("</doStatement>")
}

func (ce *CompilationEngine) compileReturn() {
	ce.print("<returnStatement>")
	ce.indent++

	ce.process("return")
	if ce.tokenizer.CurrentToken().Literal == ";" {
		ce.vmwriter.WritePush(vmwriter.CONSTANT, 0)
	} else {
		ce.compileExpression()
	}
	ce.vmwriter.WriteReturn()
	ce.process(";")

	ce.indent--
	ce.print("</returnStatement>")
}

func (ce *CompilationEngine) compileExpression() {
	ce.print("<expression>")
	ce.indent++

	ce.compileTerm()

	operand := []token.Symbol{token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.AND, token.PIPE, token.LESS_THAN, token.GREATER_THAN, token.EQUAL}
	for slices.Contains(operand, token.Symbol(ce.tokenizer.CurrentToken().Literal)) {
		op := ce.tokenizer.CurrentToken().Literal
		ce.process(op)
		ce.compileTerm()
		switch token.Symbol(op) {
		case token.PLUS:
			ce.vmwriter.WriteArithmetic(vmwriter.ADD)
		case token.MINUS:
			ce.vmwriter.WriteArithmetic(vmwriter.SUB)
		case token.ASTERISK:
			ce.vmwriter.WriteCall("Math.multiply", 2)
		case token.SLASH:
			ce.vmwriter.WriteCall("Math.divide", 2)
		case token.AND:
			ce.vmwriter.WriteArithmetic(vmwriter.AND)
		case token.PIPE:
			ce.vmwriter.WriteArithmetic(vmwriter.OR)
		case token.LESS_THAN:
			ce.vmwriter.WriteArithmetic(vmwriter.LT)
		case token.GREATER_THAN:
			ce.vmwriter.WriteArithmetic(vmwriter.GT)
		case token.EQUAL:
			ce.vmwriter.WriteArithmetic(vmwriter.EQ)
		}
	}

	ce.indent--
	ce.print("</expression>")
}

func (ce *CompilationEngine) compileTerm() {
	ce.print("<term>")
	ce.indent++

	constants := []token.TokenType{token.INT_CONST, token.STRING_CONST}
	keywordConstants := []token.Keyword{token.TRUE, token.FALSE, token.NULL, token.THIS}
	if slices.Contains(constants, ce.tokenizer.CurrentToken().Type) {
		switch ce.tokenizer.CurrentToken().Type {
		case token.INT_CONST:
			value, err := strconv.Atoi(ce.tokenizer.CurrentToken().Literal)
			if err != nil {
				panic("expected integer constant but got " + ce.tokenizer.CurrentToken().Literal)
			}
			ce.vmwriter.WritePush(vmwriter.CONSTANT, value)
		// TODO: need to check
		case token.STRING_CONST:
			ce.vmwriter.WritePush(vmwriter.CONSTANT, len(ce.tokenizer.CurrentToken().Literal))
			ce.vmwriter.WriteCall("String.new", 1)
			for _, c := range ce.tokenizer.CurrentToken().Literal {
				ce.vmwriter.WritePush(vmwriter.CONSTANT, int(c))
				ce.vmwriter.WriteCall("String.appendChar", 2)
			}
		}
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	} else if slices.Contains(keywordConstants, token.Keyword(ce.tokenizer.CurrentToken().Literal)) {
		ce.print(ce.tokenizer.CurrentToken().Xml())
		ce.tokenizer.Advance()
	} else if ce.tokenizer.CurrentToken().Literal == "(" {
		ce.process("(")
		ce.compileExpression()
		ce.process(")")
	} else if ce.tokenizer.CurrentToken().Literal == "-" || ce.tokenizer.CurrentToken().Literal == "~" {
		op := ce.tokenizer.CurrentToken().Literal
		ce.process(op)
		ce.compileTerm()
		if op == "-" {
			ce.vmwriter.WriteArithmetic(vmwriter.NEG)
		} else {
			ce.vmwriter.WriteArithmetic(vmwriter.NOT)
		}
	} else {
		name := ce.tokenizer.CurrentToken().Literal
		if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
			panic("expected identifier but got " + name)
		}
		ce.tokenizer.Advance()

		if ce.tokenizer.CurrentToken().Literal == "[" {
			ce.process("[")
			ce.compileExpression()
			ce.process("]")
		} else if ce.tokenizer.CurrentToken().Literal == "." {
			ce.process(".")
			subroutineName := ce.tokenizer.CurrentToken().Literal
			if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER {
				panic("expected identifier but got " + subroutineName)
			}
			ce.print("<identifier> " + "name: " + subroutineName + ", category: subroutine, index: -1, usage: using" + " </identifier>")
			ce.tokenizer.Advance()
			ce.process("(")
			n := ce.compileExpressionList()
			ce.process(")")
			ce.vmwriter.WriteCall(name+"."+subroutineName, n)
		} else if ce.tokenizer.CurrentToken().Literal == "(" {
			ce.process("(")
			ce.compileExpressionList()
			ce.process(")")
		} else {
			if ce.subroutineST.IndexOf(name) != -1 {
				ce.print("<identifier> " + "name: " + name + ", category: " + string(ce.subroutineST.KindOf(name)) + ", index: " + strconv.Itoa(ce.subroutineST.IndexOf(name)) + ", usage: using" + " </identifier>")
				ce.vmwriter.WritePush(kindSegmentMap[ce.subroutineST.KindOf(name)], ce.subroutineST.IndexOf(name))
			} else if ce.classST.IndexOf(name) != -1 {
				ce.print("<identifier> " + "name: " + name + ", category: " + string(ce.classST.KindOf(name)) + ", index: " + strconv.Itoa(ce.classST.IndexOf(name)) + ", usage: using" + " </identifier>")
				ce.vmwriter.WritePush(kindSegmentMap[ce.classST.KindOf(name)], ce.classST.IndexOf(name))
			} else {
				panic("undefined variable " + name)
			}
		}
	}

	ce.indent--
	ce.print("</term>")
}

func (ce *CompilationEngine) compileExpressionList() int {
	count := 0
	ce.print("<expressionList>")
	ce.indent++

	for ce.tokenizer.CurrentToken().Literal != ")" {
		ce.compileExpression()
		count++
		if ce.tokenizer.CurrentToken().Literal == "," {
			ce.process(",")
		}
	}

	ce.indent--
	ce.print("</expressionList>")
	return count
}

func (ce *CompilationEngine) process(str string) string {
	if ce.tokenizer.CurrentToken().Literal != str {
		panic("expected " + str + " but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
	return str
}

func (ce *CompilationEngine) processType() string {
	types := []token.Keyword{token.INT, token.CHAR, token.BOOLEAN}
	t := ce.tokenizer.CurrentToken().Literal
	if ce.tokenizer.CurrentToken().Type != token.IDENTIFIER && !slices.Contains(types, token.Keyword(t)) {
		panic("expected type but got " + ce.tokenizer.CurrentToken().Literal)
	}
	ce.print(ce.tokenizer.CurrentToken().Xml())
	ce.tokenizer.Advance()
	return t
}

func (ce *CompilationEngine) print(str string) {
	indentation := ""
	for i := 0; i < ce.indent; i++ {
		indentation += "  "
	}
	ce.XML += indentation + str + "\n"
}
