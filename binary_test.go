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

func ExampleIECSymbolByteSize() {
	// legit UnitSymbol
	a := IECSymbolByteSize(MiB, 10)
	fmt.Printf("%.f\n", a)
	// hack in a bad UnitSymbol to test default path
	iecUnitExponentMap["foo"] = 9
	b := IECSymbolByteSize("foo", 10)
	fmt.Printf("%.f\n", b)
	// reset hack
	delete(iecUnitExponentMap, "foo")

	// Output:
	// 10485760
	// 0
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

func TestIECByteSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnit, 0, len(iecUnitExponentMap))
	// Setup test cases based out of what is in IECUnitExponentMap
	for k, v := range iecUnitExponentMap {
		tu, _ := NewIECUnit(rand.Float64(), k)
		if tu == nil {
			break
		}
		u := testIECUnit{Unit: *tu}
		exp := float64(v * 10)
		bytes := float64(math.Exp2(exp) * u.Unit.size)
		switch k {
		case Bit:
			u.Expected = u.Unit.size * 8
		case Byte:
			u.Expected = u.Unit.size
		case Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib:
			u.Expected = bytes * 0.125
		case KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
			u.Expected = bytes
		default:
			u.Expected = float64(0)
		}
		tests = append(tests, u)
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

func ExampleBitSize() {
	a, _ := NewIECUnit(10.0, MiB)
	b, _ := NewIECUnit(10.0, Mib)
	fmt.Printf("%.f\n", a.BitSize())
	fmt.Printf("%.f\n", b.BitSize())
	// Output:
	// 83886080
	// 10485760
}

func ExampleIECSymbolBitSize() {
	// legit UnitSymbol
	a := IECSymbolBitSize(MiB, IECSymbolByteSize(MiB, 10))
	fmt.Printf("%.f\n", a)
	// hack in a bad UnitSymbol to test default path
	iecUnitExponentMap["foo"] = 9
	b := IECSymbolBitSize("foo", IECSymbolByteSize("foo", 10))
	fmt.Printf("%.f\n", b)
	// reset hack
	delete(iecUnitExponentMap, "foo")

	// Output:
	// 83886080
	// 0
}

func TestIECBitSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECUnit, 0, len(iecUnitExponentMap))
	// Setup test cases based out of what is in IECUnitExponentMap
	for k := range iecUnitExponentMap {
		tu, _ := NewIECUnit(rand.Float64()*10, k)
		if tu == nil {
			break
		}
		u := testIECUnit{Unit: *tu}
		bytes := tu.ByteSize()
		switch k {
		case Bit:
			u.Expected = u.Unit.size
		case Byte,
			Kib, Mib, Gib, Tib, Pib, Eib, Zib, Yib,
			KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB:
			u.Expected = float64(bytes * 8)
		default:
			u.Expected = float64(0)
		}
		tests = append(tests, u)
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

func TestIECSizeInUnit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testIECSizeInUnit, 0, len(iecUnitExponentMap))
	// Setup test cases based out of what is in IECUnitExponentMap
	for k := range iecUnitExponentMap {
		tu, _ := NewIECUnit(rand.Float64()*10, k)
		if tu == nil {
			break
		}
		for l, v := range iecUnitExponentMap {
			u := testIECSizeInUnit{unit: *tu, to: l}
			var (
				diffExp = float64(u.unit.exponent - v)
				left    = IECSymbolByteSize(u.unit.symbol, u.unit.size)
				right   = IECSymbolByteSize(l, u.unit.size)
			)
			if diffExp > 0 {
				u.expected = right * diffExp
			} else {
				u.expected = (left / right) * u.unit.size
			}
			tests = append(tests, u)
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
		u := findNearestIECUnitSymbols(i)
		assert.Equal(t, v, u)
	}
}

func ExampleIECUnit_Add() {
	var (
		c, f, i *IECUnit
		ok      bool
	)
	// Test the same byte symbol
	a, _ := NewIECUnit(2, MiB)
	b, _ := NewIECUnit(2, MiB)
	if c, ok = a.Add(b).(*IECUnit); !ok {
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
	if f, ok = d.Add(e).(*IECUnit); !ok {
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
	if i, ok = g.Add(h).(*IECUnit); !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", i))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.2f %s\n",
		g.size, g.symbol,
		h.size, h.symbol,
		i.size, i.symbol,
	)
	// Output:
	// 2 MiB + 2 MiB = 4 MiB
	// 2 Mib + 2 Mib = 4 Mib
	// 2 Mib + 2 MiB = 2.25 MiB
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
			tur, _ := NewIECUnit(rand.Float64()*10, l)
			u := testIECUnitAdd{left: tul, right: tur}
			left := tul.ByteSize()
			right := tur.ByteSize()
			total := left + right
			nexp := int(math.Round(math.Log2(total) / 10))
			nsym := findLargestIECUnitSymbol(tul.symbol, tur.symbol, nexp)
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
