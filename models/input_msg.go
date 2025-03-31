package models

import "sync"

type InputMsg struct {
	Input *string
	Msg   string
	Wg    *sync.WaitGroup
}

func (i InputMsg) String() string {
	return i.Msg
}
