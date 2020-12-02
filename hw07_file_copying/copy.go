package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenSourceFile        = errors.New("source file can't be opened")
	ErrGetStatOfSourceFile   = errors.New("can not get stat of source file")
	ErrSeekPosition          = errors.New("can not seek position")
	ErrCreateFile            = errors.New("can not create target file")
	ErrReadFromSourceFile    = errors.New("can not read from source file")
)

const maxChunkSize = 1 * 1024 * 1024 // 1 Mb -> B

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func evalBytesCountForCopy(size int64, offset int64, limit int64) int64 {
	bytesCount := size - offset
	if limit > 0 {
		bytesCount = min(bytesCount, limit)
	}

	return bytesCount
}

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fFrom, err := os.Open(fromPath)
	if err != nil {
		return ErrOpenSourceFile
	}
	defer fFrom.Close()

	info, err := os.Stat(fromPath)
	if err != nil {
		return ErrGetStatOfSourceFile
	}

	fSize := info.Size()

	if offset > fSize {
		return ErrOffsetExceedsFileSize
	}
	if fSize == 0 {
		return ErrUnsupportedFile
	}

	chunkSize := min(fSize, maxChunkSize)
	buff := make([]byte, chunkSize)
	_, err = fFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrSeekPosition
	}

	fTo, err := os.Create(toPath)
	if err != nil {
		return ErrCreateFile
	}
	defer fTo.Close()

	bytesCountForCopy := evalBytesCountForCopy(fSize, offset, limit)

	bar := pb.Full.Start64(bytesCountForCopy)
	fToBarProxy := bar.NewProxyWriter(fTo)

	var totalReadBytes int64 // сколько всего прочитано байт из источниа

	for {
		n, err := fFrom.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return ErrReadFromSourceFile
		}
		readBytes := int64(n)

		totalReadBytes += readBytes

		if totalReadBytes >= bytesCountForCopy {
			_, _ = fToBarProxy.Write(buff[0 : totalReadBytes-(totalReadBytes-bytesCountForCopy)])
			break
		} else {
			_, _ = fToBarProxy.Write(buff)
		}
	}

	bar.Finish()

	return nil
}
