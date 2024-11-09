package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("emptyCommand", func(t *testing.T) {
		env := Environment{}
		returnCode := RunCmd([]string{}, env)

		require.Equal(t, returnCode, 1)
	})

	t.Run("success", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": {Value: "123", NeedRemove: false},
		}
		cmd := []string{"sh", "-c", "echo $TEST_VAR"}
		returnCode := RunCmd(cmd, env)

		require.Equal(t, returnCode, 0)
	})

	t.Run("removeEnv", func(t *testing.T) {
		err := os.Setenv("REMOVE_VAR", "should_be_removed")
		require.NoError(t, err)

		value := os.Getenv("REMOVE_VAR")
		require.Equal(t, "should_be_removed", value)

		env := Environment{
			"REMOVE_VAR": {Value: "", NeedRemove: true},
		}
		cmd := []string{"sh", "-c", "echo $REMOVE_VAR"}

		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)

		value = os.Getenv("REMOVE_VAR")
		require.Equal(t, "", value)
	})
}
