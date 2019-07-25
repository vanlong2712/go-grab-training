package main

import "fmt"

type FileWords struct {
	FileName string
	Words    []Word
}

type Word struct {
	character string
	count     int64
}

func (w Word) print(prefix string) {
	fmt.Printf(prefix + "character %s: %d\n", w.character, w.count)
}
