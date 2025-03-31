package view

import (
	"cmd/tool/models"
	"cmd/tool/style"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	consoleLength = 5
)

type model struct {
	spinner        spinner.Model
	runs           []uint
	results        []models.ResultMsg
	terminalOut    []models.TerminalOut
	quitting       bool
	textInput      textinput.Model
	input          models.InputMsg
	waitingInput   bool
	addedVariables [][]string
	sshConnection  models.SshConnectionMsg
}

func NewModel() model {
	s := spinner.New()
	s.Style = style.SpinnerStyle
	s.Spinner = spinner.Dot
	ti := textinput.New()
	ti.Focus()
	return model{
		spinner:     s,
		results:     make([]models.ResultMsg, consoleLength),
		terminalOut: make([]models.TerminalOut, consoleLength),
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
			m.addedVariables = append(m.addedVariables, []string{
				m.input.Msg,
				m.textInput.Value(),
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
	case models.ResultMsg:
		if !m.quitting {
			if msg.StartText {
				m.results = append(m.results, msg)
			} else {
				m.results[len(m.results)-1] = msg
			}
			return m, nil
		}
	case models.SshConnectionMsg:
		m.sshConnection = msg
		return m, nil
	case models.TerminalOut:
		m.terminalOut = append(m.terminalOut[1:], msg)
		return m, nil
	case models.InputMsg:
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
		s += sshConnectionView(&m)
	}

	// Input Screen
	if m.waitingInput {
		s += inputView(&m)
	} else {
		s += spinnerView(&m)
	}

	return style.AppStyle.Render(s)
}
