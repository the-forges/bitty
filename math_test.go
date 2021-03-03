package bitty

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUnits(t *testing.T) {
	tt := []struct {
		l              string
		r              string
		expected       string
		expectedSymbol UnitSymbol
		formatter      string
		expectedErr    error
		msg            string
	}{
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
		assert.NoError(t, err)

		right, err := Parse(test.r)
		assert.NoError(t, err)

		actualUnit, err := AddUnits(left, right)
		assert.Equal(t, test.expectedErr, err)

		actual := fmt.Sprintf(test.formatter, actualUnit.SizeInUnit(test.expectedSymbol), actualUnit.Symbol())
		assert.Equal(t, test.expected, actual)
	}

}

func TestSubtractUnits(t *testing.T) {
	tt := []struct {
		l              string
		r              string
		expected       string
		expectedSymbol UnitSymbol
		formatter      string
		expectedErr    error
		msg            string
	}{
		{
			l:              "100 GB",
			r:              "50 GB",
			expected:       "50 GB",
			expectedSymbol: GB,
			formatter:      "%0.f %s",
			expectedErr:    nil,
			msg:            "can add units of the same standard",
		},
		{
			l:              "1 GB",
			r:              "1 GiB",
			expected:       "-0.073742 GB",
			expectedSymbol: GB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add units of the different standards",
		},
		{
			l:              "1 GiB",
			r:              "1 GB",
			expected:       "0.068677 GiB",
			expectedSymbol: GiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add units of the different standards, reversed",
		},
		{
			l:              "100000 TiB", // 10,000
			r:              "500000 TB",  // 50,000
			expected:       "-354747.350886 TiB",
			expectedSymbol: TiB,
			formatter:      "%f %s",
			expectedErr:    nil,
			msg:            "can add large numbers",
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
