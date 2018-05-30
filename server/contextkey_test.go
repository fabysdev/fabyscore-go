package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	var testKey = &ContextKey{"test"}

	assert.Equal(t, "fabyscore-go_test", testKey.String())
	assert.Equal(t, "test", testKey.name)
}
