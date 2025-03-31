package view

import "cmd/tool/style"

func sshConnectionView(m *model) string {
	var s string

	s += style.Title.Render("SSH-Connection:")
	s += "\n\n"
	s += m.sshConnection.String()
	s += "\n\n"

	return s
}
