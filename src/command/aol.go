package command

import (
	"slices"
	"strconv"
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
	aol.lock.Lock()
	defer aol.lock.Unlock()

	l := make([]Command, len(aol.c_list))
	copy(l, aol.c_list)
	return l
}

func (aol *AppendOnlyCommandList) GetAllAfterTimestamp(timestamp time.Time) []Command {
	aol.lock.Lock()
	defer aol.lock.Unlock()

	l := make([]Command, 0)

	for _, command := range aol.c_list {
		commandTimestamp := command.Timestamp

		if commandTimestamp.After(timestamp) {
			l = append(l, command)
		}
	}

	return l
}

func (aol *AppendOnlyCommandList) Reindex() {
	aol.lock.Lock()
	defer aol.lock.Unlock()

	nl := make([]Command, 0) //TODO: Pick better capacity size
	keyToLatestCommandMap := make(map[string]Command)

	slices.SortFunc[[]Command](aol.c_list, func(c1, c2 Command) int {
		if c1.Timestamp.After(c2.Timestamp) {
			return 1
		} else if c1.Timestamp.Before(c2.Timestamp) {
			return -1
		} else {
			if c1.Seed > c1.Seed {
				return 1
			} else {
				return -1
			}
		}
	})

	currentTimestamp := time.Now()

	for _, command := range aol.c_list {
		if command.Name == CommandNameDelete {
			delete(keyToLatestCommandMap, command.Key)
		} else if command.Name == CommandNameSet {
			epochTimestampMillis, _ := strconv.ParseInt(command.Params["expiresAt"], 10, 64) //TODO: Handle error
			expirationTimestamp := time.UnixMilli(epochTimestampMillis)

			if currentTimestamp.Before(expirationTimestamp) {
				keyToLatestCommandMap[command.Key] = command
			}
		}
	}

	for _, command := range keyToLatestCommandMap {
		nl = append(nl, command)
	}

	aol.c_list = nl
}
