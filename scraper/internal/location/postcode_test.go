package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPostcode(t *testing.T) {
	t.Run("should match on valid", func(t *testing.T) {
		ok, err := isPostcode("SW1A 2AA")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("should match on no space", func(t *testing.T) {
		ok, err := isPostcode("SW1A2AA")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("should not match invalid", func(t *testing.T) {
		ok, err := isPostcode("Downing Street, London, SW1A 2AA")
		require.NoError(t, err)
		assert.False(t, ok)
	})
}
