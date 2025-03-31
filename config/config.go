package config

import (
	"cmd/tool/models"
	"cmd/tool/utils"
	"flag"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func GetConfig(p *tea.Program) (models.Config, error) {
	var config models.Config
	var configFile models.ConfigFile

	var runIndexList UintSlice
	var yamlFile string

	flag.Var(&runIndexList, "runIndexes", "Comma-separated list of unsigned integers, witch run numver can go")
	flag.StringVar(&yamlFile, "yamlFile", "commands.yaml", "Path of the yaml file")

	flag.Parse()

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return config, err
	}

	// Parse the YAML data
	err = yaml.Unmarshal(data, &configFile)
	if err != nil {
		return config, err
	}

	configFile = editVariables(configFile, runIndexList, p)

	config = models.Config{
		ConfigFile:   configFile,
		RunIndexList: runIndexList,
	}

	return config, nil
}

func editVariables(config models.ConfigFile, runListIndex UintSlice, p *tea.Program) models.ConfigFile {
	for {
		edited := false

		for key, value := range config.Variables {
			newValue := utils.CheckVariables(value, config, p)
			if strings.Compare(value, newValue) != 0 {
				edited = true
				config.Variables[key] = newValue
			}
		}

		for i, cmd := range config.RunCmds {
			// Only check the variables of the cmds that are beeing run
			if len(runListIndex) == 0 || slices.Contains(runListIndex, cmd.RunIndex) {
				for y, value := range cmd.Params {
					newValue := utils.CheckVariables(value, config, p)
					if strings.Compare(value, newValue) != 0 {
						edited = true
						config.RunCmds[i].Params[y] = newValue
					}

				}
			}
		}

		if !edited {
			break
		}
	}

	return config
}
