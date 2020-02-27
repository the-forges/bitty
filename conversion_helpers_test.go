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
		testNewUnit{
			IEC,
			1,
			MiB,
			&IECUnit{1, MiB, 2},
			nil,
		},
		testNewUnit{
			IEC,
			1,
			GiB,
			&IECUnit{1, GiB, 3},
			nil,
		},
		testNewUnit{
			SI,
			1,
			MB,
			nil,
			fmt.Errorf("%s is currently not a supported standard", string(SI)),
		},
		testNewUnit{
			UnitStandard(50),
			1,
			MB,
			nil,
			fmt.Errorf("%s is currently not a supported standard", string(50)),
		},
	}
	for _, u := range tt {
		nu, err := NewUnit(u.std, u.size, u.sym)
		assert.Equal(t, u.err, err)
		assert.Equal(t, u.expected, nu)
	}
}
