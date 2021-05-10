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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 || offset < 0 {
		panic("Limit or offset is under zero")
	}

	fileSorce, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("opening of a sorce file is failed: %w", err)
	}
	defer fileSorce.Close()

	infoFileSorce, err := fileSorce.Stat()
	if err != nil {
		return fmt.Errorf("file info reading of a sorce file is failed: %w", err)
	}
	if infoFileSorce.Size() <= offset && offset != 0 {
		return ErrOffsetExceedsFileSize
	}
	if infoFileSorce.Size() == 0 {
		return ErrUnsupportedFile
	}

	fileDestination, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("creating of a destination file is failed: %w", err)
	}
	defer fileDestination.Close()

	err = fileDestination.Chmod(infoFileSorce.Mode())
	if err != nil {
		return fmt.Errorf("chmod from a sorce file to dest file is failed: %w", err)
	}

	var bytesToCopy int64
	if limit != 0 && limit < infoFileSorce.Size()-offset {
		bytesToCopy = limit
	} else {
		bytesToCopy = infoFileSorce.Size() - offset
	}

	bar := pb.New(int(bytesToCopy)).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	_, err = fileSorce.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("offset setting for a sorce file is failed: %w", err)
	}
	readerSorce := bar.NewProxyReader(fileSorce)
	_, err = io.CopyN(fileDestination, readerSorce, bytesToCopy)
	if err != nil {
		return fmt.Errorf("copying of bytes is failed: %w", err)
	}

	bar.Finish()
	return nil
}
