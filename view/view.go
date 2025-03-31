package view

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.UnsetMargins()
	durationStyle = dotStyle
	appStyle      = lipgloss.NewStyle()
	Error         = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	title         = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("10")).Width(60).PaddingLeft(20)
	consoleLength = 5
)

type SshConnectionMsg struct {
	User    string
	Address string
	Pending bool
	Success bool
	Error   string
}

func (s SshConnectionMsg) String() string {
	str := s.User + "\t" + s.Address
	if s.Pending {
		str += "\t⏳"
	} else if s.Success {
		str += "\t✅"
	} else {
		str += "\t❌"
		str += "\n"
		str += Error.Render(s.Error)
	}
	return str
}

type ResultMsg struct {
	Duration  time.Duration
	Text      string
	StartText bool
}

func (r ResultMsg) String() string {
	if r.Duration == 0 && r.Text == "" {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("%s %s", r.Text,
		durationStyle.MaxHeight(1).Render(r.Duration.String()))
}

type TerminalOut struct {
	Text string
}

func (t TerminalOut) String() string {
	if t.Text == "" {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("%s", t.Text)
}

type Input struct {
	Input *string
	Msg   string
	Wg    *sync.WaitGroup
}

func (i Input) String() string {
	return i.Msg
}

type addedVars struct {
	input string
	msg   string
}

type model struct {
	spinner        spinner.Model
	runs           []uint
	results        []ResultMsg
	terminalOut    []TerminalOut
	quitting       bool
	textInput      textinput.Model
	input          Input
	waitingInput   bool
	addedVariables []addedVars
	sshConnection  SshConnectionMsg
}

func NewModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	ti := textinput.New()
	ti.Focus()
	return model{
		spinner:     s,
		results:     make([]ResultMsg, consoleLength),
		terminalOut: make([]TerminalOut, consoleLength),
		textInput:   ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			*m.input.Input = m.textInput.Value()
			m.addedVariables = append(m.addedVariables, addedVars{
				input: m.textInput.Value(),
				msg:   m.input.Msg,
			})
			m.input.Wg.Done()
			m.waitingInput = false
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, nil
		}
	case bool:
		m.quitting = true
		return m, tea.Quit
	case ResultMsg:
		if !m.quitting {
			if msg.StartText {
				m.results = append(m.results, msg)
			} else {
				m.results[len(m.results)-1] = msg
			}
			return m, nil
		}
	case SshConnectionMsg:
		m.sshConnection = msg
		return m, nil
	case TerminalOut:
		m.terminalOut = append(m.terminalOut[1:], msg)
		return m, nil
	case Input:
		m.input = msg
		m.waitingInput = true
		m.textInput.SetValue("")
		m.textInput.Placeholder = msg.Msg
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	s += header()
	s += "\n\n"

	if m.sshConnection.Address != "" {
		s += title.Render("SSH-Connection:")
		s += "\n\n"
		s += m.sshConnection.String()
		s += "\n\n"
	}

	// Input Screen
	if m.waitingInput {

		if len(m.addedVariables) > 0 {
			s += lipgloss.NewStyle().Bold(true).Render("Inputs:")
			s += "\n"
			for _, addVar := range m.addedVariables {
				s += dotStyle.Render(addVar.msg + ": ")
				s += addVar.input
				s += "\n"
			}
		}

		s += "\n"
		s += "Variable input: \n"

		s += m.textInput.View()
		s += "\n"
		s += helpStyle.Render("enter add variable • ctrl+c quit")
	} else {
		s += spinnerView(&m)
	}

	return appStyle.Render(s)
}

func spinnerView(m *model) string {
	var s string

	if len(m.addedVariables) > 0 {
		s += title.Render("Inputs:")
		s += "\n\n"
		for _, addVar := range m.addedVariables {
			s += dotStyle.Render(addVar.msg + ": ")
			s += addVar.input
			s += "\n"
		}
	}

	s += "\n\n"

	s += title.Render("Commands:")
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

	s += title.Render("Terminal out:")
	s += "\n\n"
	for i := len(m.terminalOut) - consoleLength; i < len(m.terminalOut); i++ {
		s += m.terminalOut[i].String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("ctrl+c quit")
	}

	if m.quitting {
		s += "\n"
	}

	return s
}

func header() string {
	return `
	___                                          _ _____            _ 
   / __\___  _ __ ___  _ __ ___   __ _ _ __   __| /__   \___   ___ | |
  / /  / _ \| '_ ' _ \| '_ ' _ \ / _' | '_ \ / _' | / /\/ _ \ / _ \| |
 / /__| (_) | | | | | | | | | | | (_| | | | | (_| |/ / | (_) | (_) | |
 \____/\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|\__,_|\/   \___/ \___/|_|
																	  
 `
}
