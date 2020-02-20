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

// IECUnit handles binary units as dictated by IEC Standards
type IECUnit struct {
	// size is the size as measured by the symbol (UnitSymbol), which is equivalent to:
	// 		b(2^10)^n(1/8)
	// 		B(2^10)^n
	size     float64
	symbol   UnitSymbol
	exponent int
}

// NewIECUnit returns a *IECUnit with the proper exponent included
func NewIECUnit(size float64, symbol UnitSymbol) (*IECUnit, error) {
	if exp, ok := iecUnitExponentMap[symbol]; ok {
		return &IECUnit{size, symbol, exp}, nil
	}
	return nil, NewErrUnitSymbolNotSupported(symbol)
}

// IECSymbolBitSize takes a UnitSymbol and some bytes, returning the calculated
// bits
func IECSymbolBitSize(symbol UnitSymbol, bytes float64) float64 {
	switch symbol {
	case Bit:
		return float64(bytes / 8)
	case Byte,
		Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib,
		KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		return float64(bytes * 8)
	default:
		return float64(0)
	}
}

// BitSize returns the size of the Unit measured in bits
func (u *IECUnit) BitSize() float64 {
	return IECSymbolBitSize(u.symbol, u.ByteSize())
}

// IECSymbolByteSize takes a UnitSymbol and a size, returning the calculated
// bytes
func IECSymbolByteSize(symbol UnitSymbol, size float64) float64 {
	uexp, ok := iecUnitExponentMap[symbol]
	if !ok {
		return float64(0)
	}
	exp := float64(uexp * 10)
	bytes := float64(math.Exp2(exp) * size)
	switch symbol {
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
	return IECSymbolByteSize(u.symbol, u.size)
}

// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
func (u *IECUnit) SizeInUnit(symbol UnitSymbol) float64 {
	_, uok := iecUnitExponentMap[u.symbol]
	exp, ok := iecUnitExponentMap[symbol]
	if uok && ok {
		var (
			diffExp = float64(u.exponent - exp)
			left    = IECSymbolByteSize(u.symbol, u.size)
			right   = IECSymbolByteSize(symbol, u.size)
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

func findNearestIECUnitSymbols(exp int) []UnitSymbol {
	return iecExponentUnitMap[exp]
}

func findLargestIECUnitSymbol(left, right UnitSymbol, exp int) UnitSymbol {
	if left == right {
		return left
	}
	syms := findNearestIECUnitSymbols(exp)
	if syms[1] == left {
		return left
	}
	return right
}

// Add attempts to add one Unit to another
func (u *IECUnit) Add(unit Unit) Unit {
	var (
		ru *IECUnit
		ok bool
	)
	if ru, ok = unit.(*IECUnit); !ok {
		return u
	}
	left := u.ByteSize()
	right := ru.ByteSize()
	total := left + right
	nexp := int(math.Round(math.Log2(total) / 10))
	nsym := findLargestIECUnitSymbol(u.symbol, ru.symbol, nexp)
	size := sizeInIECUnit(nsym, total)
	nu, _ := NewIECUnit(size, nsym)
	return nu
}

// Subtract attempts to subtract one Unit from another
func (u *IECUnit) Subtract(units Unit) Unit {
	return nil
}

// Multiply attempts to multiply one Unit by another
func (u *IECUnit) Multiply(unit Unit) Unit {
	return nil
}

// Divide attempts to divide one Unit by another
func (u *IECUnit) Divide(unit Unit) Unit {
	return nil
}
