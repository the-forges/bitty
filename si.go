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
