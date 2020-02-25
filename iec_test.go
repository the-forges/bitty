package main

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

func Test_findNearestIECUnitSymbols(t *testing.T) {
	for i, v := range iecExponentUnitMap {
		u := findNearestIECUnitSymbols(uint(i))
		assert.Equal(t, v, u)
	}
}

func Test_findGCDIECUnitSymbol_knownExponents(t *testing.T) {
	var (
		res UnitSymbol
		err error
	)
	for i, outerSyms := range iecExponentUnitMap {
		outerL, outerR := outerSyms[0], outerSyms[1]
		// Test that the right is always greater than the left for an equal exponent
		res, err = findGCDIECUnitSymbol(outerL, outerR, nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		res, err = findGCDIECUnitSymbol(outerR, outerL, nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		// Test for bad unit symbol
		res, err = findGCDIECUnitSymbol(UnitSymbol("BADFOO"), outerR, nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		res, err = findGCDIECUnitSymbol(outerR, UnitSymbol("BADFOO"), nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		res, err = findGCDIECUnitSymbol(UnitSymbol("BADFOO"), outerL, nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		res, err = findGCDIECUnitSymbol(outerL, UnitSymbol("BADFOO"), nil)
		assert.NoError(t, err)
		assert.Equal(t, outerR, res)
		// Test for all known combinations of iec unit symbols
		for j, innerSyms := range iecExponentUnitMap {
			// Test that the greater exponent in a mixed set always wins
			innerL, innerR := innerSyms[0], innerSyms[1]
			if i > j {
				res, err := findGCDIECUnitSymbol(outerL, innerL, nil)
				assert.NoError(t, err)
				assert.Equal(t, innerR, res)
				res, err = findGCDIECUnitSymbol(outerR, innerR, nil)
				assert.NoError(t, err)
				assert.Equal(t, innerR, res)
			} else {
				res, err := findGCDIECUnitSymbol(outerL, innerL, nil)
				assert.NoError(t, err)
				assert.Equal(t, outerR, res)
				res, err = findGCDIECUnitSymbol(outerR, innerR, nil)
				assert.NoError(t, err)
				assert.Equal(t, outerR, res)
			}

		}
	}
}

func Test_findGCDIECUnitSymbol_calculatedExponents(t *testing.T) {
	exp := 0
	for i, outerSyms := range iecExponentUnitMap {
		outerL, outerR := outerSyms[0], outerSyms[1]
		for j, innerSyms := range iecExponentUnitMap {
			innerL, innerR := innerSyms[0], innerSyms[1]
			var expected UnitSymbol
			if i >= exp && j >= exp {
				syms := iecExponentUnitMap[exp]
				expected = syms[1]
			} else if i < exp && j >= exp {
				expected = outerR
			} else if i >= exp && j < exp {
				expected = innerR
			} else {
				if i < j {
					expected = outerR
				} else {
					expected = innerR
				}
			}
			nexp := uint(exp)

			res, err := findGCDIECUnitSymbol(outerL, innerL, &nexp)
			assert.NoErrorf(t, err, "i: %d, j: %d, nexp: %d", i, j, nexp)
			assert.Equal(t, expected, res)

			res, err = findGCDIECUnitSymbol(outerR, innerR, &nexp)
			assert.NoErrorf(t, err, "i: %d, j: %d, nexp: %d", i, j, nexp)
			assert.Equal(t, expected, res)
		}
	}
}

func ExampleIECUnit_Add() {
	// Test the same byte symbol
	a, _ := NewIECUnit(2, MiB)
	b, _ := NewIECUnit(2, MiB)
	c, ok := a.Add(b).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", c))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.f %s\n",
		a.size, a.symbol,
		b.size, b.symbol,
		c.size, c.symbol,
	)
	// Test the same bit symbol
	d, _ := NewIECUnit(2, Mib)
	e, _ := NewIECUnit(2, Mib)
	f, ok := d.Add(e).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", f))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.f %s\n",
		d.size, d.symbol,
		e.size, e.symbol,
		f.size, f.symbol,
	)
	// Test mixed bit/byte symbol
	g, _ := NewIECUnit(2, Mib)
	h, _ := NewIECUnit(2, MiB)
	i, ok := g.Add(h).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", i))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.2f %s\n",
		g.size, g.symbol,
		h.size, h.symbol,
		i.size, i.symbol,
	)
	// Test mixed byte/bit symbol
	j, _ := NewIECUnit(2, MiB)
	k, _ := NewIECUnit(2, Mib)
	l, ok := j.Add(k).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", l))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.2f %s\n",
		j.size, j.symbol,
		k.size, k.symbol,
		l.size, l.symbol,
	)
	// Output:
	// 2 MiB + 2 MiB = 4 MiB
	// 2 Mib + 2 Mib = 4 Mib
	// 2 Mib + 2 MiB = 2.25 MiB
	// 2 MiB + 2 Mib = 2.25 MiB
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
				nsym UnitSymbol
				err  error
				exp  *uint
			)
			tur, _ := NewIECUnit(rand.Float64()*10, l)
			u := testIECUnitAdd{left: tul, right: tur}
			left := tul.ByteSize()
			right := tur.ByteSize()
			total := left + right
			nexp := uint(math.Round(math.Log2(total) / 10))
			if nexp > tul.exponent && nexp > tur.exponent {
				e := uint(nexp)
				exp = &e
			}
			if tul.symbol != tur.symbol {
				nsym, err = findGCDIECUnitSymbol(tul.symbol, tur.symbol, exp)
				assert.NoError(t, err)
			} else {
				nsym = tul.symbol
			}
			size := sizeInIECUnit(nsym, total)
			u.expected, _ = NewIECUnit(size, nsym)
			tests = append(tests, u)
		}
	}
	// Add a couple of bad entries for negative testing
	gu, _ := NewIECUnit(rand.Float64()*10, MiB)
	bu := &IECUnit{rand.Float64() * 10, UnitSymbol("FooBar"), 30}
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
	tests = append(tests, bul, bur)
	// Run through all the tests
	for _, tst := range tests {
		u, ok := tst.left.Add(tst.right).(*IECUnit)
		assert.Equal(t, true, ok)
		assert.Equal(t, tst.expected, u)
	}
}

