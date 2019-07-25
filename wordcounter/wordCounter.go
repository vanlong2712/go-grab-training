package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const chunkSize int32 = 32 * 1024


func WordCounter(filename string, words []Word, waitGroup *sync.WaitGroup, out chan<- FileWords) {
	defer waitGroup.Done()

	fmt.Sprintf("wordReader started for %s", filename)

	// Attempt to open the file
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Sprintf("Error opening file %s", filename)
		return
	}

	// Start processing data from the file
	for {
		buf := make([]byte, chunkSize) // Create a buffer
		bytes, err := file.Read(buf)   // Read data into the buffer
		if bytes > 0 {
			str := string(buf) // Convert bytes to a string
			wordsFromFile := ConvertStringToArray(str)
			wordChannel := make(chan Word)
			wg := sync.WaitGroup{}
			for _, word := range words {
				wg.Add(1)
				go CountCharaterFromWords(word, wordsFromFile, &wg, wordChannel)
			}
			go MonitorWordsChannel(&wg, wordChannel)
			var tempWords []Word
			for rs := range wordChannel {
				tempWords = append(tempWords, rs)
			}
			words = tempWords
		}

		// Hit the end of file
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Sprintf("Error reading file %s: %s", filename, err)
			break
		}
	}

	out <- FileWords{file.Name(), words}
	file.Close()
}


