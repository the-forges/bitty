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

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// NewUnit takes a UnitStandard, float64, and UnitSymbol, returning a valid Unit
func NewUnit(std UnitStandard, size float64, sym UnitSymbol) (Unit, error) {
	switch std {
	case IEC:
		return NewIECUnit(size, sym)
	case SI:
		return NewSIUnit(size, sym)
	default:
		return nil, fmt.Errorf("%s is currently not a supported standard", string(std))
	}
}

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

// FindStandardBySymbol takes a unit symbol, searches for a symbol pair that
// matches, and returns the standard for that pair
func FindStandardBySymbol(sym UnitSymbol) (UnitStandard, bool) {
	for _, p := range unitSymbolPairs {
		if p.Least() == sym || p.Greatest() == sym {
			return p.Standard(), true
		}
	}
	return UnitStandard(0), false
}

// FindExponentBySymbol takes a symbol and returns the exponent
func FindExponentBySymbol(sym UnitSymbol) (int, bool) {
	s, ok := FindStandardBySymbol(sym)
	if !ok {
		return 0, false
	}
	pair, ok := FindUnitSymbolPairBySymbol(s, sym)
	if !ok {
		return 0, false
	}
	return pair.Exponent(), true
}

// FindGreatestUnitSymbol finds the greatest of two unit symbols for a given
// exponent by standard
func FindGreatestUnitSymbol(std UnitStandard, exp int) (UnitSymbol, bool) {
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	if !ok {
		return Byte, false
	}
	return pair.Greatest(), true
}

// FindLeastUnitSymbol finds the least of two unit symbols for a given
// exponent by standard
func FindLeastUnitSymbol(std UnitStandard, exp int) (UnitSymbol, bool) {
	pair, ok := FindUnitSymbolPairByExponent(std, exp)
	if !ok {
		return Byte, false
	}
	return pair.Least(), true
}

// UnitSymbolToByteSize converts the size from one unit into bytes
func UnitSymbolToByteSize(std UnitStandard, sym UnitSymbol, size float64) float64 {
	var exp, bytes, errVal float64
	pair, ok := FindUnitSymbolPairBySymbol(std, sym)
	if !ok {
		return errVal
	}
	switch std {
	case IEC:
		exp = float64(pair.Exponent() * 10)
		bytes = float64(math.Exp2(exp) * size)
	case SI:
		exp = float64(pair.Exponent())
		bytes = float64(math.Pow10(int(exp)) * size)
	default:
		return errVal
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
	case db, hb, kb, Mb, Gb, Tb, Pb, Eb, Zb, Yb:
		return float64(bytes * 0.125)
	case dB, hB, kB, MB, GB, TB, PB, EB, ZB, YB:
		return float64(bytes)
	default:
		return errVal
	}
}

// BytesToUnitSymbolSize converts bytes to the best unit size as a float64
func BytesToUnitSymbolSize(std UnitStandard, sym UnitSymbol, size float64) float64 {
	var exp float64
	pair, ok := FindUnitSymbolPairBySymbol(std, sym)
	if !ok {
		return float64(0)
	}
	switch std {
	case IEC:
		exp = math.Pow(2, float64(pair.Exponent()*10))
	case SI:
		exp = math.Pow10(int(pair.Exponent()))
	}
	switch sym {
	case Bit:
		return float64(size * 8)
	case Byte:
		return size / exp
	case db, hb, kb, Mb, Gb, Tb, Pb, Eb, Zb, Yb,
		dB, hB, kB, MB, GB, TB, PB, EB, ZB, YB:
		return size / exp
	case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
		return (size / (size * 0.125)) / 2
	case KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return size / exp
	}
	return float64(0)
}

// Parse parses a string representation of a unit size in the format of
// "<size><unit symbol>" or "<size> <unit symbol>" in order to instantiate and
// return a Unit with the correct standard, exponent, size, and symbol
func Parse(s string) (Unit, error) {
	stderr := fmt.Errorf("%s could not be found in the standard", s)
	parseerr := fmt.Errorf("%s could not be parsed", s)
	r := regexp.MustCompile(`^(\d+)\s{0,}(\w+)$`)
	if r.MatchString(s) {
		m := r.FindStringSubmatch(s)
		size, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, parseerr
		}
		symbol := UnitSymbol(m[2])
		standard, ok := FindStandardBySymbol(symbol)
		if !ok {
			return nil, stderr
		}
		return NewUnit(standard, float64(size), symbol)
	}
	return nil, parseerr
}

// ConvertUnitStd takes a unit from one standard and converts it to another
func ConvertUnitStd(u Unit, std UnitStandard) (Unit, error) {
	var (
		sym       UnitSymbol
		bytes     = u.ByteSize()
		size, exp float64
		ok        bool
	)
	switch std {
	case IEC:
		exp = math.Round(math.Log2(bytes) / 10)
		size = bytes / math.Pow(2, exp)
	case SI:
		exp = math.Floor(math.Log10(bytes))
		size = bytes / math.Pow10(int(exp))
	default:
		return nil, NewErrUnitStandardNotSupported(std)
	}
	if size < 1 {
		if sym, ok = FindLeastUnitSymbol(std, int(exp)); !ok {
			return nil, fmt.Errorf(ErrUnitExponentNotSupportedf, exp)
		}
	} else {
		if sym, ok = FindGreatestUnitSymbol(std, int(exp)); !ok {
			return nil, fmt.Errorf(ErrUnitExponentNotSupportedf, exp)
		}
	}
	return NewUnit(std, size, sym)
}
