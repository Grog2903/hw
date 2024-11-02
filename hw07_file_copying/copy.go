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
	sourceFile, fileSize, err := openSourceFile(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if fileSize == 0 {
		return createEmptyDestinationFile(toPath)
	}

	tmpFile, err := os.CreateTemp("", "*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	bytesToCopy := getBytesToCopy(fileSize, offset, limit)

	progressBar := pb.Simple.Start64(bytesToCopy)
	defer progressBar.Finish()

	if err = copySourceToTemp(sourceFile, tmpFile, offset, bytesToCopy, progressBar); err != nil {
		return err
	}

	return renameTempFile(toPath, tmpFile.Name())
}

func renameTempFile(path string, name string) error {
	return os.Rename(name, path)
}

func copySourceToTemp(
	sourceFile *os.File,
	tmpFile *os.File,
	offset int64,
	bytesToCopy int64,
	progressBar *pb.ProgressBar,
) error {
	if offset > 0 {
		if _, err := sourceFile.Seek(offset, io.SeekStart); err != nil {
			return errors.New("seek source file offset")
		}
	}

	progressBarReader := progressBar.NewProxyReader(sourceFile)
	if _, err := io.CopyN(tmpFile, progressBarReader, bytesToCopy); err != nil {
		return errors.New("copy file")
	}

	return nil
}

func getBytesToCopy(fileSize int64, offset int64, limit int64) int64 {
	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	return bytesToCopy
}

func openSourceFile(path string) (*os.File, int64, error) {
	sourceFile, err := os.OpenFile(path, os.O_RDONLY, 0o400)
	if err != nil {
		return nil, 0, ErrUnsupportedFile
	}

	stat, err := sourceFile.Stat()
	if err != nil {
		return nil, 0, errors.New("get source file stat")
	}
	fileSize := stat.Size()

	return sourceFile, fileSize, nil
}

func createEmptyDestinationFile(path string) error {
	dstFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create destination file to path: %s", path)
	}
	defer dstFile.Close()
	return nil
}
