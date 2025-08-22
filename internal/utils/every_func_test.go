package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEveryFunc(t *testing.T) {
	t.Run("Test every func", func(t *testing.T) {
		input := []bool{false, false, false, true}

		res := Every(input, func(b bool) bool {
			return b == true
		})

		assert.Equal(t, false, res)
	})
}
