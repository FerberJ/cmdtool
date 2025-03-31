package view

import (
	"cmd/tool/style"
)

func inputView(m *model) string {
	var s string

	if len(m.addedVariables) > 0 {
		s += style.Title.Render("Inputs:")
		s += "\n\n"
		s += lipTable(m)
	}

	s += "\n\n"
	s += "Variable input: \n"

	s += m.textInput.View()
	s += "\n"
	s += style.HelpStyle.Render("enter add variable • ctrl+c quit")
	return s
}
