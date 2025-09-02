package cachemap

import (
	"sync"
	"time"
)

type CacheMap struct {
	cmap              map[string]Value
	lastTombstonedMap map[string]time.Time
	lock              sync.RWMutex
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		cmap:              make(map[string]Value),
		lastTombstonedMap: make(map[string]time.Time),
		lock:              sync.RWMutex{},
	}
}

// TODO: Handle tombstoning in a better way
func (c *CacheMap) Set(key string, value string, exp time.Time, timestamp time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.setKey(key, value, exp, timestamp)
}

func (c *CacheMap) SetWithoutLock(key string, value string, exp time.Time, timestamp time.Time) {
	c.setKey(key, value, exp, timestamp)
}

func (c *CacheMap) setKey(key string, value string, exp time.Time, timestamp time.Time) {
	if lastTombestonedAt, ok := c.lastTombstonedMap[key]; ok {
		if timestamp.Before(lastTombestonedAt) {
			return
		}
	}

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
	c.deleteKey(key, timestamp)
}

func (c *CacheMap) DeleteWithoutLock(key string, timestamp time.Time) {
	c.deleteKey(key, timestamp)
}

func (c *CacheMap) deleteKey(key string, timestamp time.Time) {
	if lastTombestonedAt, ok := c.lastTombstonedMap[key]; ok {
		if timestamp.Before(lastTombestonedAt) {
			return
		}
	}

	val, exists := c.cmap[key]

	if exists && val.LastUpdated.Before(timestamp) {
		c.lastTombstonedMap[key] = timestamp
		delete(c.cmap, key)
	}
}

func (c *CacheMap) Get(key string, timestamp time.Time) *Value {
	c.lock.Lock()
	defer c.lock.Unlock()

	if lastTombestonedAt, ok := c.lastTombstonedMap[key]; ok {
		if timestamp.Before(lastTombestonedAt) {
			return nil
		}
	}

	val, exists := c.cmap[key]

	if !exists {
		return nil
	}

	if val.ExpiresAt.After(timestamp) {
		return &val
	} else {
		delete(c.cmap, key)
	}

	return nil
}
