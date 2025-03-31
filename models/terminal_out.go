package models

import (
	"cmd/tool/style"
	"fmt"
	"strings"
)

type TerminalOut struct {
	Text string
}

func (t TerminalOut) String() string {
	if t.Text == "" {
		return style.DotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("%s", t.Text)
}
