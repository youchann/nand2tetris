package symboltable

import "fmt"

type Table struct {
	symbols map[string]int
}

func getInitialSymbolTable() map[string]int {
	initialSymbolTable := map[string]int{
		"SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4,
		"SCREEN": 16384, "KBD": 24576,
	}
	// initialize Register Address
	for i := 0; i < 16; i++ {
		initialSymbolTable[fmt.Sprintf("R%d", i)] = i
	}
	return initialSymbolTable
}

func New() *Table {
	return &Table{
		symbols: getInitialSymbolTable(),
	}
}

func (t *Table) AddEntry(symbol string, address int) {
	t.symbols[symbol] = address
}

func (t *Table) Contains(symbol string) bool {
	_, ok := t.symbols[symbol]
	return ok
}

func (t *Table) GetAddress(symbol string) int {
	return t.symbols[symbol]
}
