package vmwriter

import "strconv"

type Segment string

const (
	CONSTANT Segment = "constant"
	ARGUMENT Segment = "argument"
	LOCAL    Segment = "local"
	STATIC   Segment = "static"
	THIS     Segment = "this"
	THAT     Segment = "that"
	POINTER  Segment = "pointer"
	TEMP     Segment = "temp"
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

func (w *VMWriter) WritePush(segment Segment, index int) {
	w.write("push " + string(segment) + " " + strconv.Itoa(index))
}

func (w *VMWriter) WritePop(segment Segment, index int) {
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
	w.hasIndent = false
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
