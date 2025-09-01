package filesystem

type AppendOnlyCommandFile struct {
}

func (aof *AppendOnlyCommandFile) Read() error {
	return nil
}

func (aof *AppendOnlyCommandFile) Append(newCommand *StateChange) error {
	return nil
}
