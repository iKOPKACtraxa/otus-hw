package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
// Also returns exit code of command and 1 if got a setenv problem.
func RunCmd(cmd []string, env Environment) (returnCode int, err error) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		err := os.Setenv(k, v.Value)
		if err != nil {
			return 1, fmt.Errorf("setenv at key:%v value:%v of environment from dir has got an error: %w", k, v.Value, err)
		}
	}
	cmdToRun := exec.Command(cmd[2], cmd[3:]...) // #nosec G204
	cmdToRun.Stdout = os.Stdout
	err = cmdToRun.Run()
	var p *exec.ExitError
	if errors.As(err, &p) {
		return err.(*exec.ExitError).ProcessState.ExitCode(), err // nolint:errorlint
	}
	return 0, nil
}
