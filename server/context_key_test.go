package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	var testKey = &ContextKey{"test"}

	assert.Equal(t, "fabyscore-go_test", testKey.String())
	assert.Equal(t, "test", testKey.Name)
}

func TestContextKeyUnique(t *testing.T) {
	key1 := &ContextKey{"key1"}
	key2 := &ContextKey{"key2"}

	assert.NotEqual(t, key1, key2)
}
