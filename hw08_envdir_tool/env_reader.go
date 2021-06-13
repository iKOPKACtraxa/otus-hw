package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
// Files with "=" in names is skipped.
func ReadDir(dir string) (Environment, error) {
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading of dir entries has got an error: %w", err)
	}
	envMap := make(Environment, len(dirEntry))
	for _, v := range dirEntry {
		if strings.Contains(v.Name(), "=") {
			continue
		}
		file, err := os.Open(dir + "/" + v.Name())
		if err != nil {
			return nil, fmt.Errorf("opening of dir entry has got an error: %w", err)
		}
		if fileInfo, err := file.Stat(); fileInfo.Size() == 0 || err != nil {
			if err != nil {
				return nil, fmt.Errorf("reading of file stat of dir entry has got an error: %w", err)
			}
			envMap[v.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		br := bufio.NewReader(file)
		line, err := br.ReadBytes(byte('\n'))
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading of first line of dir entry has got an error: %w", err)
		}
		line = bytes.TrimRight(line, "\n \t")
		line = bytes.ReplaceAll(line, []byte{0}, []byte{10})
		envMap[v.Name()] = EnvValue{Value: string(line)}
	}
	return envMap, nil
}
