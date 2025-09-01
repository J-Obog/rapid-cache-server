package filesystem

import (
	"bytes"
	"encoding/binary"
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

	buf2 := make([]byte, binary.MaxVarintLen32)
	binary.BigEndian.PutUint32(buf2, uint32(buf.Len()))

	_, err := aof.file.Write(append(buf2, buf.Bytes()...))
	return err
}
