package view

import (
	"cmd/tool/style"
)

func inputView(m *model) string {
	var s string

	s += style.Title.Render("Inputs:")
	s += "\n\n"

	s += inputVars(m)
	s += "\n"
	s += "Variable input: \n"

	s += m.textInput.View()
	s += "\n"
	s += style.HelpStyle.Render("enter add variable • ctrl+c quit")
	return s
}
