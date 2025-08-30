package command

import "time"

type AppendOnlyCommandList struct {
	c_list  []Command
	top_ptr int64
}

func (aol *AppendOnlyCommandList) Append(command Command) {

}

func (aol *AppendOnlyCommandList) GetAll() []Command {
	return nil
}

func (aol *AppendOnlyCommandList) GetAllAfterTimestamp(timestamp time.Time) []Command {
	return nil
}
