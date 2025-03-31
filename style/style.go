package style

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle = lipgloss.NewStyle().SetString("Error: ").Foreground(lipgloss.Color("9")).Bold(true)
)
