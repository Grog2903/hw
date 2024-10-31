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
	ErrZeroFileSize          = errors.New("zero file size")
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

	tmpFile, err := os.CreateTemp("./", "e-*.txt")
	if err != nil {
		return errors.New("create temp file")
	}
	defer os.Remove(tmpFile.Name())

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if fileSize == 0 {
		destinationFile, err := os.Create(toPath)
		if err != nil {
			return fmt.Errorf("create destination file to path: %s", toPath)
		}
		defer destinationFile.Close()

		return nil
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

	_, err = io.CopyN(tmpFile, sourceFile, bytesToCopy)
	if err != nil {
		return errors.New("copy file")
	}

	_, err = tmpFile.Seek(0, io.SeekStart)
	if err != nil {
		return errors.New("seek temp file")
	}

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create destination file to path: %s", toPath)
	}
	defer destinationFile.Close()

	progressBar := pb.Simple.Start64(bytesToCopy)
	progressBarReader := progressBar.NewProxyReader(tmpFile)

	_, err = io.Copy(destinationFile, progressBarReader)
	if err != nil {
		os.Remove(toPath)
		return errors.New("copy file")
	}

	progressBar.Finish()

	return nil
}
