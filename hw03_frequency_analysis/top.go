package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

const targetWordsCount = 10

type wordFrequency struct {
	word      string
	frequency int
}

func calculateWordsFrequency(words []string) map[string]int {
	wordsAndFrequency := make(map[string]int)

	for _, word := range words {
		wordsAndFrequency[word]++
	}

	return wordsAndFrequency
}

func splitMapToKeyValueSlice(wordsAndFrequencyMap map[string]int) []wordFrequency {
	result := make([]wordFrequency, len(wordsAndFrequencyMap))

	var i int
	for key, value := range wordsAndFrequencyMap {
		item := wordFrequency{key, value}
		result[i] = item
		i++
	}

	return result
}

func Top10(initialText string) []string {
	words := strings.Fields(initialText)

	wordsAndFrequencyMap := calculateWordsFrequency(words)
	wordFrequencySlice := splitMapToKeyValueSlice(wordsAndFrequencyMap)

	sort.Slice(wordFrequencySlice, func(i, j int) bool {
		return wordFrequencySlice[i].frequency > wordFrequencySlice[j].frequency
	})

	var count int
	if len(wordFrequencySlice) < targetWordsCount {
		count = len(wordFrequencySlice)
	} else {
		count = targetWordsCount
	}

	mostFrequencyWords := make([]string, count)
	for i := 0; i < count; i++ {
		mostFrequencyWords[i] = wordFrequencySlice[i].word
	}

	return mostFrequencyWords
}
