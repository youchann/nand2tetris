package symboltable

import "github.com/youchann/nand2tetris/11-1_symboltable/token"

type row struct {
	name  string
	Type  string // type is a reserved word in Go
	kind  token.VariableKind
	index int
}

type symbolTable struct {
	table map[string]row
}

func New() *symbolTable {
	return &symbolTable{
		table: map[string]row{},
	}
}

func (st *symbolTable) Reset() {
	st.table = map[string]row{}
}

func (st *symbolTable) Define(name, Type string, kind token.VariableKind) {
	index := st.VarCount(kind)
	st.table[name] = row{name, Type, kind, index}
}

func (st *symbolTable) VarCount(kind token.VariableKind) int {
	count := 0
	for _, r := range st.table {
		if r.kind == kind {
			count++
		}
	}
	return count
}

func (st *symbolTable) KindOf(name string) token.VariableKind {
	r, ok := st.table[name]
	if !ok {
		return token.NONE
	}
	return r.kind
}

func (st *symbolTable) TypeOf(name string) string {
	r, ok := st.table[name]
	if !ok {
		return ""
	}
	return r.Type
}

func (st *symbolTable) IndexOf(name string) int {
	r, ok := st.table[name]
	if !ok {
		return -1
	}
	return r.index
}
