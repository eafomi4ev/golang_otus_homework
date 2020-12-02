package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const inputFileName = "./testdata/input.txt"
const outputFileName = "./testdata/result.txt"

func TestCopy(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(outputFileName)
	})

	t.Run("Should return nil if all are ok", func(t *testing.T) {
		err := Copy(inputFileName, outputFileName, 0, 0)

		require.Nil(t, err)
	})

	t.Run("offset 0, limit 0", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 0, 0)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit0.txt")

		require.Equal(t, expectedText, resultText)
	})

	t.Run("offset 0, limit 10", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 0, 10)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit10.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 0, limit 1000", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 0, 1000)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 0, limit 10000", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 0, 10000)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit10000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 100, limit 1000", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 100, 1000)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset100_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("offset 6000, limit 1000", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 6000, 1000)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset6000_limit1000.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("Offset more than file size. Should return error", func(t *testing.T) {
		err := Copy(inputFileName, outputFileName, 0, 0)

		require.Error(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Limit more than file size. Should copy file up to EOF", func(t *testing.T) {
		Copy(inputFileName, outputFileName, 0, 9999999999)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit0.txt")

		require.Equal(t, string(expectedText), string(resultText))
	})

	t.Run("Limit more than file size and offset isn't 0. Should copy file up to EOF", func(t *testing.T) {
		offset := int64(6600)

		Copy(inputFileName, outputFileName, offset, 9999999999)

		resultText, _ := ioutil.ReadFile(outputFileName)
		expectedText, _ := ioutil.ReadFile("./testdata/out_offset0_limit0.txt")

		require.Equal(t, string(expectedText[offset:]), string(resultText))
	})

	t.Run("File with zero size. Should return error", func(t *testing.T) {
		err := Copy("/dev/urandom", outputFileName, 0, 0)

		require.Error(t, ErrUnsupportedFile, err)
	})
}
