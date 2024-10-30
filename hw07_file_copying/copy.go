package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o400)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer sourceFile.Close()

	stat, err := sourceFile.Stat()
	if err != nil {
		return errors.New("get source file stat")
	}
	fileSize := stat.Size()

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create destination file to path: %s", toPath)
	}
	defer destinationFile.Close()

	if offset > stat.Size() {
		os.Remove(toPath)
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		_, err := sourceFile.Seek(offset, io.SeekStart)
		if err != nil {
			os.Remove(toPath)
			return errors.New("seek source file offset")
		}
	}

	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	progressBar := pb.Simple.Start64(bytesToCopy)
	progressBarReader := progressBar.NewProxyReader(io.LimitReader(sourceFile, bytesToCopy))

	_, err = io.CopyN(destinationFile, progressBarReader, bytesToCopy)
	if err != nil {
		os.Remove(toPath)
		return errors.New("copy file")
	}

	progressBar.Finish()

	return nil
}
