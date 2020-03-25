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

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testIECUnit struct {
	Unit     IECUnit
	Expected float64
}

func ExampleNewIECUnit() {
	a, _ := NewIECUnit(10.0, Mib)
	b, _ := NewIECUnit(1.0, "GiB")
	_, cerr := NewIECUnit(3.0, "")
	_, derr := NewIECUnit(32.0, "fooBar")
	fmt.Printf("%v\n", a)
	fmt.Printf("%v\n", b)
	fmt.Printf("%v\n", cerr)
	fmt.Printf("%v\n", derr)
	// Output:
	// &{10 Mib 2}
	// &{1 GiB 3}
	// unit symbol not supported: empty symbol
	// unit symbol not supported: fooBar
}

func ExampleIECUnit_ByteSize() {
	a, _ := NewIECUnit(10.0, MiB)
	b, _ := NewIECUnit(10.0, Mib)
	fmt.Printf("%.f\n", a.ByteSize())
	fmt.Printf("%.f\n", b.ByteSize())
	// Output:
	// 10485760
	// 1310720
}

func generateTestIECUnitByteSize(t *testing.T, sym UnitSymbol) testIECUnit {
	u, err := NewIECUnit(rand.Float64(), sym)
	if err != nil {
		t.Error(err)
	}
	l := testIECUnit{Unit: *u}
	le := float64(u.exponent * 10)
	lb := float64(math.Exp2(le) * l.Unit.size)
	switch sym {
	case Bit:
		l.Expected = l.Unit.size * 8
	case Byte:
		l.Expected = l.Unit.size
	case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
		l.Expected = lb * 0.125
	case KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		l.Expected = lb
	default:
		l.Expected = float64(0)
	}
	return l
}

func TestIEC_ByteSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnit, 0, len(unitSymbolPairs))
	for _, p := range unitSymbolPairs {
		if p.Standard() != IEC {
			break
		}
		l := generateTestIECUnitByteSize(t, p.Least())
		r := generateTestIECUnitByteSize(t, p.Greatest())
		tests = append(tests, l, r)
	}
	// Add a bad entry for negative testing
	bu := testIECUnit{
		Unit:     IECUnit{rand.Float64(), UnitSymbol("FooBar"), 30},
		Expected: float64(0),
	}
	tests = append(tests, bu)
	// Run through all the tests
	for _, tst := range tests {
		assert.Equal(t, tst.Expected, tst.Unit.ByteSize())
	}
}

func ExampleIECUnit_BitSize() {
	a, _ := NewIECUnit(10.0, MiB)
	b, _ := NewIECUnit(10.0, Mib)
	fmt.Printf("%.f\n", a.BitSize())
	fmt.Printf("%.f\n", b.BitSize())
	// Output:
	// 83886080
	// 10485760
}

func generateTestIECUnitBitSize(t *testing.T, sym UnitSymbol) testIECUnit {
	tu, err := NewIECUnit(rand.Float64()*10, sym)
	if err != nil {
		t.Error(err)
	}
	l := testIECUnit{Unit: *tu}
	bytes := tu.ByteSize()
	switch sym {
	case Bit:
		l.Expected = l.Unit.size
	case Byte,
		Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib,
		KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
		l.Expected = float64(bytes * 8)
	default:
		l.Expected = float64(0)
	}
	return l
}

func TestIECUnit_BitSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnit, 0, len(unitSymbolPairs))
	// Setup test cases based out of what is in IECUnitExponentMap
	for _, p := range unitSymbolPairs {
		if p.Standard() != IEC {
			break
		}
		l := generateTestIECUnitBitSize(t, p.Least())
		r := generateTestIECUnitBitSize(t, p.Greatest())
		tests = append(tests, l, r)
	}
	// Add a bad entry for negative testing
	bu := testIECUnit{
		Unit:     IECUnit{rand.Float64() * 10, UnitSymbol("FooBar"), 30},
		Expected: float64(0),
	}
	tests = append(tests, bu)
	// Run through all the tests
	for _, tst := range tests {
		assert.Equal(t, tst.Expected, tst.Unit.BitSize())
		if t.Failed() {
			fmt.Printf("size: %f, symbol: %s, bits: %f, expected: %f\n",
				tst.Unit.size, tst.Unit.symbol,
				tst.Unit.BitSize(), tst.Expected,
			)
		}
	}
}

