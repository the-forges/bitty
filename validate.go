package bitty

import "fmt"

// ValidateSymbol checks that a symbol is valid
func ValidateSymbol(sym UnitSymbol) bool {
	str := fmt.Sprintf("%d %s", 0, sym)
	if _, err := Parse(str); err != nil {
		return false
	}
	return true
}

// ValidateSymbols validates all symbols, returning a tuple of booleans
func ValidateSymbols(l, r UnitSymbol) (bool, bool) {
	return ValidateSymbol(l), ValidateSymbol(r)
}
