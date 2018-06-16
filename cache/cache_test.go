package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNotFound(t *testing.T) {
	c := New()
	c.Set("a", "valueA")
	c.Set("b", "valueB")

	item, found := c.Get("test")
	assert.Nil(t, item)
	assert.False(t, found)
}

func TestGet(t *testing.T) {
	c := New()
	c.Set("test", "value")
	c.Set("a", "valueA")
	c.Set("b", "valueB")

	item, found := c.Get("test")
	assert.True(t, found)
	assert.Equal(t, "value", item)
}

func TestExpiry(t *testing.T) {
	c := New()
	c.Set("test", "value", Expire(50*time.Millisecond))
	c.Set("a", "valueA")
	c.Set("b", "valueB", Expire(2*time.Second))

	time.Sleep(30 * time.Millisecond)

	item, found := c.Get("test")
	assert.True(t, found)
	assert.Equal(t, "value", item)

	time.Sleep(30 * time.Millisecond)

	item, found = c.Get("test")
	assert.Nil(t, item)
	assert.False(t, found)

	item, found = c.Get("a")
	assert.True(t, found)
	assert.Equal(t, "valueA", item)

	item, found = c.Get("b")
	assert.True(t, found)
	assert.Equal(t, "valueB", item)
}

func TestDelete(t *testing.T) {
	c := New()
	c.Set("test", "value")
	c.Set("a", "valueA")
	c.Set("b", "valueB", Expire(2*time.Second))

	item, found := c.Get("test")
	assert.True(t, found)
	assert.Equal(t, "value", item)

	c.Delete("test")

	item, found = c.Get("test")
	assert.Nil(t, item)
	assert.False(t, found)

	item, found = c.Get("b")
	assert.True(t, found)
	assert.Equal(t, "valueB", item)
}

func TestClear(t *testing.T) {
	c := New()
	c.Set("test", "value")
	c.Set("a", "valueA")
	c.Set("b", "valueB", Expire(2*time.Second))

	item, found := c.Get("test")
	assert.True(t, found)
	assert.Equal(t, "value", item)

	c.Clear()

	item, found = c.Get("test")
	assert.Nil(t, item)
	assert.False(t, found)

	item, found = c.Get("b")
	assert.Nil(t, item)
	assert.False(t, found)

	item, found = c.Get("a")
	assert.Nil(t, item)
	assert.False(t, found)
}

func TestKeys(t *testing.T) {
	c := New()
	c.Set("test", "value")
	c.Set("a", "valueA")
	c.Set("b", "valueB", Expire(2*time.Second))

	keys := c.Keys()

	assert.Len(t, keys, 3)
	assert.ElementsMatch(t, []string{"test", "a", "b"}, keys)
}

func TestDeleteExpired(t *testing.T) {
	c := New()
	c.Set("test", "value")
	c.Set("a", "valueA")
	c.Set("b", "valueB", Expire(10*time.Millisecond))
	assert.Len(t, c.Keys(), 3)

	c.DeleteExpired()
	assert.Len(t, c.Keys(), 3)

	time.Sleep(50 * time.Millisecond)

	c.DeleteExpired()
	assert.Len(t, c.Keys(), 2)
}

func TestNewWithCleanup(t *testing.T) {
	c, stop := NewWithCleanup(50 * time.Millisecond)

	c.Set("test", "value", Expire(10*time.Millisecond))
	c.Set("a", "valueA")

	assert.Len(t, c.Keys(), 2)

	time.Sleep(60 * time.Millisecond)

	assert.Len(t, c.Keys(), 1)

	stop <- true

	c.Set("test", "value", Expire(10*time.Millisecond))
	assert.Len(t, c.Keys(), 2)

	time.Sleep(60 * time.Millisecond)
	assert.Len(t, c.Keys(), 2)
}
