package token

type CommandType string

const (
	C_ARITHMETIC CommandType = "C_ARITHMETIC"
	C_PUSH       CommandType = "C_PUSH"
	C_POP        CommandType = "C_POP"
	C_LABEL      CommandType = "C_LABEL"
	C_GOTO       CommandType = "C_GOTO"
	C_IF         CommandType = "C_IF"
	C_FUNCTION   CommandType = "C_FUNCTION"
	C_RETURN     CommandType = "C_RETURN"
	C_CALL       CommandType = "C_CALL"
)

type Segment string

const (
	SEGMENT_LOCAL    Segment = "local"
	SEGMENT_ARGUMENT Segment = "argument"
	SEGMENT_THIS     Segment = "this"
	SEGMENT_THAT     Segment = "that"
	SEGMENT_POINTER  Segment = "pointer"
	SEGMENT_TEMP     Segment = "temp"
	SEGMENT_CONSTANT Segment = "constant"
	SEGMENT_STATIC   Segment = "static"
)

type CommandSymbol string

const (
	PUSH     CommandSymbol = "push"
	POP      CommandSymbol = "pop"
	LABEL    CommandSymbol = "label"
	GOTO     CommandSymbol = "goto"
	IF_GOTO  CommandSymbol = "if-goto"
	CALL     CommandSymbol = "call"
	FUNCTION CommandSymbol = "function"
	RETURN   CommandSymbol = "return"
	ADD      CommandSymbol = "add"
	SUB      CommandSymbol = "sub"
	NEG      CommandSymbol = "neg"
	EQ       CommandSymbol = "eq"
	GT       CommandSymbol = "gt"
	LT       CommandSymbol = "lt"
	AND      CommandSymbol = "and"
	OR       CommandSymbol = "or"
	NOT      CommandSymbol = "not"
)
