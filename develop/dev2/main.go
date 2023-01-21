package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	s := PackedString("qwe\\45")
	fmt.Println(s.Unpack())

}

type PackedString string

func (s PackedString) Unpack() string {
	var lastRune, lastLetter rune
	var result, num strings.Builder
	var esc bool
	result.Reset()
	num.Reset()
	lastRune = 0
	lastLetter = 0
	for i, curRune := range s {

		if unicode.IsDigit(curRune) && i == 0 {
			return ""
		}

		if unicode.IsLetter(curRune) {

			if unicode.IsDigit(lastRune) {
				numRunes, err := strconv.Atoi(num.String())
				if err != nil {
					log.Fatal(err)
				}
				for j := 0; j < numRunes-1; j++ {
					result.WriteRune(lastLetter)
				}
				num.Reset()
			}

			result.WriteRune(curRune)
			lastLetter = curRune
			lastRune = curRune
		}

		if unicode.IsDigit(curRune) {

			if esc {
				result.WriteRune(curRune)
				lastLetter = curRune
				lastRune = curRune
				esc = false
			} else {

				if unicode.IsLetter(lastRune) {
					num.Reset()
				}
				num.WriteRune(curRune)
				lastRune = curRune

				if i == utf8.RuneCountInString(string(s))-1 {
					numRunes, err := strconv.Atoi(num.String())
					if err != nil {
						log.Fatal(err)
					}
					for j := 0; j < numRunes-1; j++ {
						result.WriteRune(lastLetter)
					}
				}
			}

		}
		if curRune == '\\' {
			if lastRune == '\\' {
				result.WriteRune(curRune)
				lastLetter = curRune
				lastRune = curRune
				esc = false

			} else {
				esc = true
				lastRune = curRune
			}
		}
	}

	return result.String()
}
