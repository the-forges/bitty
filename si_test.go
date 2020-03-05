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
	a, _ := NewSIUnit(10.0, MB)
	fmt.Printf("%.f\n", a.ByteSize())
	// Output:
	// 10000000
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
