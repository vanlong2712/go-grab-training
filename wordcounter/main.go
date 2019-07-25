
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {
	// For who does not know :D
	// PLEASE RUN THIS IN TERMINAL: go run *.go data/story1.txt data/story2.txt
	if len(os.Args) < 2 {
		fmt.Println("Expected at least one filename.")
		os.Exit(1)
	}

	var words []Word

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(`input words you want to count. example: "a b ab c"`)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 || len(strings.Fields(text)) == 0 {
			fmt.Println("Expected at least one character.")
			os.Exit(1)
		}
		words = ConvertStringToArrayTypeWord(text)
		break
	}
	filenames := os.Args[1:]

	wordCountResult := make(chan FileWords)
	wg := sync.WaitGroup{}

	for i := 0; i < len(filenames); i++ {
		wg.Add(1)
		go WordCounter(filenames[i], words, &wg, wordCountResult)

	}

	go MonitorFileWordsChannel(&wg, wordCountResult)

	var finalWords []Word

	for result := range wordCountResult {
		a := result.Words
		sort.Slice(a, func(i, j int) bool {
			return a[i].character < a[j].character
		})
		fmt.Println("filename", result.FileName)
		for _, w := range a {
			finalWords = append(finalWords, w)
			w.print("\t")
		}
	}

	finalWords = uniqueWordSlice(finalWords)
	for _, v := range finalWords {
		v.print("Total ")
	}
}
