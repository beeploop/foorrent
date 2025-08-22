package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeOverlap(t *testing.T) {
	t.Run("Test if two ranges is overlapping", func(t *testing.T) {
		tests := []struct {
			range1   Range
			range2   Range
			expected bool
		}{
			{
				range1:   Range{1, 2},
				range2:   Range{2, 3},
				expected: false,
			},
			{
				range1:   Range{4, 6},
				range2:   Range{2, 4},
				expected: false,
			},
			{
				range1:   Range{1, 5},
				range2:   Range{2, 3},
				expected: true,
			},
			{
				range1:   Range{2, 3},
				range2:   Range{1, 5},
				expected: true,
			},
			{
				range1:   Range{2, 3},
				range2:   Range{4, 5},
				expected: false,
			},
			{
				range1:   Range{0, 1000},
				range2:   Range{1000, 2000},
				expected: false,
			},
			{
				range1:   Range{0, 1000},
				range2:   Range{900, 2000},
				expected: true,
			},
		}

		for _, test := range tests {
			result := IsOverlapping(test.range1, test.range2)

			assert.Equal(t, test.expected, result, "expected %s and %s to have result of %s", test.range1, test.range2, test.expected)
		}
	})
}
