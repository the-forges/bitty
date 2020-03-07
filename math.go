package bitty

import (
	"fmt"
	"math"
)

// AddUnits takes to units with valid symbols, sums them, then returns a new unit
func AddUnits(lu, ru Unit) (Unit, error) {
	var (
		nexp     int
		nsym     UnitSymbol
		size     float64
		lok, rok bool
	)
	// validate that the units can be added to each other
	lok, rok = ValidateSymbols(lu.Symbol(), ru.Symbol())
	if !lok && !rok {
		return nil, fmt.Errorf("unable to add units with invalid symbols together: %f %s + %f %s", lu.Symbol(), ru.Symbol())
	}
	if lok && !rok {
		return lu, nil
	}
	if rok && !lok {
		return ru, nil
	}
	// take the generic byte sizes for easy addition
	lbyte, rbyte := lu.ByteSize(), ru.ByteSize()
	tbyte := lbyte + rbyte
	if tbyte > 0 {
		nexp = int(math.Round(math.Log2(tbyte) / 10))
	}
	// determine the dominate symbol (prefer the left) and normalize units
	if lu.Symbol() != ru.Symbol() {
		//ru
	}
	if lu.Exponent() >= ru.Exponent() {
		nexp = lu.Exponent()
	} else {
		nexp = ru.Exponent()
	}
}
