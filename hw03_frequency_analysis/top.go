package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordFreq struct {
	word string
	freq int
}

func Top10(inputString string) []string {
	// slice := strings.Fields(inputString) // если без *
	slice := strings.Fields(strings.ToLower(inputString)) // если со *
	wordFreqMap := make(map[string]int, len(slice))
	for _, v := range slice {
		v = strings.Trim(v, `,.!?-`) // если со *
		if v != "" {
			wordFreqMap[v]++
		}
	}
	wordFreqSlice := make([]wordFreq, 0, len(wordFreqMap))
	for k, v := range wordFreqMap {
		// unitToAdd := wordFreq{word: k, freq: v} удалить
		wordFreqSlice = append(wordFreqSlice, wordFreq{word: k, freq: v})
	}
	sort.Slice(wordFreqSlice, func(i, j int) bool {
		if wordFreqSlice[i].freq == wordFreqSlice[j].freq {
			return wordFreqSlice[i].word < wordFreqSlice[j].word
		} else {
			return wordFreqSlice[i].freq > wordFreqSlice[j].freq
		}
	})
	wordFreqSliceLen := 0
	if len(wordFreqSlice) > 10 {
		wordFreqSliceLen = 10
	} else {
		wordFreqSliceLen = len(wordFreqSlice)
	}
	wordFreqSlice10 := make([]string, 0, wordFreqSliceLen)
	for i := 0; i < wordFreqSliceLen; i++ {
		wordFreqSlice10 = append(wordFreqSlice10, wordFreqSlice[i].word)
	}
	return wordFreqSlice10
}

//сделать комментарии, проверить тесты
