package cachemap

import "time"

type Value struct {
	Value       string
	ExpiresAt   time.Time
	LastUpdated time.Time
}

type CacheMap struct {
	cmap map[string]Value
}

func (c *CacheMap) Set(key string, value string, exp time.Time, timestamp time.Time) {

}

func (c *CacheMap) Delete(key string, timestamp time.Time) {

}
