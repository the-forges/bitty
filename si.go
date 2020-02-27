package main

// IECUnitSymbolPair represents a binary unit symbol pair as defined by the 9th
// edition SI standard
type SIUnitSymbolPair struct {
	least, greatest UnitSymbol
	exponent        int
}

func NewSIUnitSymbolPair(l, r UnitSymbol, e int) UnitSymbolPair {
	return &SIUnitSymbolPair{least: l, greatest: r, exponent: e}
}

func (pair *SIUnitSymbolPair) Standard() UnitStandard {
	return SI
}

func (pair *SIUnitSymbolPair) Exponent() int {
	return pair.exponent
}

func (pair *SIUnitSymbolPair) Least() UnitSymbol {
	return pair.least
}

func (pair *SIUnitSymbolPair) Greatest() UnitSymbol {
	return pair.greatest
}
