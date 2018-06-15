package cache

import (
	"sync"
	"time"
)

// Cache is basically a thread-safe map[string]interface{}.
type Cache struct {
	mu    sync.RWMutex
	items map[string]interface{}
	inc   uint8
}

// New returns a new Cache intance.
func New() *Cache {
	return &Cache{
		mu:    sync.RWMutex{},
		items: make(map[string]interface{}),
		inc:   0,
	}
}

// NewWithCleanup returns a new Cache intance and a bool channel to stop the cleanup.
// Starts a cleanup go routine to delete expired items.
func NewWithCleanup(cleanupInterval time.Duration) (*Cache, chan bool) {
	stop := make(chan bool)
	c := New()

	go c.cleanup(cleanupInterval, stop)

	return c, stop
}

// Set adds an item to the cache, replacing an existing one.
func (c *Cache) Set(key string, value interface{}, options ...ItemOption) {
	c.mu.Lock()

	for _, option := range options {
		value = option(value)
	}

	c.items[key] = value

	c.inc++
	if c.inc == 0 {
		c.inc = 1
	}

	c.mu.Unlock()
}

// Get returns an item or nil and a bool indicating whether the key was found.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()

	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	c.mu.RUnlock()

	if itm, ok := item.(ExpiryItem); ok {
		if time.Now().UnixNano() > itm.Expiration {
			return nil, false
		}

		return itm.Value, true
	}

	return item, true
}

// Delete removes the key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()

	delete(c.items, key)

	c.inc++
	if c.inc == 0 {
		c.inc = 1
	}

	c.mu.Unlock()
}

// Clear deletes all items.
func (c *Cache) Clear() {
	c.mu.Lock()

	c.items = map[string]interface{}{}

	c.inc++
	if c.inc == 0 {
		c.inc = 1
	}

	c.mu.Unlock()
}

// Keys returns a list of all keys.
func (c *Cache) Keys() []string {
	c.mu.RLock()

	keys := make([]string, len(c.items))

	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}

	c.mu.RUnlock()

	return keys
}

// DeleteExpired deletes all ExpiryItems which are expired.
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()

	inc := false
	for key, item := range c.items {
		if itm, ok := item.(ExpiryItem); ok && now > itm.Expiration {
			delete(c.items, key)
			inc = true
		}
	}

	if inc {
		c.inc++
		if c.inc == 0 {
			c.inc = 1
		}
	}

	c.mu.Unlock()
}

// cleanup is a endless loop which deletes the expired items.
// The stop channel is used to break the loop.
func (c *Cache) cleanup(interval time.Duration, stop chan bool) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
