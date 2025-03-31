package utils

import (
	"cmd/tool/models"
	"cmd/tool/view"
	"fmt"
	"regexp"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func CheckVariables(value string, config models.ConfigFile, p *tea.Program) string {
	if strings.Contains(value, "{{remoteUser}}") {
		value = strings.ReplaceAll(value, "{{remoteUser}}", config.RemoteUser)
	}

	if strings.Contains(value, "{{remoteHost}}") {
		value = strings.ReplaceAll(value, "{{remoteHost}}", config.RemoteHost)
	}

	inRe := regexp.MustCompile(`\{{(input:)([^}]+)\}}`)
	inMatches := inRe.FindAllStringSubmatch(value, -1)

	if len(inMatches) > 0 {
		var wg sync.WaitGroup
		wg.Add(1)

		var input = ""
		p.Send(view.Input{Input: &input, Msg: inMatches[0][2], Wg: &wg})

		wg.Wait()
		value = input
	}

	re := regexp.MustCompile(`\{{([^}]+)\}}`)
	matches := re.FindAllStringSubmatch(value, -1)

	for _, match := range matches {
		if len(match) == 2 {
			inMatches := inRe.FindAllStringSubmatch(config.Variables[match[1]], -1)
			if len(inMatches) == 0 {
				value = strings.Replace(value, fmt.Sprintf("{{%s}}", match[1]), config.Variables[match[1]], 1)
			}
		}
	}

	return value
}
