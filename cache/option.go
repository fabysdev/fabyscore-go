package cache

import "time"

// ItemOption is a function type to modify the value.
type ItemOption func(value interface{}) interface{}

// ExpiryItem defines an item with expiration.
type ExpiryItem struct {
	Value      interface{}
	Expiration int64
}

// Expire returns an ItemOption for creating an ExpiryItem.
func Expire(d time.Duration) ItemOption {
	return func(value interface{}) interface{} {
		return ExpiryItem{
			Value:      value,
			Expiration: time.Now().Add(d).UnixNano(),
		}
	}
}
