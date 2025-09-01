package filesystem

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"os"
)

type WriteOperationAOF struct {
	file *os.File
}

func (aof *WriteOperationAOF) Open(filePath string) error {
	filePtr, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) //TODO: set with correct permission

	if err != nil {
		return err
	}

	aof.file = filePtr
	return nil
}

func (aof *WriteOperationAOF) Close() error {
	return aof.file.Close()
}

func (aof *WriteOperationAOF) Read() ([]WriteOperation, error) {
	changes := make([]WriteOperation, 0)

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

		var change WriteOperation

		dec := gob.NewDecoder(bytes.NewBuffer(buf2))
		if err := dec.Decode(change); err != nil {
			return nil, err
		}

		changes = append(changes, change)
	}

	return changes, nil
}

func (aof *WriteOperationAOF) Append(newStateChange WriteOperation) error {
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
