package bitty

import (
	"fmt"
	"testing"

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

type testNewUnit struct {
	std      UnitStandard
	size     float64
	sym      UnitSymbol
	expected Unit
	err      error
}

func TestNewUnit(t *testing.T) {
	tt := []testNewUnit{
		{
			IEC,
			1,
			MiB,
			&IECUnit{1, MiB, 2},
			nil,
		},
		{
			IEC,
			1,
			GiB,
			&IECUnit{1, GiB, 3},
			nil,
		},
		{
			SI,
			1,
			MB,
			&SIUnit{1, MB, 6},
			nil,
		},
		{
			UnitStandard(50),
			1,
			MB,
			nil,
			fmt.Errorf("%v is currently not a supported standard", 50),
		},
	}
	for _, u := range tt {
		nu, err := NewUnit(u.std, u.size, u.sym)
		assert.Equal(t, u.err, err)
		assert.Equal(t, u.expected, nu)
	}
}

type parseExampleData struct {
	input    string
	expected Unit
	err      error
}

func TestParse(t *testing.T) {
	a, _ := NewIECUnit(1, MiB)
	b, _ := NewSIUnit(1, MB)
	c, _ := NewIECUnit(1.64, Kib)
	d, _ := NewSIUnit(-1, Mb)
	e, _ := NewSIUnit(1, MB)
	erra := NewErrUnitCouldNotBeParsed("one MiB")
	errb := NewErrUnitStandardNotSupported(UnitStandard(2))
	tt := []parseExampleData{
		{"1 MiB", a, nil},
		{"1 MB", b, nil},
		{"1.64 Kib", c, nil},
		{"-1 Mb", d, nil},
		{"1MB", e, nil},
		{"one MiB", nil, erra},
		{"1 Bab", nil, errb},
	}
	for _, d := range tt {
		u, err := Parse(d.input)
		if err != nil || d.err != nil {
			assert.Error(t, d.err, err)
		}
		assert.Equal(t, d.expected, u)
	}
}
