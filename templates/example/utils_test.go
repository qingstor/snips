package example

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert.True(t, strings.HasPrefix(Utils(), "Hi there"))
}
