package utils

import (
	"fmt"
	"unicode/utf8"
)

func CountCharacters(characters string) (count int){
	count = 0
	for len(characters) > 0 {
		r, size := utf8.DecodeLastRuneInString(characters)
		count++
		fmt.Printf("%c %v\n", r, size)
		characters = characters[:len(characters)-size]
	}
	return
}
