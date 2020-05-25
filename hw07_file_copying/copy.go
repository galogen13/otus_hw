package main

import (
	"errors"
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
//nolint:funlen
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
		return err
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
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

	resFile, _ := os.Create(toPath)
	if err != nil {
		return err
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
	barWriter := bar.NewProxyWriter(resFile)

	for offset < fileSize {
		_, err := file.Seek(offset, 0)
		if err != nil {
			return err
		}
		bufSize := int64(math.Min(float64(bufSize), float64(fileSize-offset)))

		written, err := io.CopyN(barWriter, file, bufSize)
		if err != nil {
			return err
		}
		offset += written

		bar.Write()
	}
	resFile.Close()
	bar.Finish()
	return nil
}
