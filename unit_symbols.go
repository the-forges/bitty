package main

// UnitSymbol represents the measurement symbol of a binary measurement as dictated by the SI
type UnitSymbol string

// SI symbols for binary measurements including notation for former IEC measurements
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

// UnitStandard represents a standard for unit measurement. Currently SI 9th
// edition is the supported standard, with SI notation for IEC binary and
// decimal formats
type UnitStandard int

// Unit Standard enums
const (
	SI UnitStandard = iota
	IEC
)

// UnitSymbolPair holds the least and greatest UnitSymbol for a given standard
// and exponent
type UnitSymbolPair interface {
	Standard() UnitStandard
	Exponent() int
	Least() UnitSymbol
	Greatest() UnitSymbol
}

var unitSymbolPairs = []UnitSymbolPair{
	NewBaseUnitSymbolPair(SI),
	NewBaseUnitSymbolPair(IEC),
	NewIECUnitSymbolPair(Kib, KiB, 1),
	NewIECUnitSymbolPair(Mib, MiB, 2),
	NewIECUnitSymbolPair(Gib, GiB, 3),
	NewIECUnitSymbolPair(Tib, TiB, 4),
	NewIECUnitSymbolPair(Pib, PiB, 5),
	NewIECUnitSymbolPair(Eib, EiB, 6),
	NewIECUnitSymbolPair(Zib, ZiB, 7),
	NewIECUnitSymbolPair(Yib, YiB, 8),
	NewSIUnitSymbolPair(Kb, KB, 1),
	NewSIUnitSymbolPair(Mb, MB, 2),
	NewSIUnitSymbolPair(Gb, GB, 3),
	NewSIUnitSymbolPair(Tb, TB, 4),
	NewSIUnitSymbolPair(Pb, PB, 5),
	NewSIUnitSymbolPair(Eb, EB, 6),
	NewSIUnitSymbolPair(Zb, ZB, 7),
	NewSIUnitSymbolPair(Yb, YB, 8),
}
