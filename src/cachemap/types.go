package cachemap

import "time"

type Value struct {
	Val         string
	ExpiresAt   time.Time
	LastUpdated time.Time
}
