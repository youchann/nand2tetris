package vmwriter

import "strconv"

type segment string

const (
	CONSTANT segment = "constant"
	ARGUMENT segment = "argument"
	LOCAL    segment = "local"
	STATIC   segment = "static"
	THIS     segment = "this"
	THAT     segment = "that"
	POINTER  segment = "pointer"
	TEMP     segment = "temp"
)

type command string

const (
	ADD command = "add"
	SUB command = "sub"
	NEG command = "neg"
	EQ  command = "eq"
	GT  command = "gt"
	LT  command = "lt"
	AND command = "and"
	OR  command = "or"
	NOT command = "not"
)

type VMWriter struct {
	className string
	Code      string
	hasIndent bool
}

func New(n string) *VMWriter {
	return &VMWriter{
		className: n,
		Code:      "",
		hasIndent: false,
	}
}

func (w *VMWriter) WritePush(segment segment, index int) {
	w.write("push " + string(segment) + " " + strconv.Itoa(index))
}

func (w *VMWriter) WritePop(segment segment, index int) {
	w.write("pop " + string(segment) + " " + strconv.Itoa(index))
}

func (w *VMWriter) WriteArithmetic(command command) {
	w.write(string(command))
}

func (w *VMWriter) WriteLabel(label string) {
	w.hasIndent = false
	w.write("label " + label)
	w.hasIndent = true
}

func (w *VMWriter) WriteGoto(label string) {
	w.write("goto " + label)
}

func (w *VMWriter) WriteIf(label string) {
	w.write("if-goto " + label)
}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.write("call " + name + " " + strconv.Itoa(nArgs))
}

func (w *VMWriter) WriteFunction(name string, nLocals int) {
	w.write("function " + w.className + "." + name + " " + strconv.Itoa(nLocals))
	w.hasIndent = true
}

func (w *VMWriter) WriteReturn() {
	w.write("return")
}

func (w *VMWriter) write(str string) {
	if w.hasIndent {
		w.Code += "    "
	}
	w.Code += str + "\n"
}
