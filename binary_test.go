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

func ExampleIECSizeInUnit() {
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

func ExampleIECUnit_Add() {
	var (
		tmp Unit
		err error
		c   *IECUnit
		ok  bool
	)
	a, _ := NewIECUnit(2, MiB)
	b, _ := NewIECUnit(2, MiB)
	if tmp, err = a.Add(b); err != nil {
		panic(err)
	}
	if c, ok = tmp.(*IECUnit); !ok {
		panic(fmt.Errorf("Unit not *IECUnit: %v", c))
	}
	fmt.Printf(
		"%.f %s + %.f %s = %.f %s\n",
		a.size, a.symbol,
		b.size, b.symbol,
		c.size, c.symbol,
	)
	// Output:
	// 2 MiB + 2 MiB = 4 MiB
}
