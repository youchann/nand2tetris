package symboltable

import "github.com/youchann/nand2tetris/11-2_vmwriter/token"

type row struct {
	name  string
	Type  string // type is a reserved word in Go
	kind  token.VariableKind
	index int
}

type SymbolTable struct {
	table map[string]row
}

func New() *SymbolTable {
	return &SymbolTable{
		table: map[string]row{},
	}
}

func (st *SymbolTable) Reset() {
	st.table = map[string]row{}
}

func (st *SymbolTable) Define(name, Type string, kind token.VariableKind) {
	index := st.VarCount(kind)
	st.table[name] = row{name, Type, kind, index}
}

func (st *SymbolTable) VarCount(kind token.VariableKind) int {
	count := 0
	for _, r := range st.table {
		if r.kind == kind {
			count++
		}
	}
	return count
}

func (st *SymbolTable) KindOf(name string) token.VariableKind {
	r, ok := st.table[name]
	if !ok {
		return token.NONE
	}
	return r.kind
}

func (st *SymbolTable) TypeOf(name string) string {
	r, ok := st.table[name]
	if !ok {
		return ""
	}
	return r.Type
}

func (st *SymbolTable) IndexOf(name string) int {
	r, ok := st.table[name]
	if !ok {
		return -1
	}
	return r.index
}
