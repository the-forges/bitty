package bitty

import (
	"fmt"
	"math"
)

// TODO: Refactor Add into a cross unit generic AddUnits
// AddUnits takes to units with valid symbols, sums them, then returns a new unit
func AddUnits(lu, ru Unit) (Unit, error) {
	var (
		nexp                 int
		nsym                 UnitSymbol
		size                 float64
		lok, rok             bool
		leastsym, greatsym   UnitSymbol
		leastsize, greatsize float64
	)
	// validate that the units can be added to each other
	lok, rok = ValidateSymbols(lu.Symbol(), ru.Symbol())
	if !lok && !rok {
		return nil, fmt.Errorf("unable to add units with invalid symbols together: %f %s + %f %s", lu.Size(), lu.Symbol(), ru.Size(), ru.Symbol())
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
	// get unit symbols
	leastsym, lok = FindLeastUnitSymbol(lu.Standard(), nexp)
	greatsym, lok = FindGreatestUnitSymbol(lu.Standard(), nexp)
	// get unit symbol sizes
	if lok {
		leastsize = BytesToUnitSymbolSize(lu.Standard(), leastsym, tbyte)
		greatsize = BytesToUnitSymbolSize(lu.Standard(), greatsym, tbyte)
		// prefer unit symbol sizes greater than 1
		if greatsize < 1 {
			nsym = leastsym
			size = leastsize
		} else {
			nsym = greatsym
			size = greatsize
		}
		if size >= 1000 {
			nsym, size := 
		}
		switch lu.Standard() {
		case IEC:
			return NewIECUnit(size, nsym)
		case SI:
			return NewSIUnit(size, nsym)
		}
	}
	return nil, NewErrUnitStandardNotSupported(lu.Standard())
}
