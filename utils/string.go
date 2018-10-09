package utils

import (
	"unicode/utf8"
)

func CountCharacters(characters string) (count int){
	count = 0
	for len(characters) > 0 {
		_, size := utf8.DecodeLastRuneInString(characters)
		count++
		characters = characters[:len(characters)-size]
	}
	return
}
