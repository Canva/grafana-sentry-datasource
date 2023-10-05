package sentry

import (
	"sync"
	"time"
)

const (
	defaultCacheTTL = 5 * time.Minute
)

type Cache struct {
	mu         sync.Mutex
	data       map[string]interface{}
	expiration map[string]time.Time
}

func NewCache() *Cache {
	return &Cache{
		data:       make(map[string]interface{}),
		expiration: make(map[string]time.Time),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	c.expiration[key] = time.Now().Add(ttl)
}

func (c *Cache) Get(key string) (interface{}, bool) {
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
