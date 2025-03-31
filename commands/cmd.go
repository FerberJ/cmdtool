package commands

import (
	"bufio"
	"cmd/tool/models"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

func ExecCmd(cmd string, args []string, paramsCount uint, params []string, p *tea.Program) error {
	if int(paramsCount) == len(params) {

		for argI, arg := range args {
			count := strings.Count(arg, "%s")

			iRemoved := 0
			for i := 0; i < count; i++ {
				realI := i - iRemoved
				args[argI] = strings.Replace(args[argI], "%s", params[realI], 1)
				params = append(params[:realI], params[realI+1:]...)
				iRemoved += 1
			}
		}

		for _, param := range params {
			args = append(args, param)
		}

		execCmd := exec.Command(cmd, args...)
		execCmd.Stdin = os.Stdin

		// var outBuffer, errBuffer bytes.Buffer
		stdoutPipe, _ := execCmd.StdoutPipe()
		stderrPipe, _ := execCmd.StderrPipe()

		latestError := ""

		// Capture and send output in real-time
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				p.Send(models.TerminalOut{Text: scanner.Text()})
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				latestError += "" + scanner.Text()
				p.Send(models.TerminalOut{Text: scanner.Text()})
			}
		}()

		// Start execution
		if err := execCmd.Start(); err != nil {
			return err
		}

		// Ensure all output is captured
		if err := execCmd.Wait(); err != nil {
			return fmt.Errorf("%v", latestError)
		}

		return nil
	}

	return fmt.Errorf("the number of params (%v) doesn't match the specifications: %v", len(params), paramsCount)
}

func ExecSsh(client *ssh.Client, cmd string, args []string, paramsCount uint, params []string) error {
	if int(paramsCount) == len(params) {
		cmdArr := make([]string, 0)
		cmdArr = append(cmdArr, cmd)

		for argI, arg := range args {
			count := strings.Count(arg, "%s")

			iRemoved := 0
			for i := 0; i < count; i++ {
				realI := i - iRemoved
				args[argI] = strings.Replace(args[argI], "%s", params[realI], 1)
				params = append(params[:realI], params[realI+1:]...)
				iRemoved += 1
			}
		}

		for _, arg := range args {
			cmdArr = append(cmdArr, arg)
		}
		for _, param := range params {
			cmdArr = append(cmdArr, param)
		}
		cmdStr := strings.Join(cmdArr, " ")

		return executeCmd(client, cmdStr)
	}

	return fmt.Errorf("the number of params (%v) doenst match the specifications: %v", len(params), paramsCount)
}

func executeCmd(client *ssh.Client, cmdString string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmdString)
	if err != nil {
		return fmt.Errorf("%s: %s", err, output)
	}

	return nil
}
