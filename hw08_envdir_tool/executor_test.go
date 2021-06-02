package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	initEnv Environment = Environment{
		"HELLO": EnvValue{Value: "SHOULD_REPLACE"},
		"FOO":   EnvValue{Value: "SHOULD_REPLACE"},
		"UNSET": EnvValue{Value: "SHOULD_REMOVE"},
		"ADDED": EnvValue{Value: "from original env"},
		"EMPTY": EnvValue{Value: "SHOULD_BE_EMPTY"},
	}
	osArgs            []string = []string{"go-envdir", "testdata/env", "bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
	osArgsForExitCode []string = []string{"go-envdir", "", "cd", "wrongPath"}
	refString         string   = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`
)

func TestRunCmd(t *testing.T) {
	t.Run("Regular test", func(t *testing.T) {
		for k, v := range initEnv {
			err := os.Setenv(k, v.Value)
			require.NoError(t, err, "setting of initEnv has got an error", err)
		}

		env, err := ReadDir(osArgs[1])
		require.NoError(t, err, "reading of env files has got an error", err)

		readerFromPipe, writerToPipe, err := os.Pipe()
		require.NoError(t, err, "making of pipe has got an error", err)
		realStdout := os.Stdout
		os.Stdout = writerToPipe

		returnCode, err := RunCmd(osArgs, env)
		require.NoError(t, err, "RunCmd has got an error", err)
		require.Truef(t, returnCode == 0, "return code is not a 0 it is: %v", returnCode)

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, readerFromPipe)
			readerFromPipe.Close()
			outC <- buf.String()
		}()
		writerToPipe.Close()
		os.Stdout = realStdout
		outputOfRunCmd := <-outC
		close(outC)
		require.Truef(t, refString == outputOfRunCmd, "program output and reference string is not equal, needs:\n%v\ngot\n%v", refString, outputOfRunCmd)
	})
}

func TestRunCmdForExitcodes(t *testing.T) {
	t.Run("Test for exitcodes", func(t *testing.T) {
		returnCode, _ := RunCmd(osArgsForExitCode, nil)
		require.Truef(t, returnCode == 1, "return code is not 1")
	})
}
