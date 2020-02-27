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
	db   UnitSymbol = "db"
	hb   UnitSymbol = "hb"
	kb   UnitSymbol = "kb"
	Mb   UnitSymbol = "Mb"
	Gb   UnitSymbol = "Gb"
	Tb   UnitSymbol = "Tb"
	Pb   UnitSymbol = "Pb"
	Eb   UnitSymbol = "Eb"
	Zb   UnitSymbol = "Zb"
	Yb   UnitSymbol = "Yb"
	dB   UnitSymbol = "dB"
	hB   UnitSymbol = "hB"
	kB   UnitSymbol = "kB"
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
	NewSIUnitSymbolPair(db, dB, 1),
	NewSIUnitSymbolPair(hb, hB, 2),
	NewSIUnitSymbolPair(kb, kB, 3),
	NewSIUnitSymbolPair(Mb, MB, 6),
	NewSIUnitSymbolPair(Gb, GB, 9),
	NewSIUnitSymbolPair(Tb, TB, 12),
	NewSIUnitSymbolPair(Pb, PB, 15),
	NewSIUnitSymbolPair(Eb, EB, 18),
	NewSIUnitSymbolPair(Zb, ZB, 21),
	NewSIUnitSymbolPair(Yb, YB, 24),
}
