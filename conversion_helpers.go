package bitty

/*
	Copyright 2020 IBM

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
func FindUnitSymbolPairByExponent(std UnitStandard, exp int) (UnitSymbolPair, bool) {
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
func FindGreatestUnitSymbol(std UnitStandard, exp int) (UnitSymbol, bool) {
	var sym UnitSymbol
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	if !ok {
		return sym, false
	}
	return pair.Greatest(), true
}

// FindLeastUnitSymbol finds the least of two unit symbols for a given
// exponent by standard
func FindLeastUnitSymbol(std UnitStandard, exp int) (UnitSymbol, bool) {
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	return pair.Least(), ok
}

// UnitSymbolToByteSize converts the size from one unit into bytes
func UnitSymbolToByteSize(std UnitStandard, sym UnitSymbol, size float64) float64 {
	var exp, bytes float64
	pair, ok := FindUnitSymbolPairBySymbol(std, sym)
	if !ok {
		return float64(0)
	}
	if std == IEC {
		exp = float64(pair.Exponent() * 10)
		bytes = float64(math.Exp2(exp) * size)
	} else {
		exp = float64(pair.Exponent())
		bytes = float64(math.Pow10(int(exp)) * size)
	}
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

// BytesToUnitSymbolSize ...
func BytesToUnitSymbolSize(std UnitStandard, sym UnitSymbol, size float64) float64 {
	var exp float64
	pair, ok := FindUnitSymbolPairBySymbol(std, sym)
	if !ok {
		return float64(0)
	}
	if std == IEC {
		exp = float64(pair.Exponent() * 10)
	} else {
		exp = float64(pair.Exponent())
	}
	switch sym {
	case Bit:
		return float64(size * 8)
	case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
		return (size / (size * 0.125)) / 2
	case Byte, KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return size / math.Pow(2, float64(exp))
	default:
		return float64(0)
	}
}
