package filesystem

import "github.com/J-Obog/rapid-cache-server/src/command"

type AppendOnlyCommandFile struct {
}

func (aof *AppendOnlyCommandFile) Read() error {
	return nil
}

func (aof *AppendOnlyCommandFile) Append(newCommand *command.Command) error {
	return nil
}
