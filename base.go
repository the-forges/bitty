package bitty

/*
	Copyright 2020 IBM

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

// Symbolic enables Unit to provide a standard, exponent, and symbol
type Symbolic interface {
	// Standard returns the standard for the unit as defined by the SI brochure,
	// 9th Edition, page 145:
	// 	https://www.bipm.org/utils/common/pdf/si-brochure/SI-Brochure-9.pdf
	Standard() UnitStandard
	// Exponent returns the supported exponent as an int
	// To calculate the value from the standard, symbol, and size the formulas
	// are:
	// 	- Given the Standard is SI, value as v is equal to (10^e)size
	// 	- Given the Standard is IEC, value as v is equal to (2^(e*10))size
	Exponent() int
	// Symbol returns the supported symbol as a UnitSymbol
	Symbol() UnitSymbol
}

// Sizer enables Unit to get a size measured by bit, byte, or arbitrary data
// Unit kind
type Sizer interface {
	// Size returns the size of the Unit
	Size() float64
	// BitSize returns the size of the Unit measured in bits
	BitSize() float64
	// ByteSize returns the size of the Unit measured in bytes
	ByteSize() float64
	// SizeInUnit returns the size of the Unit measured in an arbitrary
	// UnitSymbol from Bit up to YiB or YB
	SizeInUnit(UnitSymbol) float64
}

// Calculator enables Units to be calculated against each other
// All returns are diminshing or increasing UnitSymbol measurements as defined
// by the SI and IEC
type Calculator interface {
	// Add attempts to add one Unit to another
	Add(Unit) (Unit, error)
	// Subtract attempts to subtract one Unit from another
	Subtract(Unit) Unit
	// Multiply attempts to multiply one Unit by another
	Multiply(Unit) Unit
	// Divide attempts to divide one Unit from another
	Divide(Unit) Unit
}

// Unit enables Unit kinds to interact with each other
type Unit interface {
	Symbolic
	Sizer
	Calculator
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
func (b *BaseUnitSymbolPair) Exponent() int {
	return 0
}
func (b *BaseUnitSymbolPair) Least() UnitSymbol {
	return Bit
}
func (b *BaseUnitSymbolPair) Greatest() UnitSymbol {
	return Byte
}
