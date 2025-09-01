package filesystem

import (
	"time"
)

type StateChangeType string

const (
	StateChangeKeyUpdate StateChangeType = "KEY_UPDATE"
	StateChangeKeyDelete StateChangeType = "KEY_DELETE"
)

type StateChange interface {
	ChangeType() StateChangeType
}

type KeyUpdate struct {
	Timestamp time.Time
	Seed      string
	Key       string
	Val       string
	ExpiresAt time.Time
}

func (KeyUpdate) ChangeType() StateChangeType {
	return StateChangeKeyUpdate
}

type KeyDelete struct {
	Timestamp time.Time
	Seed      string
	Key       string
}

func (*KeyDelete) ChangeType() StateChangeType {
	return StateChangeKeyUpdate
}
