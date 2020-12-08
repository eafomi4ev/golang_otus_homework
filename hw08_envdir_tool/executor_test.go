package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs := Environment{
		"BAR": EnvValue{
			Value:      "bar",
			NeedRemove: false,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
		"FOO": EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		},
		"HELLO": EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		},
		"UNSET": EnvValue{
			Value:      "",
			NeedRemove: true,
		},
	}

	t.Run("test with argument", func(t *testing.T) {
		cmd := []string{"./testdata/echo.sh", "foo=42"}
		exitCode := RunCmd(cmd, envs)

		require.Zero(t, exitCode)
	})

	t.Run("test without arguments", func(t *testing.T) {
		cmd := []string{"./testdata/echo.sh"}
		exitCode := RunCmd(cmd, envs)

		require.Zero(t, exitCode)
	})

	t.Run("test without arguments", func(t *testing.T) {
		cmd := []string{"./test.sh"}
		exitCode := RunCmd(cmd, envs)

		require.Zero(t, exitCode)
	})

	t.Run("handle correct exit code", func(t *testing.T) {
		cmd := []string{"./testdata/exit42.sh"}
		exitCode := RunCmd(cmd, envs)

		require.Equal(t, 42, exitCode)
	})
}
