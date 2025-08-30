package command

type CommandName int

const (
	Set CommandName = iota
	Delete
)

type Command struct {
	Name      CommandName
	Key       string
	Timestamp int64
	Params    map[string]string
}
