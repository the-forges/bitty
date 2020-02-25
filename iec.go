package main

import (
	"errors"
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
	exponent        uint
}

func NewIECUnitSymbolPair(l, r UnitSymbol, e uint) UnitSymbolPair {
	return &IECUnitSymbolPair{least: l, greatest: r, exponent: e}
}

func (pair *IECUnitSymbolPair) Standard() UnitStandard {
	return IEC
}

func (pair *IECUnitSymbolPair) Exponent() uint {
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
	exponent uint
}

// NewIECUnit returns a *IECUnit with the proper exponent included
func NewIECUnit(size float64, sym UnitSymbol) (*IECUnit, error) {
	if pair, ok := FindUnitSymbolPairBySymbol(IEC, sym); ok {
		return &IECUnit{size, sym, pair.Exponent()}, nil
	}
	return nil, NewErrUnitSymbolNotSupported(sym)
}

// BitSize returns the size of the Unit measured in bits
func (u *IECUnit) BitSize() float64 {
	switch u.symbol {
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

func sizeInIECUnit(symbol UnitSymbol, bytes float64) float64 {
	uexp, ok := iecUnitExponentMap[symbol]
	if !ok {
		return float64(0)
	}
	exp := float64(uexp * 10)
	switch symbol {
	case Bit:
		return float64(bytes * 8)
	case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
		return (bytes / (bytes * 0.125)) / 2
	case Byte, KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return bytes / math.Pow(2, float64(exp))
	default:
		return float64(0)
	}
}

// ByteSize returns the size of the Unit measured in bytes
func (u *IECUnit) ByteSize() float64 {
	return UnitSymbolByteSize(IEC, u.symbol, u.size)
}

// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
func (u *IECUnit) SizeInUnit(symbol UnitSymbol) float64 {
	_, uok := FindUnitSymbolPairBySymbol(IEC, u.symbol)
	p, ok := FindUnitSymbolPairBySymbol(IEC, symbol)
	if uok && ok {
		var (
			left    = u.ByteSize()
			right   = UnitSymbolByteSize(IEC, symbol, u.size)
			diffExp = float64(u.exponent) - float64(p.Exponent())
		)
		if diffExp > 0 {
			return right * diffExp
		}
		if left > 0.0 && right > 0.0 {
			return float64((left / right) * u.size)
		}
	}
	return float64(0)
}

func findNearestIECUnitSymbols(exp uint) []UnitSymbol {
	return iecExponentUnitMap[exp]
}

func findGCDIECUnitSymbol(left, right UnitSymbol, optexp *uint) (UnitSymbol, error) {
	var (
		lexp, rexp *int
		exp        uint
		err        = errors.New("no available IEC unit symbol matched")
	)
	// Get initial exponent for unit symbols
	if e, ok := iecUnitExponentMap[left]; ok {
		lexp = &e
	}
	if e, ok := iecUnitExponentMap[right]; ok {
		rexp = &e
	}
	if lexp == nil && rexp == nil {
		return Bit, err
	}
	// Set an initial value for the exponent, validating for nil results on
	// either or both sides of an exponent search
	if lexp != nil && rexp == nil {
		exp = uint(*lexp)
	} else if lexp == nil && rexp != nil {
		exp = uint(*rexp)
	} else {
		if *lexp > *rexp {
			exp = uint(*rexp)
		} else {
			exp = uint(*lexp)
		}
	}
	// Update the exponent if the optional exponent is passed in as an argument
	// and is lower than the current exponent
	if optexp != nil {
		e := *optexp
		if e < exp {
			exp = e
			syms := iecExponentUnitMap[e]
			left, right = syms[0], syms[1]
		}
	}
	// Get common unit symbols based on the final exponent
	syms := findNearestIECUnitSymbols(exp)
	var li, ri int
	// Look for left and right matches to the result and save their index if
	// matched
	for i, v := range syms {
		n := i + 1
		if left == v {
			li = n
		}
		if right == v {
			ri = n
		}
	}
	// If there were no matches assume the highest symbol and return a bad
	// symbol error
	if li == 0 && ri == 0 {
		return syms[1], err
	}
	// If the left and right indexes are equal we need to choose the right
	// symbol as the greatest of the pairs for the exponent, meaning index of 1
	if li == ri {
		return syms[1], nil
	}
	// Find the valid indexes and return the greatest
	if li > 0 && ri > 0 {
		if li > ri {
			return syms[li-1], nil
		}
		return syms[ri-1], nil
	}
	if li > 0 && ri == 0 {
		if li == 1 {
			return syms[li], nil
		}
		return syms[li-1], nil
	}
	if ri == 1 && li == 0 {
		return syms[ri], nil
	}
	return syms[ri-1], nil
}

// Add attempts to add one Unit to another
func (u *IECUnit) Add(unit Unit) Unit {
	var (
		ru   *IECUnit
		ok   bool
		exp  *uint
		nsym UnitSymbol
		err  error
	)
	if ru, ok = unit.(*IECUnit); !ok {
		return u
	}
	left := u.ByteSize()
	right := ru.ByteSize()
	total := left + right
	nexp := uint(math.Round(math.Log2(total) / 10))
	if nexp > u.exponent && nexp > ru.exponent {
		e := uint(nexp)
		exp = &e
	}
	if u.symbol != ru.symbol {
		nsym, err = findGCDIECUnitSymbol(u.symbol, ru.symbol, exp)
		if err != nil {
			return u
		}
	} else {
		nsym = u.symbol
	}
	size := sizeInIECUnit(nsym, total)
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
	if u.symbol != ru.symbol || uint(nexp) != u.exponent {
		lsym, ok = FindLeastUnitSymbol(IEC, uint(nexp))
		if !ok {
			return u
		}
		gsym, ok = FindGreatestUnitSymbol(IEC, uint(nexp))
		if !ok {
			return u
		}
	} else {
		lsym, gsym = u.symbol, u.symbol
	}
	smlSize := sizeInIECUnit(lsym, total)
	lrgSize := sizeInIECUnit(gsym, total)

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
