package main

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

var iecExponentUnitMap = [][]UnitSymbol{
	[]UnitSymbol{Bit, Byte},
	[]UnitSymbol{Kib, KiB},
	[]UnitSymbol{Mib, MiB},
	[]UnitSymbol{Gib, GiB},
	[]UnitSymbol{Tib, TiB},
	[]UnitSymbol{Pib, PiB},
	[]UnitSymbol{Eib, EiB},
	[]UnitSymbol{Zib, ZiB},
	[]UnitSymbol{Yib, YiB},
}

// IECUnitSymbolPair represents a binary unit symbol pair as defined by the 9th
// edition SI standard
type IECUnitSymbolPair struct {
	least, greatest UnitSymbol
	exponent        int
}

func NewIECUnitSymbolPair(l, r UnitSymbol, e int) UnitSymbolPair {
	return &IECUnitSymbolPair{least: l, greatest: r, exponent: e}
}

func (pair *IECUnitSymbolPair) Standard() UnitStandard {
	return IEC
}

func (pair *IECUnitSymbolPair) Exponent() int {
	return pair.exponent
}

func (pair *IECUnitSymbolPair) Least() UnitSymbol {
	return pair.least
}

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

func (u *IECUnit) Standard() UnitStandard {
	return IEC
}

func (u *IECUnit) Exponent() int {
	return u.exponent
}

func (u *IECUnit) Symbol() UnitSymbol {
	return u.symbol
}

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

func findNearestIECUnitSymbols(exp int) []UnitSymbol {
	return iecExponentUnitMap[exp]
}

// Add attempts to add one Unit to another
func (u *IECUnit) Add(unit Unit) Unit {
	var (
		ru   *IECUnit
		exp  int
		nsym UnitSymbol
		size float64
	)
	// Validate both sides for valid symbols
	_, upok := FindUnitSymbolPairBySymbol(IEC, u.Symbol())
	_, rpok := FindUnitSymbolPairBySymbol(IEC, unit.Symbol())
	if !upok && !rpok {
		n, _ := NewIECUnit(u.Size(), Byte)
		return n
	}
	if upok && !rpok {
		return u
	}
	if !upok && rpok {
		return unit
	}
	ru = unit.(*IECUnit)
	// Lets get adding
	left := u.ByteSize()
	right := ru.ByteSize()
	total := left + right
	nexp := int(math.Round(math.Log2(total) / 10))
	if nexp > u.Exponent() && nexp > ru.Exponent() {
		exp = nexp
	} else {
		if u.Exponent() >= ru.Exponent() {
			exp = u.Exponent()
		} else {
			exp = ru.Exponent()
		}
	}
	nsym, _ = FindGreatestUnitSymbol(IEC, exp)
	lsym, _ := FindLeastUnitSymbol(IEC, exp)
	size = BytesToUnitSymbolSize(IEC, nsym, total)
	lsize := BytesToUnitSymbolSize(IEC, lsym, total)
	if size < 1 {
		nsym = lsym
		size = lsize
	}
	nu, _ := NewIECUnit(size, nsym)
	return nu
}

// Subtract attempts to subtract one Unit from another
func (u *IECUnit) Subtract(unit Unit) Unit {
	var (
		ru         *IECUnit
		ok         bool
		total      float64
		nexp       int
		lsym, gsym UnitSymbol
	)
	if ru, ok = unit.(*IECUnit); !ok {
		return u
	}
	left := u.ByteSize()
	right := ru.ByteSize()
	if left > right {
		total = left - right
	} else {
		total = right - left
	}
	if total > 0 {
		nexp = int(math.Round(math.Log2(total) / 10))
	} else {
		nexp = 0
	}
	if u.Symbol() != ru.Symbol() || nexp != u.Exponent() {
		lsym, ok = FindLeastUnitSymbol(IEC, nexp)
		if !ok {
			return u
		}
		gsym, ok = FindGreatestUnitSymbol(IEC, nexp)
		if !ok {
			return u
		}
	} else {
		lsym, gsym = u.Symbol(), u.Symbol()
	}
	smlSize := BytesToUnitSymbolSize(IEC, lsym, total)
	lrgSize := BytesToUnitSymbolSize(IEC, gsym, total)

	if lrgSize > 0 {
		nu, err := NewIECUnit(lrgSize, gsym)
		if err != nil {
			return u
		}
		return nu
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
