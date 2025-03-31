package models

import (
	"cmd/tool/style"
	"fmt"
	"strings"
	"time"
)

type ResultMsg struct {
	Duration  time.Duration
	Text      string
	StartText bool
}

func (r ResultMsg) String() string {
	if r.Duration == 0 && r.Text == "" {
		return style.DotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("%s %s", r.Text,
		style.DurationStyle.MaxHeight(1).Render(r.Duration.String()))
}
