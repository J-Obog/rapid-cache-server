package command

import "time"

type CommandName int

const (
	Set CommandName = iota
	Delete
)

type Command struct {
	Name      CommandName
	Key       string
	Timestamp time.Time
	Params    map[string]string
}
