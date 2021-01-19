package bitty

import "math"

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

// SIUnitSymbolPair represents a base 10 decimal unit symbol pair as defined by the 9th
// edition SI standard
type SIUnitSymbolPair struct {
	least, greatest UnitSymbol
	exponent        int
}

// NewSIUnitSymbolPair takes a UnitStandard and returns a new UnitSymbolPair
func NewSIUnitSymbolPair(l, r UnitSymbol, e int) UnitSymbolPair {
	return &SIUnitSymbolPair{least: l, greatest: r, exponent: e}
}

// Standard returns the UnitStandard of a SIUnitSymbolPair: SI
func (pair *SIUnitSymbolPair) Standard() UnitStandard {
	return SI
}

// Exponent returns the exponent of a SIUnitSymbolPair
func (pair *SIUnitSymbolPair) Exponent() int {
	return pair.exponent
}

// Least returns the least UnitSymbol of a SIUnitSymbolPair
func (pair *SIUnitSymbolPair) Least() UnitSymbol {
	return pair.least
}

// Greatest returns the greatest UnitSymbol of a SIUnitSymbolPair
func (pair *SIUnitSymbolPair) Greatest() UnitSymbol {
	return pair.greatest
}

// SIUnit handles binary units as dictated by SI Standards
type SIUnit struct {
	// size is the size as measured by the symbol (UnitSymbol), which is equivalent to:
	// 		b(2^10)^n(1/8)
	// 		B(2^10)^n
	size     float64
	symbol   UnitSymbol
	exponent int
}

// NewSIUnit returns a *SIUnit with the proper exponent included
func NewSIUnit(size float64, sym UnitSymbol) (*SIUnit, error) {
	if pair, ok := FindUnitSymbolPairBySymbol(SI, sym); ok {
		return &SIUnit{size, sym, pair.Exponent()}, nil
	}
	return nil, NewErrUnitSymbolNotSupported(sym)
}

// Standard returns the UnitStandard of a SIUnit: SI
func (u *SIUnit) Standard() UnitStandard {
	return SI
}

// Exponent returns the exponent of a SIUnit
func (u *SIUnit) Exponent() int {
	return u.exponent
}

// Symbol returns the UnitSymbol of a SIUnit
func (u *SIUnit) Symbol() UnitSymbol {
	return u.symbol
}

// Size returns the size of an SIUnit
func (u *SIUnit) Size() float64 {
	return u.size
}

// BitSize returns the size of the Unit measured in bits
func (u *SIUnit) BitSize() float64 {
	switch u.Symbol() {
	case Bit:
		return float64(u.ByteSize() / 8)
	case Byte,
		Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib,
		KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return float64(u.ByteSize() * 8)
	default:
		return float64(0)
	}
}

// ByteSize returns the size of the Unit measured in bytes
func (u *SIUnit) ByteSize() float64 {
	return UnitSymbolToByteSize(SI, u.Symbol(), u.Size())
}

// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
func (u *SIUnit) SizeInUnit(symbol UnitSymbol) float64 {
	_, uok := FindUnitSymbolPairBySymbol(SI, u.Symbol())
	p, ok := FindUnitSymbolPairBySymbol(SI, symbol)
	if uok && ok {
		var (
			left    = u.ByteSize()
			right   = UnitSymbolToByteSize(SI, symbol, u.Size())
			diffExp = float64(u.Exponent()) - float64(p.Exponent())
		)
		if diffExp > 0 {
			return right * diffExp
		}
		if left > 0.0 && right > 0.0 {
			return float64((left / right) * u.Size())
		}
	}
	return float64(0)
}

// Add attempts to add one Unit to another
func (u *SIUnit) Add(unit Unit) Unit {
	var (
		nexp     int
		nsym     UnitSymbol
		size     float64
		lok, rok bool
	)
	// Validate both sides for valid symbols
	lok, rok = ValidateSymbols(u.Symbol(), unit.Symbol())
	if !lok && !rok {
		nu, _ := NewSIUnit(0, Byte)
		return nu
	}
	if lok && !rok {
		return u
	}
	if rok && !lok {
		return unit
	}
	// Lets get adding
	left := u.ByteSize()
	right := unit.ByteSize()
	total := left + right
	if total > 0 {
		nexp = int(math.Round(math.Log2(total) / 10))
	}
	if u.Exponent() >= unit.Exponent() {
		nexp = u.Exponent()
	} else {
		nexp = unit.Exponent()
	}
	lsym, ok := FindLeastUnitSymbol(SI, nexp)
	gsym, ok := FindGreatestUnitSymbol(SI, nexp)
	if !ok {
		nu, _ := NewSIUnit(0, Byte)
		return nu
	}
	smallSize := BytesToUnitSymbolSize(SI, lsym, total)
	lrgSize := BytesToUnitSymbolSize(SI, gsym, total)
	if lrgSize < 1 {
		nsym = lsym
		size = smallSize
	} else {
		nsym = gsym
		size = lrgSize
	}
	nu, _ := NewSIUnit(size, nsym)
	return nu
}

// Subtract attempts to subtract one Unit from another
func (u *SIUnit) Subtract(unit Unit) Unit {
	var (
		ok, lok, rok, neg bool
		total             float64
		nexp              int
		lsym, gsym        UnitSymbol
	)
	lok, rok = ValidateSymbols(u.Symbol(), unit.Symbol())
	if !lok && !rok {
		nu, _ := NewSIUnit(0, Byte)
		return nu
	}
	if lok && !rok {
		return u
	}
	if rok && !lok {
		return unit
	}
	left := u.ByteSize()
	right := unit.ByteSize()
	if left >= right {
		total = left - right
	} else {
		total = right - left
		neg = true
	}
	if total > 0 {
		l := math.Log10(total)
		nexp = int(math.Floor(l))
	}
	lsym, ok = FindLeastUnitSymbol(SI, nexp)
	gsym, ok = FindGreatestUnitSymbol(SI, nexp)
	if !ok {
		nu, _ := NewSIUnit(0, Byte)
		return nu
	}
	smlSize := BytesToUnitSymbolSize(SI, lsym, total)
	lrgSize := BytesToUnitSymbolSize(SI, gsym, total)
	if lrgSize >= 0 {
		if neg {
			lrgSize = -lrgSize
		}
		nu, err := NewSIUnit(lrgSize, gsym)
		if err != nil {
			return u
		}
		return nu
	}
	if neg {
		smlSize = -smlSize
	}
	nu, err := NewSIUnit(smlSize, lsym)
	if err != nil {
		return u
	}
	return nu
}

// Multiply attempts to multiply one Unit by another
func (u *SIUnit) Multiply(unit Unit) Unit {
	return nil
}

// Divide attempts to divide one Unit by another
func (u *SIUnit) Divide(unit Unit) Unit {
	return nil
}
