package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

const (
	fromPathOfSourceFile = "testdata/input.txt"
	fromPathRand         = "/dev/urandom"
	ErrorIs              = "ErrorIs"
	ErrorAs              = "ErrorAs"
	dirTest              = "testdata"
	thisIsWrongPath      = "thisIsWrongPath"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		offset    int64
		limit     int64
		refFile   string
		nameOfSub string
	}{
		{0, 0, "testdata/out_offset0_limit0.txt", "offset=0, limit=0"},
		{0, 10, "testdata/out_offset0_limit10.txt", "offset=0, limit=10"},
		{0, 1000, "testdata/out_offset0_limit1000.txt", "offset=0, limit=1000"},
		{0, 10000, "testdata/out_offset0_limit10000.txt", "offset=0, limit=10000"},
		{100, 1000, "testdata/out_offset100_limit1000.txt", "offset=100, limit=1000"},
		{6000, 1000, "testdata/out_offset6000_limit1000.txt", "offset=6000, limit=1000"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprint("Various copying of the file. Subtest: ", tc.nameOfSub), func(t *testing.T) {
			tempFile, err := os.CreateTemp(dirTest, "")
			require.NoError(t, err, "at temp file creation got an error: ", err)
			err = Copy(fromPathOfSourceFile, tempFile.Name(), tc.offset, tc.limit)
			require.NoError(t, err, "at copying got an error: ", err)
			cmp := equalfile.New(nil, equalfile.Options{})
			equal, err := cmp.CompareFile(tempFile.Name(), tc.refFile)
			require.NoError(t, err, "in comparing got an error: ", err)
			require.Truef(t, equal, "Source and dest files not equal")
			os.Remove(tempFile.Name())
		})
	}
}

func TestCopyForErrors(t *testing.T) {
	tests := []struct {
		branch    string
		offset    int64
		limit     int64
		fromPath  string
		err       error
		nameOfSub string
	}{
		{ErrorIs, 6618, 0, fromPathOfSourceFile, ErrOffsetExceedsFileSize, "offset=6618, error=ErrOffsetExceedsFileSize"},
		{ErrorIs, 0, 0, fromPathRand, ErrUnsupportedFile, "fromPath='/dev/urandom', error=ErrUnsupportedFile"},
		{ErrorAs, 0, 0, thisIsWrongPath, &fs.PathError{}, "fromPath='thisIsWrongPath', error=&fs.PathError{}"},
		{ErrorIs, -1, 0, fromPathOfSourceFile, ErrLimitOrOffsetIsUnder0, "offset=-1, error=ErrLimitOrOffsetIsUnder0"},
		{ErrorIs, 0, -1, fromPathOfSourceFile, ErrLimitOrOffsetIsUnder0, "limit=-1, error=ErrLimitOrOffsetIsUnder0"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprint("Test for errors. Subtest: ", tc.nameOfSub), func(t *testing.T) {
			tempFile, err := os.CreateTemp(dirTest, "")
			require.NoError(t, err, "at temp file creation got an error: ", err)
			err = Copy(tc.fromPath, tempFile.Name(), tc.offset, tc.limit)
			if tc.branch == ErrorIs {
				require.ErrorIs(t, err, tc.err, "another error received: ", err)
			}
			if tc.branch == ErrorAs {
				require.ErrorAs(t, err, &tc.err, "another error received: ", err)
			}
			os.Remove(tempFile.Name())
		})
	}
}
