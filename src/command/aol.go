package command

import (
	"sync"
	"time"
)

type AppendOnlyCommandList struct {
	c_list  []Command
	top_ptr int
	lock    sync.RWMutex
}

func NewAppendOnlyCommandList() *AppendOnlyCommandList {
	return &AppendOnlyCommandList{
		c_list:  make([]Command, 0),
		top_ptr: -1,
		lock:    sync.RWMutex{},
	}
}

func (aol *AppendOnlyCommandList) Append(command Command) {
	aol.lock.Lock()
	defer aol.lock.Unlock()

}

func (aol *AppendOnlyCommandList) GetAll() []Command {
	return nil
}

func (aol *AppendOnlyCommandList) GetAllAfterTimestamp(timestamp time.Time) []Command {
	return nil
}

func (aol *AppendOnlyCommandList) Reindex() {

}
