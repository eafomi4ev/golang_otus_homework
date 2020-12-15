package main

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dirName := "./testdata/env"
	expectedEnvs := Environment{
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

	t.Run("positive tests", func(t *testing.T) {
		t.Run("read envs from dir", func(t *testing.T) {
			envs, err := ReadDir(dirName)

			require.Equal(t, expectedEnvs, envs)
			require.NoError(t, err)
		})

		t.Run("use only files names and ignore dirs", func(t *testing.T) {
			dirIgnoreName := "./testdata/env/ingnoredir"
			os.Mkdir(dirIgnoreName, os.ModeDir)
			defer os.RemoveAll(dirIgnoreName)

			envs, err := ReadDir(dirName)

			require.NotContains(t, envs, "ingnoredir")
			require.NoError(t, err)
		})
	})

	t.Run("negative tests", func(t *testing.T) {
		t.Run("read envs from not existing dir", func(t *testing.T) {
			envs, err := ReadDir("./notExistingDir")

			require.Nil(t, envs)
			require.IsType(t, new(os.PathError), errors.Cause(err))
		})

		t.Run("read envs from from dir with '=' in the name", func(t *testing.T) {
			fileIncorrectName := "./testdata/env/FT=42"
			os.Create(fileIncorrectName)
			defer os.Remove(fileIncorrectName)

			envs, err := ReadDir("./testdata/env")

			require.Nil(t, envs)
			require.IsType(t, ErrIncorrectFileName, errors.Cause(err))
		})
	})
}
