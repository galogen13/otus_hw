package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyValueFrom        = errors.New("empty value -from")
	ErrEmptyValueTo          = errors.New("empty value -to")
	ErrFromEqualsTo          = errors.New("-from and -to are equal")
)

// Copy .
func Copy(fromPath string, toPath string, offset, limit int64) error {
	const bufSize int64 = 1024 * 1024

	if fromPath == "" {
		return ErrEmptyValueFrom
	}

	if toPath == "" {
		return ErrEmptyValueTo
	}

	if fromPath == toPath {
		return ErrFromEqualsTo
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open file info: %v", err)
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fi.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileSize
	}

	if limit+offset < fileSize {
		fileSize = limit + offset
	}

	tmpl := "{{counters . }} {{percent . }} \n"
	bar := pb.New64(limit)
	bar.SetTemplateString(tmpl)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	bar.SetWriter(os.Stdout)
	bar.Set(pb.Static, true)
	bar.Start()
	bar.Write()

	resFile, _ := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	barReader := bar.NewProxyWriter(resFile)

	for offset < fileSize {
		_, err := file.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("failed to seek: %v", err)
		}
		bufSize := int64(math.Min(float64(bufSize), float64(fileSize-offset)))
		buf := make([]byte, bufSize)
		read, err := file.Read(buf)
		offset += int64(read)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read: %v", err)
		}

		_, err = barReader.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}
		bar.Write()
	}
	resFile.Close()
	bar.Finish()
	return nil
}
