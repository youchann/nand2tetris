package symboltable

type Kind string

const (
	STATIC    Kind = "STATIC"
	FIELD     Kind = "FIELD"
	ARGUMENT  Kind = "ARGUMENT"
	VAR_LOCAL Kind = "VAR"
	NONE      Kind = "NONE"
)

var KindMap = map[string]Kind{
	"static":   STATIC,
	"field":    FIELD,
	"argument": ARGUMENT,
	"var":      VAR_LOCAL,
}

type row struct {
	name  string
	Type  string // type is a reserved word in Go
	kind  Kind
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

func (st *SymbolTable) Define(name, Type string, kind Kind) {
	index := st.VarCount(kind)
	st.table[name] = row{name, Type, kind, index}
}

func (st *SymbolTable) VarCount(kind Kind) int {
	count := 0
	for _, r := range st.table {
		if r.kind == kind {
			count++
		}
	}
	return count
}

func (st *SymbolTable) KindOf(name string) Kind {
	r, ok := st.table[name]
	if !ok {
		return NONE
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
