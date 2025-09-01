package filesystem

type AppendOnlyStateChangeFile struct {
}

func (aof *AppendOnlyStateChangeFile) Read() ([]StateChange, error) {
	return nil, nil
}

func (aof *AppendOnlyStateChangeFile) Append(newCommand StateChange) error {
	return nil
}
