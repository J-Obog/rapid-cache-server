package cachemap

import (
	"sync"
	"time"
)

type CacheMap struct {
	cmap map[string]Value
	lock sync.RWMutex
}

// TODO: Handle tombstoning in a better way
func (c *CacheMap) Set(key string, value string, exp time.Time, timestamp time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	val, exists := c.cmap[key]

	if !exists || val.LastUpdated.Before(timestamp) {
		c.cmap[key] = Value{
			Val:         value,
			ExpiresAt:   exp,
			LastUpdated: timestamp,
		}
	}
}

func (c *CacheMap) Delete(key string, timestamp time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	val, exists := c.cmap[key]

	if exists && val.LastUpdated.Before(timestamp) {
		delete(c.cmap, key)
	}
}
