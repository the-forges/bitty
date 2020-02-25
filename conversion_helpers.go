package main

import "math"

// FindUnitSymbolPairBySymbol takes a UnitStandard and a symbol in order to
// find and return the UnitSymbolPair for that standard and symbol, or false
// if the UnitSymbolPair cannot be found.
func FindUnitSymbolPairBySymbol(std UnitStandard, sym UnitSymbol) (UnitSymbolPair, bool) {
	var stdMatch, symMatch bool
	for _, p := range unitSymbolPairs {
		stdMatch = p.Standard() == std
		symMatch = (p.Least() == sym || p.Greatest() == sym)
		if stdMatch && symMatch {
			return p, true
		}
	}
	return nil, false
}

// FindUnitSymbolPairByExponent takes a UnitStandard and an exponent in order to
// find and return the UnitSymbolPair for that standard and exponent, or false
// if the UnitSymbolPair cannot be found.
func FindUnitSymbolPairByExponent(std UnitStandard, exp uint) (UnitSymbolPair, bool) {
	var stdMatch, expMatch bool
	for _, p := range unitSymbolPairs {
		stdMatch = p.Standard() == std
		expMatch = p.Exponent() == exp
		if stdMatch && expMatch {
			return p, true
		}
	}
	return nil, false
}

// FindGreatestUnitSymbol finds the greatest of two unit symbols for a given
// exponent by standard
func FindGreatestUnitSymbol(std UnitStandard, exp uint) (UnitSymbol, bool) {
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	return pair.Greatest(), ok
}

// FindLeastUnitSymbol finds the least of two unit symbols for a given
// exponent by standard
func FindLeastUnitSymbol(std UnitStandard, exp uint) (UnitSymbol, bool) {
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	return pair.Least(), ok
}

func UnitSymbolByteSize(std UnitStandard, sym UnitSymbol, size float64) float64 {
	pair, ok := FindUnitSymbolPairBySymbol(std, sym)
	if !ok {
		return float64(0)
	}
	exp := float64(pair.Exponent() * 10)
	bytes := float64(math.Exp2(exp) * size)
	switch sym {
	case Bit:
		return float64(size * 8)
	case Byte:
		return float64(size)
	case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
		return float64(bytes * 0.125)
	case KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return float64(bytes)
	default:
		return float64(0)
	}
}
