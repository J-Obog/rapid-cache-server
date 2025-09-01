package filesystem

import (
	"time"
)

type WriteOperation interface {
	OperationType() string
}

type SetKeyOperation struct {
	Timestamp time.Time
	Seed      string
	Key       string
	Val       string
	ExpiresAt time.Time
}

func (SetKeyOperation) OperationType() string {
	return "KEY_SET"
}

type DeleteKeyOperation struct {
	Timestamp time.Time
	Seed      string
	Key       string
}

func (DeleteKeyOperation) OperationType() string {
	return "KEY_DELETE"
}
