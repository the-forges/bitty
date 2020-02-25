package main

// Unit enables Unit kinds to interact with each other
type Unit interface {
	Sizer
	Calculator
}

// Sizer enables Unit to get a size measured by bit, byte, or arbitrary data Unit kind
type Sizer interface {
	// BitSize returns the size of the Unit measured in bits
	BitSize() float64
	// ByteSize returns the size of the Unit measured in bytes
	ByteSize() float64
	// SizeInUnit returns the size of the Unit measured in an arbitrary UnitSymbol from Bit up to YiB or YB
	SizeInUnit(UnitSymbol) float64
}

// Calculator enables Units to be calculated against each other
// All returns are diminshing or increasing UnitSymbol measurements as defined by the SI and IEC
type Calculator interface {
	// Add attempts to add one Unit to another
	Add(Unit) Unit
	// Subtract attempts to subtract one Unit from another
	Subtract(Unit) Unit
	// Multiply attempts to multiply one Unit by another
	Multiply(Unit) Unit
	// Divide attempts to divide one Unit from another
	Divide(Unit) Unit
}

// BaseSymbolPair represents the bit and byte pairs necassary for
type BaseUnitSymbolPair struct {
	standard *UnitStandard
}

func NewBaseUnitSymbolPair(std UnitStandard) UnitSymbolPair {
	return &BaseUnitSymbolPair{standard: &std}
}
func (b *BaseUnitSymbolPair) Standard() UnitStandard {
	if b.standard == nil {
		return SI
	}
	return *b.standard
}
func (b *BaseUnitSymbolPair) Exponent() uint {
	return 0
}
func (b *BaseUnitSymbolPair) Least() UnitSymbol {
	return Bit
}
func (b *BaseUnitSymbolPair) Greatest() UnitSymbol {
	return Byte
}
