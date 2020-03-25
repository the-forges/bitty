package bitty

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

type testSIUnit struct {
	Unit     SIUnit
	Expected float64
}

func ExampleNewSIUnit() {
	a, _ := NewSIUnit(10.0, MB)
	fmt.Printf("%v\n", a)
	// Output:
	// &{10 MB 6}
}

func ExampleSIUnit_ByteSize() {
	a, _ := NewSIUnit(10.0, kB)
	fmt.Printf("%.f\n", a.ByteSize())
	// Output:
	// 10000
}

func generateTestSIUnitByteSize(t *testing.T, sym UnitSymbol) testSIUnit {
	u, err := NewSIUnit(rand.Float64(), sym)
	if err != nil {
		t.Error(err)
	}
	l := testSIUnit{Unit: *u}
	le := float64(u.exponent * 10)
	lb := float64(math.Exp2(le) * l.Unit.size)
	switch sym {
	case Bit:
		l.Expected = l.Unit.size * 8
	case Byte:
		l.Expected = l.Unit.size
	case db, hb, kb, Mb, Gb, Tb, Pb, Eb, Zb, Yb:
		l.Expected = lb * 0.125
	case dB, hB, kB, MB, GB, TB, PB, EB, ZB, YB:
		l.Expected = lb
	default:
		l.Expected = float64(0)
	}
	return l
}

func TestSIUnit_ByteSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testSIUnit, 0, len(unitSymbolPairs))
	for _, p := range unitSymbolPairs {
		if p.Standard() != SI {
			break
		}
		l := generateTestSIUnitByteSize(t, p.Least())
		r := generateTestSIUnitByteSize(t, p.Greatest())
		tests = append(tests, l, r)
	}
	bu := testSIUnit{
		Unit:     SIUnit{rand.Float64(), UnitSymbol("FooBar"), 30},
		Expected: float64(0),
	}
	tests = append(tests, bu)
	for _, tst := range tests {
		assert.Equal(t, tst.Expected, tst.Unit.ByteSize())
	}
}

func ExampleSIUnit_Add() {
	// Test the same byte symbol
	a, _ := NewSIUnit(2, MB)
	b, _ := NewSIUnit(2, MB)
	c, _ := a.Add(b)
	fmt.Printf(
		"%.f %s + %.f %s = %.f %s\n",
		a.Size(), a.Symbol(),
		b.Size(), b.Symbol(),
		c.Size(), c.Symbol(),
	)
	// Output:
	// 2 MB + 2 MB = 4 MB
}

type testSIUnitAdd struct {
	left, right, expected *SIUnit
}

func generateTestSIUnitAdd(lu, ru *SIUnit) testSIUnitAdd {
	var (
		exp  int
		sym  UnitSymbol
		size float64
		tu   = testSIUnitAdd{left: lu, right: ru}
	)
	left, right := lu.ByteSize(), ru.ByteSize()
	total := left + right
	if total > 0 {
		exp = int(math.Round(math.Log2(total) / 10))
	}
	if lu.Exponent() >= ru.Exponent() {
		exp = lu.Exponent()
	} else {
		exp = ru.Exponent()
	}
	lsym, ok := FindLeastUnitSymbol(SI, exp)
	gsym, ok := FindGreatestUnitSymbol(SI, exp)
	if !ok {
		tu.expected, _ = NewSIUnit(0, Byte)
		return tu
	}
	smlSize := BytesToUnitSymbolSize(SI, lsym, total)
	lrgSize := BytesToUnitSymbolSize(SI, gsym, total)
	if lrgSize < 1 {
		sym = lsym
		size = smlSize
	} else {
		sym = gsym
		size = lrgSize
	}
	tu.expected, _ = NewSIUnit(size, sym)
	return tu
}

