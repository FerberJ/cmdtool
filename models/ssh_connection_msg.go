package models

import (
	"cmd/tool/style"
)

type SshConnectionMsg struct {
	User    string
	Address string
	Pending bool
	Success bool
	Error   string
}

func (s SshConnectionMsg) String() string {
	str := s.User + "\t" + s.Address
	if s.Pending {
		str += "\t⏳"
	} else if s.Success {
		str += "\t✅"
	} else {
		str += "\t❌"
		str += "\n"
		str += style.Error.Render(s.Error)
	}
	return str
}
