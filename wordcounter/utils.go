package main

import (
	"strings"
	"sync"
)

func findElementInWordSlice(key string, words []Word) (int, Word) {
	for i, v := range words {
		if v.character == key {
			return i, v
		}
	}
	return -1, Word{}
}

func uniqueStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueWordSlice(wordSlice []Word) []Word {
	keys := make(map[string]bool)
	var list []Word
	for _, entry := range wordSlice {
		if _, value := keys[entry.character]; !value {
			keys[entry.character] = true
			list = append(list, entry)
		} else {
			index, oldWord := findElementInWordSlice(entry.character, list)
			if index > -1 {
				oldWord.count += entry.count
			}
			list[index] = oldWord
		}
	}
	return list
}

func ConvertStringToArrayTypeWord(str string) []Word {
	a := uniqueStringSlice(strings.Fields(strings.ToLower(str)))
	var words []Word
	for _, v := range a {
		w := Word{
			character: v,
			count:     0,
		}
		words = append(words, w)
	}
	return words
}

func ConvertStringToArray(str string) []string {
	a := strings.Fields(strings.ToLower(str))
	return a
}

func MonitorFileWordsChannel(wg *sync.WaitGroup, ch chan FileWords) {
	wg.Wait()
	close(ch)
}

func MonitorWordsChannel(wg *sync.WaitGroup, ch chan Word) {
	wg.Wait()
	close(ch)
}

func CountCharaterFromWords(word Word, words []string, waitGroup *sync.WaitGroup, out chan Word) {
	defer waitGroup.Done()
	var countWord int64
	for _, r := range words {
		if strings.Contains(r, strings.ToLower(word.character)) {
			countWord++
		}
	}
	word.count = countWord
	out <- word
}