func TestSIUnit_Add(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testSIUnitAdd, 0, len(unitSymbolPairs))
	// Setup test cases based out of what is in SIUnitExponentMap
	for _, p := range unitSymbolPairs {
		if p.Standard() != SI {
			break
		}
		olu, _ := NewSIUnit(rand.Float64()*10, p.Least())
		ogu, _ := NewSIUnit(rand.Float64()*10, p.Greatest())
		for _, o := range unitSymbolPairs {
			if o.Standard() != SI {
				break
			}
			ilu, _ := NewSIUnit(rand.Float64()*10, o.Least())
			igu, _ := NewSIUnit(rand.Float64()*10, o.Greatest())
			tlu := generateTestSIUnitAdd(olu, ilu)
			tgu := generateTestSIUnitAdd(ogu, igu)
			tests = append(tests, tlu, tgu)
		}
	}
	// Add a couple of bad entries for negative testing
	s := rand.Float64() * 10
	gu, _ := NewSIUnit(s, MB)
	byteu, _ := NewSIUnit(0, Byte)
	bu := &SIUnit{s, UnitSymbol("FooBar"), 30}
	bul := testSIUnitAdd{
		left:     bu,
		right:    gu,
		expected: gu,
	}
	bur := testSIUnitAdd{
		left:     gu,
		right:    bu,
		expected: gu,
	}
	bub := testSIUnitAdd{
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

func ExampleSIUnit_Subtract() {
	var (
		c  *SIUnit
		ok bool
	)
	// Test the same byte symbol
	a, _ := NewSIUnit(10, MB)
	b, _ := NewSIUnit(10.023, MB)
	c, ok = a.Subtract(b).(*SIUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *SIUnit: %v", c))
	}
	fmt.Printf(
		"%.3f %s - %.3f %s = %.3f %s\n",
		a.size, a.symbol,
		b.size, b.symbol,
		c.size, c.symbol,
	)
	// Output:
	// 10.000 MB - 10.023 MB = -23.000 kB
}

type testSIUnitSubtract struct {
	left, right, expected *SIUnit
}

func generateTestSIUnitSubtract(lu, ru *SIUnit) testSIUnitSubtract {
	var (
		exp              int
		nsym, lsym, gsym UnitSymbol
		total            float64
		neg, ok          bool
		tu               = testSIUnitSubtract{left: lu, right: ru}
	)
	left, right := lu.ByteSize(), ru.ByteSize()
	if left >= right {
		total = left - right
	} else {
		total = right - left
		neg = true
	}
	if total > 0 {
		exp = int(math.Floor(math.Log10(total)))
	}
	lsym, ok = FindLeastUnitSymbol(SI, exp)
	gsym, ok = FindGreatestUnitSymbol(SI, exp)
	if !ok {
		tu.expected, _ = NewSIUnit(0, Byte)
		return tu
	}
	smlSize := BytesToUnitSymbolSize(SI, lsym, total)
	lrgSize := BytesToUnitSymbolSize(SI, gsym, total)
	if lrgSize > 0 {
		if neg {
			lrgSize = -lrgSize
		}
		tu.expected, _ = NewSIUnit(lrgSize, gsym)
	} else {
		if neg {
			smlSize = -smlSize
		}
		tu.expected, _ = NewSIUnit(smlSize, nsym)
	}
	return tu
}

func TestSIUnit_Subtract(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testSIUnitSubtract, 0, len(unitSymbolPairs))
	// Setup test cases based out of what is in SIUnitExponentMap
	for _, p := range unitSymbolPairs {
		if p.Standard() != SI {
			break
		}
		olu, _ := NewSIUnit(rand.Float64()*10, p.Least())
		ogu, _ := NewSIUnit(rand.Float64()*10, p.Greatest())
		for _, o := range unitSymbolPairs {
			if o.Standard() != SI {
				break
			}
			ilu, _ := NewSIUnit(rand.Float64()*10, o.Least())
			igu, _ := NewSIUnit(rand.Float64()*10, o.Greatest())
			tlu := generateTestSIUnitSubtract(olu, ilu)
			tgu := generateTestSIUnitSubtract(ogu, igu)
			tests = append(tests, tlu, tgu)
		}
	}
	// Add a couple of bad entries for negative testing
	s := rand.Float64() * 10
	gu, _ := NewSIUnit(s, MB)
	byteu, _ := NewSIUnit(0, Byte)
	bu := &SIUnit{s, UnitSymbol("FooBar"), 30}
	bul := testSIUnitSubtract{
		left:     bu,
		right:    gu,
		expected: gu,
	}
	bur := testSIUnitSubtract{
		left:     gu,
		right:    bu,
		expected: gu,
	}
	bub := testSIUnitSubtract{
		left:     bu,
		right:    bu,
		expected: byteu,
	}
	tests = append(tests, bul, bur, bub)
	// Run through all the tests
	for _, tst := range tests {
		u, ok := tst.left.Subtract(tst.right).(*SIUnit)
		assert.Equal(t, true, ok)
		assert.Equal(t, tst.expected, u)
	}
}
