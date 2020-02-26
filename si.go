package main

import "math"

// IECUnitSymbolPair represents a binary unit symbol pair as defined by the 9th
// edition SI standard
type SIUnitSymbolPair struct {
	least, greatest UnitSymbol
	exponent        int
}

func NewSIUnitSymbolPair(l, r UnitSymbol, e int) UnitSymbolPair {
	return &SIUnitSymbolPair{least: l, greatest: r, exponent: e}
}

func (pair *SIUnitSymbolPair) Standard() UnitStandard {
	return SI
}

func (pair *SIUnitSymbolPair) Exponent() int {
	return pair.exponent
}

func (pair *SIUnitSymbolPair) Least() UnitSymbol {
	return pair.least
}

func (pair *SIUnitSymbolPair) Greatest() UnitSymbol {
	return pair.greatest
}

// SIUnit handles deca units as dictated by SI Standards
type SIUnit struct {
	// size is the size as measured by the symbol (UnitSymbol), which is equivalent to:
	// 		b(2^10)^n(1/8)
	// 		B(2^10)^n
	size     float64
	symbol   UnitSymbol
	exponent int
}

// NewSIUnit returns a *SIUnit with the proper exponent included
func NewSIUnit(size float64, symbol UnitSymbol) (*SIUnit, error) {
	// if exp, ok := siUnitExponentMap[symbol]; ok {
	// 	return &SIUnit{size, symbol, exp}, nil
	// }
	return nil, NewErrUnitSymbolNotSupported(symbol)
}

// SISymbolBitSize takes a UnitSymbol and some bytes, returning the calculated
// bits
func SISymbolBitSize(symbol UnitSymbol, bytes float64) float64 {
	switch symbol {
	case Bit:
		return float64(bytes / 8)
	case Byte,
		Kb, Mb, Gb, Tb, Pb, Eb, Zb, Yb,
		KB, MB, GB, TB, PB, EB, ZB, YB:
		return float64(bytes * 8)
	default:
		return float64(0)
	}
}

// BitSize returns the size of the Unit measured in bits
func (u *SIUnit) BitSize() float64 {
	return SISymbolBitSize(u.symbol, u.ByteSize())
}

// SISymbolByteSize takes a UnitSymbol and a size, returning the calculated
// bytes
func SISymbolByteSize(symbol UnitSymbol, size float64) float64 {
	uexp := 0
	exp := float64(uexp * 10)
	bytes := float64(math.Exp2(exp) * size)
	switch symbol {
	case Bit:
		return float64(size * 8)
	case Byte:
		return float64(size)
	case Kb, Mb, Gb, Tb, Pb, Eb, Zb, Yb:
		return float64(bytes * 0.125)
	case KB, MB, GB, TB, PB, EB, ZB, YB:
		return float64(bytes)
	default:
		return float64(0)
	}
}

func sizeInSIUnit(symbol UnitSymbol, bytes float64) float64 {
	return float64(0)
}

// ByteSize returns the size of the Unit measured in bytes
func (u *SIUnit) ByteSize() float64 {
	return float64(0)
}

// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
func (u *SIUnit) SizeInUnit(symbol UnitSymbol) float64 {
	return float64(0)
}

// Add attempts to add one Unit to another
func (u *SIUnit) Add(unit Unit) Unit {
	return nil
}

// Subtract attempts to subtract one Unit from another
func (u *SIUnit) Subtract(unit Unit) Unit {
	return nil
}

// Multiply attempts to multiply one Unit by another
func (u *SIUnit) Multiply(unit Unit) Unit {
	return nil
}

// Divide attempts to divide one Unit by another
func (u *SIUnit) Divide(unit Unit) Unit {
	return nil
}
