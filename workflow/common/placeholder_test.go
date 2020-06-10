package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNextPlaceholder verifies dynamically-generated placeholder strings.
func TestNextPlaceholder(t *testing.T) {
	pg := NewPlaceholderGenerator()
	assert.Equal(t, pg.NextPlaceholder(), fmt.Sprintf("%s0", placeholderPrefix))
	assert.Equal(t, pg.NextPlaceholder(), fmt.Sprintf("%s1", placeholderPrefix))
	assert.Equal(t, pg.NextPlaceholder(), fmt.Sprintf("%s2", placeholderPrefix))

	assert.True(t, pg.IsPlaceholder(fmt.Sprintf("%s0", placeholderPrefix)))
	assert.True(t, pg.IsPlaceholder(fmt.Sprintf("%s1", placeholderPrefix)))
	assert.True(t, pg.IsPlaceholder(fmt.Sprintf("%s2", placeholderPrefix)))

	assert.False(t, pg.IsPlaceholder(fmt.Sprintf("%s3", placeholderPrefix)))
	assert.False(t, pg.IsPlaceholder(fmt.Sprintf("%saa", placeholderPrefix)))
	assert.False(t, pg.IsPlaceholder(fmt.Sprintf("%s2", placeholderPrefix+"aa")))
}
