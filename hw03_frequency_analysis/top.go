package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordFreq struct {
	word string
	freq int
}

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

// Top10 функция частотного анализа, которая получая на вход текст возвращает Топ10 встречающихся строк.
func Top10(inputString string) []string {
	var slice []string
	// предварительная подготовка из текста в срез слов, пребелы являются разделителями, регистр нижний
	if taskWithAsteriskIsCompleted {
		slice = strings.Fields(strings.ToLower(inputString))
	} else {
		slice = strings.Fields(inputString)
	}
	// для создания множества лучше использовать map
	wordFreqMap := make(map[string]int, len(slice))
	for _, v := range slice {
		if taskWithAsteriskIsCompleted {
			v = strings.Trim(v, `,.!?-`)
		}
		if v != "" {
			wordFreqMap[v]++ // подсчет частоты каждого слова
		}
	}
	// для сортировки нужно использовать срез
	wordFreqSlice := make([]wordFreq, 0, len(wordFreqMap))
	for k, v := range wordFreqMap {
		wordFreqSlice = append(wordFreqSlice, wordFreq{word: k, freq: v})
	}
	// сортировка по убыванию частоты появления слова в тексте
	// если частоты одинаковые, сортировка по возрастанию лексигрфически
	sort.Slice(wordFreqSlice, func(i, j int) bool {
		if wordFreqSlice[i].freq == wordFreqSlice[j].freq {
			return wordFreqSlice[i].word < wordFreqSlice[j].word
		}
		return wordFreqSlice[i].freq > wordFreqSlice[j].freq
	})
	// ограничение количества слов к выводу в ответ
	wordFreqSliceLenLim := 0
	if len(wordFreqSlice) > 10 {
		wordFreqSliceLenLim = 10
	} else {
		wordFreqSliceLenLim = len(wordFreqSlice)
	}
	// подготовка ответа
	wordFreqSliceAnswer := make([]string, 0, wordFreqSliceLenLim)
	for i := 0; i < wordFreqSliceLenLim; i++ {
		wordFreqSliceAnswer = append(wordFreqSliceAnswer, wordFreqSlice[i].word)
	}
	return wordFreqSliceAnswer
}
