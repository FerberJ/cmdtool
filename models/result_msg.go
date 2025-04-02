package models

import (
	"cmd/tool/style"
	"fmt"
	"strings"
	"time"
)

type ResultMsg struct {
	Duration   time.Duration
	Text       string
	IsFirstMsg bool
	Pending    bool
	Success    bool
	Error      string
}

func (r ResultMsg) String() string {
	if r.Duration == 0 && r.Text == "" {
		return style.DotStyle.Render(strings.Repeat(".", 30))
	}

	cmd := fmt.Sprintf("%s %s", r.Text,
		style.DurationStyle.MaxHeight(1).Render(r.Duration.String()))

	if !r.Pending {
		if r.Success {
			cmd = "✅ " + cmd
		} else {
			cmd = "❌ " + cmd + " "
			cmd += style.Error.Render(r.Error)
		}
	}
	return cmd
}
