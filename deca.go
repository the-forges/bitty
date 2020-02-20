package main

// Conversions for base 10 measurements
const (
	KBSizeInBytes = ByteSizeInBits * 1e3
	MBSizeInBytes = KBSizeInBytes * 1e6
	GBSizeInBytes = MBSizeInBytes * 1e9
	TBSizeInBytes = GBSizeInBytes * 1e12
	PBSizeInBytes = TBSizeInBytes * 1e15
	EBSizeInBytes = PBSizeInBytes * 1e18
	ZBSizeInBytes = EBSizeInBytes * 1e21
	YBSizeInBytes = ZBSizeInBytes * 1e24
)

// SIUnit handles deca units as dictated by IEC Standards
type SIUnit float64
