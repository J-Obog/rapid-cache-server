package command

import "time"

type CommandName int
type CommandParamKey int

const (
	CommandNameSet CommandName = iota
	CommandNameDelete
)

const (
	CommandParamKeyValue CommandParamKey = iota
	CommandParamKeyExpiresAt
)

type Command struct {
	Name      CommandName
	Key       string
	Timestamp time.Time
	Seed      string
	Params    map[CommandParamKey]string
}
