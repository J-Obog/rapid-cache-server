package filesystem

type AppendOnlyStateChangeFile struct {
}

func (aof *AppendOnlyStateChangeFile) Read() error {
	return nil
}

func (aof *AppendOnlyStateChangeFile) Append(newCommand StateChange) error {
	return nil
}
