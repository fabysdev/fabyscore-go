package cache

import (
	"math"
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

func TestInc(t *testing.T) {
	c := New()

	c.Set("test", "value")
	assert.Equal(t, uint8(1), c.inc)

	c.Set("a", "valueA")
	assert.Equal(t, uint8(2), c.inc)

	c.Delete("test")
	assert.Equal(t, uint8(3), c.inc)

	c.Clear()
	assert.Equal(t, uint8(4), c.inc)

	c.Set("test", "value", Expire(10*time.Millisecond))
	assert.Equal(t, uint8(5), c.inc)

	time.Sleep(50 * time.Millisecond)

	c.DeleteExpired()
	assert.Equal(t, uint8(6), c.inc)
}

func TestIncNotZero(t *testing.T) {
	c := New()

	c.inc = math.MaxUint8
	c.Set("test", "value")
	assert.Equal(t, uint8(1), c.inc)

	c.Set("a", "valueA")
	assert.Equal(t, uint8(2), c.inc)

	c.inc = math.MaxUint8
	c.Delete("test")
	assert.Equal(t, uint8(1), c.inc)

	c.inc = math.MaxUint8
	c.Clear()
	assert.Equal(t, uint8(1), c.inc)

	c.Set("test", "value", Expire(10*time.Millisecond))
	assert.Equal(t, uint8(2), c.inc)

	time.Sleep(50 * time.Millisecond)

	c.inc = math.MaxUint8
	c.DeleteExpired()
	assert.Equal(t, uint8(1), c.inc)
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
