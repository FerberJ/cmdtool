package style

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	Blue      = lipgloss.Color("63")
	LightGray = lipgloss.Color("241")
	Red       = lipgloss.Color("9")
	Green     = lipgloss.Color("10")
	White     = lipgloss.Color("231")
)

var (
	SpinnerStyle    = lipgloss.NewStyle().Foreground(Blue)
	HelpStyle       = lipgloss.NewStyle().Foreground(LightGray).Margin(1, 0)
	DotStyle        = HelpStyle.UnsetMargins()
	BackgroundStyle = DotStyle
	DurationStyle   = DotStyle
	AppStyle        = lipgloss.NewStyle()
	Error           = lipgloss.NewStyle().Foreground(Red)
	Title           = lipgloss.NewStyle().Bold(true).Background(Green).Width(60).Align(lipgloss.Center)
)

// Table
var (
	HeaderStyle  = lipgloss.NewStyle().Foreground(LightGray).Bold(true).Align(lipgloss.Center)
	CellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(26)
	OddRowStyle  = CellStyle.Foreground(LightGray)
	EvenRowStyle = CellStyle.Foreground(White)
)
