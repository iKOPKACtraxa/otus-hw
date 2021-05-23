package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitOrOffsetIsUnder0 = errors.New("limit or offset is under zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 || offset < 0 {
		return ErrLimitOrOffsetIsUnder0
	}

	fileSource, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("opening of a source file is failed: %w", err)
	}
	defer fileSource.Close()
	_, err = fileSource.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("offset setting for a source file is failed: %w", err)
	}

	infoFileSource, err := fileSource.Stat()
	if err != nil {
		return fmt.Errorf("file info reading of a source file is failed: %w", err)
	}
	if infoFileSource.Size() <= offset && offset != 0 {
		return ErrOffsetExceedsFileSize
	}
	if infoFileSource.Size() == 0 {
		return ErrUnsupportedFile
	}

	fileDestination, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("creating of a destination file is failed: %w", err)
	}
	defer fileDestination.Close()
	err = fileDestination.Chmod(infoFileSource.Mode())
	if err != nil {
		return fmt.Errorf("chmod from a source file to dest file is failed: %w", err)
	}

	bytesToCopy := infoFileSource.Size() - offset
	if limit != 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	bar := pb.New(int(bytesToCopy)).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()
	readerSource := bar.NewProxyReader(fileSource)

	_, err = io.CopyN(fileDestination, readerSource, bytesToCopy)
	if err != nil {
		return fmt.Errorf("copying of bytes is failed: %w", err)
	}

	bar.Finish()
	return nil
}
