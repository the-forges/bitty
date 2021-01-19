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
	"math"
)

var iecUnitExponentMap = map[UnitSymbol]int{
	Bit:  0,
	Byte: 0,
	Kib:  1,
	KiB:  1,
	Mib:  2,
	MiB:  2,
	Gib:  3,
	GiB:  3,
	Tib:  4,
	TiB:  4,
	Pib:  5,
	PiB:  5,
	Eib:  6,
	EiB:  6,
	Zib:  7,
	ZiB:  7,
	Yib:  8,
	YiB:  8,
}

// IECUnitSymbolPair represents a base 2 binary unit symbol pair as defined by the 9th
// edition SI standard
type IECUnitSymbolPair struct {
	least, greatest UnitSymbol
	exponent        int
}

// NewIECUnitSymbolPair takes a UnitStandard and returns a new UnitSymbolPair
func NewIECUnitSymbolPair(l, r UnitSymbol, e int) UnitSymbolPair {
	return &IECUnitSymbolPair{least: l, greatest: r, exponent: e}
}

// Standard returns the UnitStandard of a IECUnitSymbolPair: IEC
func (pair *IECUnitSymbolPair) Standard() UnitStandard {
	return IEC
}

// Exponent returns the exponent of a IECUnitSymbolPair
func (pair *IECUnitSymbolPair) Exponent() int {
	return pair.exponent
}

// Least returns the least UnitSymbol of a IECUnitSymbolPair
func (pair *IECUnitSymbolPair) Least() UnitSymbol {
	return pair.least
}

// Greatest returns the greatest UnitSymbol of a IECUnitSymbolPair
func (pair *IECUnitSymbolPair) Greatest() UnitSymbol {
	return pair.greatest
}

// IECUnit handles binary units as dictated by SI Standards
type IECUnit struct {
	// size is the size as measured by the symbol (UnitSymbol), which is equivalent to:
	// 		b(2^10)^n(1/8)
	// 		B(2^10)^n
	size     float64
	symbol   UnitSymbol
	exponent int
}

// NewIECUnit returns a *IECUnit with the proper exponent included
func NewIECUnit(size float64, sym UnitSymbol) (*IECUnit, error) {
	if pair, ok := FindUnitSymbolPairBySymbol(IEC, sym); ok {
		return &IECUnit{size, sym, pair.Exponent()}, nil
	}
	return nil, NewErrUnitSymbolNotSupported(sym)
}

// Standard returns the UnitStandard of a IECUnit: IEC
func (u *IECUnit) Standard() UnitStandard {
	return IEC
}

// Exponent returns the exponent of a IECUnit
func (u *IECUnit) Exponent() int {
	return u.exponent
}

// Symbol returns the UnitSymbol of a IECUnit
func (u *IECUnit) Symbol() UnitSymbol {
	return u.symbol
}

// Size returns the size of an IECUnit
func (u *IECUnit) Size() float64 {
	return u.size
}

// BitSize returns the size of the Unit measured in bits
func (u *IECUnit) BitSize() float64 {
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
func (u *IECUnit) ByteSize() float64 {
	return UnitSymbolToByteSize(IEC, u.Symbol(), u.Size())
}

// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
func (u *IECUnit) SizeInUnit(symbol UnitSymbol) float64 {
	_, uok := FindUnitSymbolPairBySymbol(IEC, u.Symbol())
	p, ok := FindUnitSymbolPairBySymbol(IEC, symbol)
	if uok && ok {
		var (
			left    = u.ByteSize()
			right   = UnitSymbolToByteSize(IEC, symbol, u.Size())
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
func (u *IECUnit) Add(unit Unit) Unit {
	var (
		nexp     int
		nsym     UnitSymbol
		size     float64
		lok, rok bool
	)
	// Validate both sides for valid symbols
	lok, rok = ValidateSymbols(u.Symbol(), unit.Symbol())
	if !lok && !rok {
		nu, _ := NewIECUnit(0, Byte)
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
	lsym, ok := FindLeastUnitSymbol(IEC, nexp)
	gsym, ok := FindGreatestUnitSymbol(IEC, nexp)
	if !ok {
		nu, _ := NewIECUnit(0, Byte)
		return nu
	}
	smallSize := BytesToUnitSymbolSize(IEC, lsym, total)
	lrgSize := BytesToUnitSymbolSize(IEC, gsym, total)
	if lrgSize < 1 {
		nsym = lsym
		size = smallSize
	} else {
		nsym = gsym
		size = lrgSize
	}
	nu, _ := NewIECUnit(size, nsym)
	return nu
}

// Subtract attempts to subtract one Unit from another
func (u *IECUnit) Subtract(unit Unit) Unit {
	var (
		ok, lok, rok, neg bool
		total             float64
		nexp              int
		lsym, gsym        UnitSymbol
	)
	lok, rok = ValidateSymbols(u.Symbol(), unit.Symbol())
	if !lok && !rok {
		nu, _ := NewIECUnit(0, Byte)
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
		nexp = int(math.Round(math.Log2(total) / 10))
	}
	lsym, ok = FindLeastUnitSymbol(IEC, nexp)
	gsym, ok = FindGreatestUnitSymbol(IEC, nexp)
	if !ok {
		nu, _ := NewIECUnit(0, Byte)
		return nu
	}
	smlSize := BytesToUnitSymbolSize(IEC, lsym, total)
	lrgSize := BytesToUnitSymbolSize(IEC, gsym, total)
	if lrgSize >= 0 {
		if neg {
			lrgSize = -lrgSize
		}
		nu, err := NewIECUnit(lrgSize, gsym)
		if err != nil {
			return u
		}
		return nu
	}
	if neg {
		smlSize = -smlSize
	}
	nu, err := NewIECUnit(smlSize, lsym)
	if err != nil {
		return u
	}
	return nu
}

// Multiply attempts to multiply one Unit by another
func (u *IECUnit) Multiply(unit Unit) Unit {
	return nil
}

// Divide attempts to divide one Unit by another
func (u *IECUnit) Divide(unit Unit) Unit {
	return nil
}
