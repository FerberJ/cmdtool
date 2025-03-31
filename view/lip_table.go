package view

import (
	"cmd/tool/style"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func lipTable(m *model) string {
	var s string

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(style.LightGray)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return style.HeaderStyle
			case row%2 == 0:
				return style.EvenRowStyle
			default:
				return style.OddRowStyle
			}
		}).
		Headers("VAR", "INPUT").
		Rows(m.addedVariables...)

	s += t.Render()
	return s
}
