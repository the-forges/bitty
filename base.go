package main

// Conforms to IEC and SI Standards
// Conversions taken from SI Brochure 9, EN, Chapter 3, page 143 (145)
// https://www.bipm.org/utils/common/pdf/si-brochure/SI-Brochure-9.pdf
const (
	BitSize        float64 = 1
	ByteSizeInBits         = BitSize * 8
)

// UnitSymbol represents the measurement symbol of a unit measurement as dictated by the SI and IEC
type UnitSymbol string

// SI and IEC symbols for unit measurements
const (
	Bit  UnitSymbol = "Bit"
	Byte UnitSymbol = "Byte"
	Kib  UnitSymbol = "Kib"
	Mib  UnitSymbol = "Mib"
	Gib  UnitSymbol = "Gib"
	Tib  UnitSymbol = "Tib"
	Pib  UnitSymbol = "Pib"
	Eib  UnitSymbol = "Eib"
	Zib  UnitSymbol = "Zib"
	Yib  UnitSymbol = "Yib"
	KiB  UnitSymbol = "KiB"
	MiB  UnitSymbol = "MiB"
	GiB  UnitSymbol = "GiB"
	TiB  UnitSymbol = "TiB"
	PiB  UnitSymbol = "PiB"
	EiB  UnitSymbol = "EiB"
	ZiB  UnitSymbol = "ZiB"
	YiB  UnitSymbol = "YiB"
	Kb   UnitSymbol = "Kb"
	Mb   UnitSymbol = "Mb"
	Gb   UnitSymbol = "Gb"
	Tb   UnitSymbol = "Tb"
	Pb   UnitSymbol = "Pb"
	Eb   UnitSymbol = "Eb"
	Zb   UnitSymbol = "Zb"
	Yb   UnitSymbol = "Yb"
	KB   UnitSymbol = "KB"
	MB   UnitSymbol = "MB"
	GB   UnitSymbol = "GB"
	TB   UnitSymbol = "TB"
	PB   UnitSymbol = "PB"
	EB   UnitSymbol = "EB"
	ZB   UnitSymbol = "ZB"
	YB   UnitSymbol = "YB"
)

// ExponentMap is used to quickly lookup an exponent given a unit symbol
var ExponentMap = map[UnitSymbol]int{
	Bit:  0,
	Byte: 0,
	Kib:  1,
	Mib:  2,
	Gib:  3,
	Tib:  4,
	Pib:  5,
	Eib:  6,
	Zib:  7,
	Yib:  8,
	KiB:  1,
	MiB:  2,
	GiB:  3,
	TiB:  4,
	PiB:  5,
	EiB:  6,
	ZiB:  7,
	YiB:  8,
	Kb:   1,
	Mb:   2,
	Gb:   3,
	Tb:   4,
	Pb:   5,
	Eb:   6,
	Zb:   7,
	Yb:   8,
	KB:   1,
	MB:   2,
	GB:   3,
	TB:   4,
	PB:   5,
	EB:   6,
	ZB:   7,
	YB:   8,
}

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
