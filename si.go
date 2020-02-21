package main

var siUnitExponentMap = map[UnitSymbol]int{
	Bit:  0,
	Byte: 0,
	Kb:   1,
	KB:   1,
	Mb:   2,
	MB:   2,
	Gb:   3,
	GB:   3,
	Tb:   4,
	TB:   4,
	Pb:   5,
	PB:   5,
	Eb:   6,
	EB:   6,
	Zb:   7,
	ZB:   7,
	Yb:   8,
	YB:   8,
}

var siExponentUnitMap = [][]UnitSymbol{
	[]UnitSymbol{Bit, Byte},
	[]UnitSymbol{Kb, KB},
	[]UnitSymbol{Mb, MB},
	[]UnitSymbol{Gb, GB},
	[]UnitSymbol{Tb, TB},
	[]UnitSymbol{Pb, PB},
	[]UnitSymbol{Eb, EB},
	[]UnitSymbol{Zb, ZB},
	[]UnitSymbol{Yb, YB},
}

// SIUnit handles deca units as dictated by IEC Standards
type SIUnit float64
