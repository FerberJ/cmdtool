package view

import "cmd/tool/style"

func spinnerView(m *model) string {
	var s string

	if len(m.addedVariables) > 0 {
		s += style.Title.Render("Inputs:")
		s += "\n\n"
		s += lipTable(m)
	}

	s += "\n\n"

	s += style.Title.Render("Commands:")
	s += "\n\n"
	if m.quitting {
		s += "Programm quitting"
	} else {
		s += m.spinner.View() + " Running commands"
	}

	s += "\n\n"

	for i := len(m.results) - consoleLength; i < len(m.results); i++ {
		s += m.results[i].String() + "\n"
	}

	s += "\n\n"

	s += style.Title.Render("Terminal out:")
	s += "\n\n"
	for i := len(m.terminalOut) - consoleLength; i < len(m.terminalOut); i++ {
		s += m.terminalOut[i].String() + "\n"
	}

	if !m.quitting {
		s += style.HelpStyle.Render("ctrl+c quit")
	}

	if m.quitting {
		s += "\n"
	}

	return s
}