func ExampleIECUnit_Subtract() {
	var (
		c, f, i *IECUnit
		ok      bool
	)
	// Test the same byte symbol
	a, _ := NewIECUnit(2, MiB)
	b, _ := NewIECUnit(2, MiB)
	c, ok = a.Subtract(b).(*IECUnit)
	if !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", c))
	}
	fmt.Printf(
		"%.f %s - %.f %s = %.f %s\n",
		a.size, a.symbol,
		b.size, b.symbol,
		c.size, c.symbol,
	)
	// Test the same bit symbol
	d, _ := NewIECUnit(2, Mib)
	e, _ := NewIECUnit(2, Mib)
	if f, ok = d.Subtract(e).(*IECUnit); !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", f))
	}
	fmt.Printf(
		"%.f %s - %.f %s = %.f %s\n",
		d.size, d.symbol,
		e.size, e.symbol,
		f.size, f.symbol,
	)
	// Test mixed bit/byte symbol
	g, _ := NewIECUnit(2, Mib)
	h, _ := NewIECUnit(2, MiB)
	if i, ok = g.Subtract(h).(*IECUnit); !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", i))
	}
	fmt.Printf(
		"%.f %s - %.f %s = %.2f %s\n",
		g.size, g.symbol,
		h.size, h.symbol,
		i.size, i.symbol,
	)
	// Output:
	// 2 MiB - 2 MiB = 0 Bit
	// 2 Mib - 2 Mib = 0 Bit
	// 2 Mib - 2 MiB = 1.75 MiB
}
