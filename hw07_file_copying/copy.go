package main

import (
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
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
		return errors.Wrap(ErrOpenSourceFile, err.Error())
	}
	defer fFrom.Close()

	info, err := os.Stat(fromPath)
	if err != nil {
		return errors.Wrap(ErrGetStatOfSourceFile, err.Error())
	}

	fSize := info.Size()
	if offset > fSize {
		return errors.Wrap(ErrOffsetExceedsFileSize, "")
	}
	if fSize == 0 {
		return errors.Wrap(ErrUnsupportedFile, "")
	}

	_, err = fFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.Wrap(ErrSeekPosition, err.Error())
	}

	fTo, err := os.Create(toPath)
	if err != nil {
		return errors.Wrap(ErrCreateFile, err.Error())
	}
	defer fTo.Close()

	bytesCountForCopy := evalBytesCountForCopy(fSize, offset, limit)
	chunkSize := min(bytesCountForCopy, maxChunkSize)
	var totalReadBytes int64 // сколько всего прочитано байт из источника

	bar := pb.Full.Start64(bytesCountForCopy)
	fToBarProxy := bar.NewProxyWriter(fTo)

	for {
		n, err := io.CopyN(fToBarProxy, fFrom, chunkSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return errors.Wrap(ErrReadFromSourceFile, err.Error())
		}

		totalReadBytes += n
		if totalReadBytes >= bytesCountForCopy {
			break
		}
	}

	bar.Finish()

	return nil
}
