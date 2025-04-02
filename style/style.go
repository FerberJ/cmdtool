package style

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	Blue      = lipgloss.Color("63")
	LightGray = lipgloss.Color("241")
	Red       = lipgloss.Color("9")
	Green     = lipgloss.Color("10")
)

var (
	SpinnerStyle    = lipgloss.NewStyle().Foreground(Blue)
	HelpStyle       = lipgloss.NewStyle().Foreground(LightGray).Margin(1, 0)
	DotStyle        = HelpStyle.UnsetMargins()
	BackgroundStyle = DotStyle
	DurationStyle   = DotStyle
	AppStyle        = lipgloss.NewStyle()
	Error           = lipgloss.NewStyle().Foreground(Red)
	Title           = lipgloss.NewStyle().Bold(true).Background(Green).Width(80).Align(lipgloss.Center)
	Header          = lipgloss.NewStyle().Bold(true).Foreground(Green).Width(80).Align(lipgloss.Center)
	InputVar        = DotStyle.Width(40)
	InputResult     = lipgloss.NewStyle().Width(60)
)
