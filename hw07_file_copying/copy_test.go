package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("bashTests", func(t *testing.T) {
		cmd := exec.Command("bash", "test.sh")

		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("Output:\n%s", string(output))
		}

		require.NoError(t, err)
	})

	t.Run("unsupportedFile", func(t *testing.T) {
		err := Copy("testdata/input1.txt", "out.txt", 0, 0)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("bigOffset", func(t *testing.T) {
		file, _ := os.Open("testdata/input.txt")
		fileInfo, _ := file.Stat()
		fileSize := fileInfo.Size()
		offset := fileSize + 1

		err := Copy("testdata/input.txt", "out.txt", offset, 0)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})
}
