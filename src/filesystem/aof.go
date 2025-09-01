package filesystem

import (
	"bytes"
	"encoding/gob"
	"os"
)

type AppendOnlyStateChangeFile struct {
	file *os.File
}

func (aof *AppendOnlyStateChangeFile) Read() ([]StateChange, error) {
	return nil, nil
}

func (aof *AppendOnlyStateChangeFile) Append(newStateChange StateChange) error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(newStateChange); err != nil {
		return err
	}

	_, err := aof.file.Write(buf.Bytes())
	return err
}
