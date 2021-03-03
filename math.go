package bitty

import (
	"fmt"
)

// AddUnits takes two units with valid symbols, sums them, then returns a new unit
// AddUnits will always default to the left unit's symbol and exponent
func AddUnits(lu, ru Unit) (Unit, error) {
	var (
		newExponent int
		newSymbol   UnitSymbol
		newSize     float64
		lok, rok    bool
	)
	// validate that the units can be added to each other
	lok, rok = ValidateSymbols(lu.Symbol(), ru.Symbol())
	if !lok && !rok {
		return nil, fmt.Errorf("unable to add units with invalid symbols together: %s + %s", lu.Symbol(), ru.Symbol())
	}
	if lok && !rok {
		return lu, fmt.Errorf("unable to add unit with invalid symbol: %s", ru.Symbol())
	}
	if rok && !lok {
		return ru, fmt.Errorf("unable to add unit with invalid symbol: %s", lu.Symbol())
	}

	leftByte, rightByte := lu.ByteSize(), ru.ByteSize()
	totalByte := leftByte + rightByte

	newSymbol = lu.Symbol()
	newExponent = lu.Exponent()
	newSize = BytesToUnitSymbolSize(lu.Standard(), newSymbol, totalByte)

	var u Unit
	switch lu.Standard() {
	case IEC:
		u = &IECUnit{
			size:     newSize,
			symbol:   newSymbol,
			exponent: newExponent,
		}
	case SI:
		u = &SIUnit{
			size:     newSize,
			symbol:   newSymbol,
			exponent: newExponent,
		}
	}
	return u, nil
}

// SubtractUnits takes two units with valid symbols, subtracts them, then returns a new unit
// SubtractUnits will always default to the left unit's symbol and exponent
func SubtractUnits(lu, ru Unit) (Unit, error) {
	var (
		newExponent int
		newSymbol   UnitSymbol
		newSize     float64
		lok, rok    bool
	)
	// validate that the units can be added to each other
	lok, rok = ValidateSymbols(lu.Symbol(), ru.Symbol())
	if !lok && !rok {
		return nil, fmt.Errorf("unable to add units with invalid symbols together: %s + %s", lu.Symbol(), ru.Symbol())
	}
	if lok && !rok {
		return lu, fmt.Errorf("unable to subtract unit with invalid symbol: %s", ru.Symbol())
	}
	if rok && !lok {
		return ru, fmt.Errorf("unable to subtract unit with invalid symbol: %s", lu.Symbol())
	}

	leftByte, rightByte := lu.ByteSize(), ru.ByteSize()
	totalByte := leftByte - rightByte

	newSymbol = lu.Symbol()
	newExponent = lu.Exponent()
	newSize = BytesToUnitSymbolSize(lu.Standard(), newSymbol, totalByte)

	var u Unit
	switch lu.Standard() {
	case IEC:
		u = &IECUnit{
			size:     newSize,
			symbol:   newSymbol,
			exponent: newExponent,
		}
	case SI:
		u = &SIUnit{
			size:     newSize,
			symbol:   newSymbol,
			exponent: newExponent,
		}
	}
	return u, nil
}
