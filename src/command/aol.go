package command

import (
	"sync"
	"time"
)

type AppendOnlyCommandList struct {
	c_list []Command
	lock   sync.RWMutex
}

func NewAppendOnlyCommandList() *AppendOnlyCommandList {
	return &AppendOnlyCommandList{
		c_list: make([]Command, 0),
		lock:   sync.RWMutex{},
	}
}

func (aol *AppendOnlyCommandList) Append(newCommand Command) {
	aol.lock.Lock()
	defer aol.lock.Unlock()
	aol.c_list = append(aol.c_list, newCommand)
}

func (aol *AppendOnlyCommandList) GetAll() []Command {
	return nil
}

func (aol *AppendOnlyCommandList) GetAllAfterTimestamp(timestamp time.Time) []Command {
	return nil
}

func (aol *AppendOnlyCommandList) Reindex() {

}
