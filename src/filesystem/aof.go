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

func NewAppendOnlyStateChangeFile() *AppendOnlyStateChangeFile {
	//os.Open
	return nil
}

func (aof *AppendOnlyStateChangeFile) Close() error {
	return aof.Close()
}

func (aof *AppendOnlyStateChangeFile) Read() ([]StateChange, error) {
	changes := make([]StateChange, 0)

	for {
		buf1 := make([]byte, binary.MaxVarintLen32)
		if numRead, _ := aof.file.Read(buf1); numRead == 0 {
			break
		}

		sz, _ := binary.Uvarint(buf1) //TODO: Handle errors

		buf2 := make([]byte, sz)

		if numRead, _ := aof.file.Read(buf2); numRead == 0 {
			break
		}

		var change StateChange

		dec := gob.NewDecoder(bytes.NewBuffer(buf2))
		if err := dec.Decode(change); err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	return changes, nil
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
