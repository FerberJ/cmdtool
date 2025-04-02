package view

import (
	"cmd/tool/style"
)

func inputVars(m *model) string {
	var s string

	for _, in := range m.addedVariables {
		s += style.InputVar.Render(in[0] + ": ")
		s += style.InputResult.Render(in[1])
		s += "\n"
	}

	return s
}