func ExampleIECUnit_SizeInUnit() {
	a, _ := NewIECUnit(10.0, MiB)
	inKiB := a.SizeInUnit(KiB)
	inGiB := a.SizeInUnit(GiB)
	inMib := a.SizeInUnit(Mib)
	fmt.Println(inKiB, inGiB, inMib)
	// Output:
	// 10240 0.009765625 80
}

type testIECSizeInUnit struct {
	unit     IECUnit
	to       UnitSymbol
	expected float64
}

func generateTestIECUnitSizeInUnit(t *testing.T, unit IECUnit, sym UnitSymbol) testIECSizeInUnit {
	u := testIECSizeInUnit{unit: unit, to: sym}
	r, err := NewIECUnit(unit.size, sym)
	if err != nil {
		t.Error(err)
	}
	var (
		left    = unit.ByteSize()
		right   = r.ByteSize()
		diffExp = float64(unit.exponent - r.exponent)
	)
	if diffExp > 0 {
		u.expected = right * diffExp
	} else {
		u.expected = (left / right) * u.unit.size
	}
	return u
}

func TestIECUnit_SizeInUnit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECSizeInUnit, 0, len(unitSymbolPairs))
	for _, p := range unitSymbolPairs {
		if p.Standard() != IEC {
			break
		}
		l, err := NewIECUnit(rand.Float64()*10, p.Least())
		if err != nil {
			t.Error(err)
			break
		}
		r, err := NewIECUnit(rand.Float64()*10, p.Greatest())
		if err != nil {
			t.Error(err)
			break
		}
		for _, rp := range unitSymbolPairs {
			lu := generateTestIECUnitSizeInUnit(t, *l, rp.Least())
			ru := generateTestIECUnitSizeInUnit(t, *r, rp.Greatest())
			tests = append(tests, lu, ru)
		}
	}
	// Add a couple of bad entries for negative testing
	bu := testIECSizeInUnit{
		unit:     IECUnit{rand.Float64() * 10, UnitSymbol("FooBar"), 30},
		to:       MiB,
		expected: float64(0),
	}
	bur := testIECSizeInUnit{
		unit:     IECUnit{rand.Float64() * 10, MiB, 30},
		to:       UnitSymbol("FooBar"),
		expected: float64(0),
	}
	tests = append(tests, bu, bur)
	// Run through all the tests
	for _, tst := range tests {
		assert.Equal(t, tst.expected, tst.unit.SizeInUnit(tst.to))
	}
}

func ExampleIECUnit_Add() {
	// Test the same byte symbol
	a, _ := NewIECUnit(2, MiB)
	b, _ := NewIECUnit(2, MiB)
	c, _ := a.Add(b)
	fmt.Printf(
		"%.f %s + %.f %s = %.f %s\n",
		a.Size(), a.Symbol(),
		b.Size(), b.Symbol(),
		c.Size(), c.Symbol(),
	)
	// Output:
	// 2 MiB + 2 MiB = 4 MiB
}

type testIECUnitAdd struct {
	left, right, expected *IECUnit
}

func TestIECUnit_Add(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnitAdd, 0, len(iecUnitExponentMap))
	// Setup test cases based out of what is in IECUnitExponentMap
	for k := range iecUnitExponentMap {
		tul, _ := NewIECUnit(rand.Float64()*10, k)
		if tul == nil {
			break
		}
		for l := range iecUnitExponentMap {
			var (
				ru   *IECUnit
				nexp int
				nsym UnitSymbol
				size float64
			)
			ru, _ = NewIECUnit(rand.Float64()*10, l)
			u := testIECUnitAdd{left: tul, right: ru}
			lok, rok := ValidateSymbols(tul.Symbol(), ru.Symbol())
			if !lok && !rok {
				nexp = 0
			}
			if lok && !rok {
				nexp = tul.Exponent()
			}
			if rok && !lok {
				nexp = ru.Exponent()
			}
			left := tul.ByteSize()
			right := ru.ByteSize()
			total := left + right
			if total > 0 {
				nexp = int(math.Round(math.Log2(total) / 10))
			}
			if tul.Exponent() >= ru.Exponent() {
				nexp = tul.Exponent()
			} else {
				nexp = ru.Exponent()
			}
			lsym, ok := FindLeastUnitSymbol(IEC, nexp)
			gsym, ok := FindGreatestUnitSymbol(IEC, nexp)
			if !ok {
				u.expected, _ = NewIECUnit(0, Byte)
			} else {
				smallSize := BytesToUnitSymbolSize(IEC, lsym, total)
				lrgSize := BytesToUnitSymbolSize(IEC, gsym, total)
				if lrgSize < 1 {
					nsym = lsym
					size = smallSize
				} else {
					nsym = gsym
					size = lrgSize
				}
				u.expected, _ = NewIECUnit(size, nsym)
			}
			tests = append(tests, u)
		}
	}
	// Add a couple of bad entries for negative testing
	s := rand.Float64() * 10
	gu, _ := NewIECUnit(s, MiB)
	byteu, _ := NewIECUnit(0, Byte)
	bu := &IECUnit{s, UnitSymbol("FooBar"), 30}
	bul := testIECUnitAdd{
		left:     bu,
		right:    gu,
		expected: gu,
	}
	bur := testIECUnitAdd{
		left:     gu,
		right:    bu,
		expected: gu,
	}
	bub := testIECUnitAdd{
		left:     bu,
		right:    bu,
		expected: byteu,
	}
	tests = append(tests, bul, bur, bub)
	// Run through all the tests
	for _, tst := range tests {
		u, _ := tst.left.Add(tst.right)
		assert.Equal(t, tst.expected, u)
	}
}

