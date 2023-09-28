package sentry

import (
	"sync"
	"time"
)

const (
	defaultCacheTTL = 5 * time.Minute
)

type Cache[T comparable] struct {
	mu         sync.Mutex
	data       map[string]T
	expiration map[string]time.Time
}

func NewCache[T comparable]() *Cache {
	return &Cache{
		data:       make(map[string]T),
		expiration: make(map[string]time.Time),
	}
}

func (c *Cache) Set[T comparable](key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	c.expiration[key] = time.Now().Add(ttl)
}

func (c *Cache) Get[T comparable](key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, found := c.data[key]
	if !found {
		return nil, false
	}
	expTime, expFound := c.expiration[key]
	if !expFound || time.Now().After(expTime) {
		delete(c.data, key)
		delete(c.expiration, key)
		return nil, false
	}
	return value, true
}
