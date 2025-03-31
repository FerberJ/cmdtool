package main

import (
	"cmd/tool/commands"
	"cmd/tool/config"
	"cmd/tool/utils"
	"cmd/tool/view"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

func containsRunIndex(runIndex uint, list []uint) bool {
	if len(list) == 0 {
		return true
	}
	for _, i := range list {
		if i == runIndex {
			return true
		}
	}
	return false
}

func main() {
	p := tea.NewProgram(view.NewModel())

	go func() {
		config, _ := config.GetConfig(p)

		remoteUser := config.ConfigFile.RemoteUser
		remoteHost := config.ConfigFile.RemoteHost
		idRsaPath := config.ConfigFile.IdRsaPath

		p.Send(view.SshConnectionMsg{User: remoteUser, Address: remoteHost, Pending: true, Success: false})
		client, err := sshConnect(remoteUser, remoteHost, "22", idRsaPath)
		if err != nil {
			p.Send(view.SshConnectionMsg{User: remoteUser, Address: remoteHost, Pending: false, Success: false, Error: err.Error()})
			p.Send(true)
			return
		}
		defer client.Close()
		p.Send(view.SshConnectionMsg{User: remoteUser, Address: remoteHost, Pending: false, Success: true})

		for _, runCmd := range config.ConfigFile.RunCmds {
			start := time.Now()

			if !containsRunIndex(runCmd.RunIndex, config.RunIndexList) {
				continue
			}

			startMsg := fmt.Sprintf("⏳ %v", runCmd.Description)
			p.Send(view.ResultMsg{Text: startMsg, Duration: 0, StartText: true})
			ticker := time.NewTicker(1 * time.Second)
			processTime := time.Now()
			quit := make(chan bool)
			go func() {
				for {
					select {
					case <-ticker.C:
						elapsedTime := time.Since(processTime)
						fullSeconds := elapsedTime.Truncate(time.Second)
						p.Send(view.ResultMsg{Text: startMsg, Duration: fullSeconds, StartText: false})
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()

			for i, param := range runCmd.Params {
				runCmd.Params[i] = utils.CheckVariables(param, config.ConfigFile, p)
			}
			cmd := config.ConfigFile.Cmds[runCmd.Cmd]

			args := make([]string, len(cmd.Args))
			copy(args, cmd.Args)

			if cmd.Type == "exec" {

				err := commands.ExecCmd(cmd.Cmd, args, cmd.Params, runCmd.Params, p)
				if err != nil {
					msg := fmt.Sprintf("❌ %v, error: %s", runCmd.Description, view.Error.Render(err.Error()))
					quit <- true
					p.Send(view.ResultMsg{Text: msg, Duration: 0, StartText: false})
					if runCmd.StopAfterFail {
						p.Send(true)
					}
					continue
					// p.Send(true)
				}
			}
			if cmd.Type == "ssh" {
				err := commands.ExecSsh(client, cmd.Cmd, args, cmd.Params, runCmd.Params)
				if err != nil {
					msg := fmt.Sprintf("❌ %v, error: %s", runCmd.Description, view.Error.Render(err.Error()))
					quit <- true
					p.Send(view.ResultMsg{Text: msg, Duration: 0, StartText: false})
					if runCmd.StopAfterFail {
						p.Send(true)
					}
					continue
					// p.Send(true)
				}
			}

			quit <- true
			stop := time.Now()
			duration := stop.Sub(start)
			msg := fmt.Sprintf("✅ %v", runCmd.Description)

			p.Send(view.ResultMsg{Text: msg, Duration: duration, StartText: false})
		}

		p.Send(true)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func sshConnect(user, host, port, keyPath string) (*ssh.Client, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
}
