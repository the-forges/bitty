package bitty

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mathTableTest struct {
	l              string
	r              string
	expected       string
	expectedSymbol UnitSymbol
	formatter      string
	expectedErr    error
	msg            string
}

type mathSadTableTest struct {
	l   Unit
	r   Unit
	msg string
}

func TestAddUnits(t *testing.T) {
	tt := []mathTableTest{
		{
			l:              "100 GB",
			r:              "50 GB",
			expected:       "150 GB",
			expectedSymbol: GB,
			formatter:      "%0.f %s",
			expectedErr:    nil,
			msg:            "can add units of the same standard",
		},
		{
			l:              "1 GB",
			r:              "1 GiB",
			expected:       "2.073742 GB",
			expectedSymbol: GB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add units of the different standards",
		},
		{
			l:              "1 GiB",
			r:              "1 GB",
			expected:       "1.931323 GiB",
			expectedSymbol: GiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add units of the different standards, reversed",
		},
		{
			l:              "100000 TiB", // 10,000
			r:              "500000 TB",  // 50,000
			expected:       "554747.350886 TiB",
			expectedSymbol: TiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add large numbers",
		},
	}

	for _, test := range tt {
		left, err := Parse(test.l)

		right, err := Parse(test.r)

		actualUnit, err := AddUnits(left, right)
		assert.Equal(t, test.expectedErr, err)

		actual := fmt.Sprintf(test.formatter, actualUnit.SizeInUnit(test.expectedSymbol), actualUnit.Symbol())
		assert.Equal(t, test.expected, actual)
	}

}

func TestAddUnitsSadPath(t *testing.T) {
	tt := []mathSadTableTest{
		{
			l: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			r: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			msg: "returns an error if both units are invalid",
		},
		{
			l: &IECUnit{
				size:   100,
				symbol: GiB,
			},
			r: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			msg: "returns an error if r is invalid",
		},
		{
			l: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			r: &IECUnit{
				size:   100,
				symbol: GiB,
			},
			msg: "returns an error if l is invalid",
		},
	}

	for _, test := range tt {
		_, err := AddUnits(test.l, test.r)
		assert.Error(t, err, test.msg)
	}

}

func TestSubtractUnits(t *testing.T) {
	tt := []mathTableTest{
		{
			l:              "100 GB",
			r:              "50 GB",
			expected:       "50 GB",
			expectedSymbol: GB,
			formatter:      "%0.f %s",
			expectedErr:    nil,
			msg:            "can subtract units of the same standard",
		},
		{
			l:              "1 GB",
			r:              "1 GiB",
			expected:       "-0.073742 GB",
			expectedSymbol: GB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can subtract units of the different standards",
		},
		{
			l:              "1 GiB",
			r:              "1 GB",
			expected:       "0.068677 GiB",
			expectedSymbol: GiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can subtract units of the different standards, reversed",
		},
		{
			l:              "100000 TiB", // 10,000
			r:              "500000 TB",  // 50,000
			expected:       "-354747.350886 TiB",
			expectedSymbol: TiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can subtract large numbers",
		},
	}

	for _, test := range tt {
		left, err := Parse(test.l)
		assert.NoError(t, err)

		right, err := Parse(test.r)
		assert.NoError(t, err)

		actualUnit, err := SubtractUnits(left, right)
		assert.Equal(t, test.expectedErr, err)

		actual := fmt.Sprintf(test.formatter, actualUnit.Size(), actualUnit.Symbol())
		assert.Equal(t, test.expected, actual)
	}

}

func TestSubtractUnitsSadPath(t *testing.T) {
	tt := []mathSadTableTest{
		{
			l: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			r: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			msg: "returns an error if both units are invalid",
		},
		{
			l: &IECUnit{
				size:   100,
				symbol: GiB,
			},
			r: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			msg: "returns an error if r is invalid",
		},
		{
			l: &IECUnit{
				size:   100,
				symbol: UnitSymbol("giib"),
			},
			r: &IECUnit{
				size:   100,
				symbol: GiB,
			},
			msg: "returns an error if l is invalid",
		},
	}

	for _, test := range tt {
		_, err := SubtractUnits(test.l, test.r)
		assert.Error(t, err, test.msg)
	}

}