func ExampleIECUnit_Subtract() {
	var (
		c  *IECUnit
		ok bool
	)
	// Test the same byte symbol
	a, _ := NewIECUnit(10, GiB)
	b, _ := NewIECUnit(10.023, GiB)
	c, ok = a.Subtract(b).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", c))
	}
	fmt.Printf(
		"%.3f %s - %.3f %s = %.3f %s\n",
		a.size, a.symbol,
		b.size, b.symbol,
		c.size, c.symbol,
	)
	// Output:
	// 10.000 GiB - 10.023 GiB = -23.552 MiB
}

type testIECUnitSubtract struct {
	left, right, expected *IECUnit
}

func TestIECUnit_Subtract(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnitSubtract, 0, len(iecUnitExponentMap))
	// Setup test cases based out of what is in IECUnitExponentMap
	for k := range iecUnitExponentMap {
		tul, _ := NewIECUnit(rand.Float64()*10, k)
		if tul == nil {
			break
		}
		for l := range iecUnitExponentMap {
			var (
				ru               *IECUnit
				nexp             int
				nsym, lsym, gsym UnitSymbol
				total            float64
				neg              bool
			)
			ru, _ = NewIECUnit(rand.Float64()*10, l)
			u := testIECUnitSubtract{left: tul, right: ru}
			lok, rok := ValidateSymbols(tul.Symbol(), ru.Symbol())
			if !lok && !rok {
				nexp = 0
			}
			if lok && !rok {
				nexp = tul.Exponent()
			}
			if rok && !lok {
				nexp = ru.Exponent()
			}
			left := tul.ByteSize()
			right := ru.ByteSize()
			if left >= right {
				total = left - right
			} else {
				total = right - left
				neg = true
			}
			if total > 0 {
				nexp = int(math.Round(math.Log2(total) / 10))
			}
			lsym, _ = FindLeastUnitSymbol(IEC, nexp)
			gsym, _ = FindGreatestUnitSymbol(IEC, nexp)
			smlSize := BytesToUnitSymbolSize(IEC, lsym, total)
			lrgSize := BytesToUnitSymbolSize(IEC, gsym, total)
			if lrgSize > 0 {
				if neg {
					lrgSize = -lrgSize
				}
				u.expected, _ = NewIECUnit(lrgSize, gsym)
			} else {
				if neg {
					smlSize = -smlSize
				}
				u.expected, _ = NewIECUnit(smlSize, nsym)
			}

			tests = append(tests, u)
		}
	}
	// Add a couple of bad entries for negative testing
	s := rand.Float64() * 10
	gu, _ := NewIECUnit(s, MiB)
	byteu, _ := NewIECUnit(0, Byte)
	bu := &IECUnit{s, UnitSymbol("FooBar"), 30}
	bul := testIECUnitSubtract{
		left:     bu,
		right:    gu,
		expected: gu,
	}
	bur := testIECUnitSubtract{
		left:     gu,
		right:    bu,
		expected: gu,
	}
	bub := testIECUnitSubtract{
		left:     bu,
		right:    bu,
		expected: byteu,
	}
	tests = append(tests, bul, bur, bub)
	// Run through all the tests
	for _, tst := range tests {
		u, ok := tst.left.Subtract(tst.right).(*IECUnit)
		assert.Equal(t, true, ok)
		assert.Equal(t, tst.expected, u)
	}
}
