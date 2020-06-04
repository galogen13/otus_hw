package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("correct parameters", func(t *testing.T) {

		result := Copy("", "test.txt", 0, 1)
		require.Equal(t, result, ErrEmptyValueFrom)

		result = Copy("test.txt", "", 0, 1)
		require.Equal(t, result, ErrEmptyValueTo)

		result = Copy("test.txt", "test.txt", 0, 1)
		require.Equal(t, result, ErrFromEqualsTo)

		tmpfile := makeTestFile()
		defer os.Remove(tmpfile.Name())
		writeStringToFile(tmpfile, "test")

		resultFileName := "testdata/result.txt"
		result = Copy(tmpfile.Name(), resultFileName, 5, 0)
		require.Equal(t, result, ErrOffsetExceedsFileSize)
		require.NoFileExists(t, resultFileName)

	})

	t.Run("unsupported file", func(t *testing.T) {

		resultFileName := "testdata/out"
		result := Copy("/dev/null", resultFileName, 0, 0)
		require.Equal(t, result, ErrUnsupportedFile)
		require.NoFileExists(t, resultFileName)

	})

	t.Run("copy offset = 0 limit = 0", func(t *testing.T) {

		resultFileName := "testdata/out"
		sourceFileName := "testdata/input.txt"
		expectedFileName := "testdata/input.txt"
		copyTest(t, sourceFileName, resultFileName, expectedFileName, 0, 0)

	})

	t.Run("copy offset = 0 limit = 10", func(t *testing.T) {

		resultFileName := "testdata/out"
		sourceFileName := "testdata/input.txt"
		expectedFileName := "testdata/out_offset0_limit10.txt"
		copyTest(t, sourceFileName, resultFileName, expectedFileName, 0, 10)

	})
}

func makeTestFile() *os.File {
	tmpfile, err := ioutil.TempFile("testdata", "test.")
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile
}

func writeStringToFile(f *os.File, content string) {

	if _, err := f.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func copyTest(t *testing.T, sourceFileName, resultFileName, expectedFileName string, offset, limit int64) {
	result := Copy(sourceFileName, resultFileName, offset, limit)
	if result == nil {
		defer os.Remove(resultFileName)
	}
	require.Nil(t, result)
	require.FileExists(t, resultFileName)

	expectedFile, err := os.Open(expectedFileName)
	defer expectedFile.Close()
	b1, err := ioutil.ReadAll(expectedFile)
	if err != nil {
		log.Fatal(err)
	}

	resultFile, err := os.Open(resultFileName)
	defer resultFile.Close()
	b2, err := ioutil.ReadAll(resultFile)
	if err != nil {
		log.Fatal(err)
	}
	require.Equal(t, string(b1), string(b2))
}
