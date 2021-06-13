package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

var envRef Environment = Environment{"BAR": EnvValue{Value: "bar", NeedRemove: false}, "EMPTY": EnvValue{Value: "", NeedRemove: false}, "FOO": EnvValue{Value: "   foo\nwith new line", NeedRemove: false}, "HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false}, "UNSET": EnvValue{Value: "", NeedRemove: true}}

func TestReadDir(t *testing.T) {
	t.Run("Regular test", func(t *testing.T) {
		env, err := ReadDir(osArgs[1])
		require.NoError(t, err, "reading of env files has got an error", err)
		require.Truef(t, cmp.Equal(envRef, env), "returned env is wrong")
	})
}
