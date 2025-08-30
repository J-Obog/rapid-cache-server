package command

import "time"

type CommandName int

const (
	CommandNameSet CommandName = iota
	CommandNameDelete
)

type Command struct {
	Name      CommandName
	Key       string
	Timestamp time.Time
	Params    map[string]string
}
