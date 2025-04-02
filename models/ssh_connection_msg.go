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
	var str string
	if s.Pending {
		str += style.InputVar.Render("⏳ " + s.User)
	} else if s.Success {
		str += style.InputVar.Render("✅ " + s.User)
	} else {
		str += style.InputVar.Render("❌ " + s.User + " ")
		str += style.Error.Render(s.Error)
	}
	str += style.InputResult.Render(s.Address)
	return str
}
