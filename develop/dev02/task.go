package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"
)

// ErrInvalidString is a custom error for strings that cannot be unpacked
var ErrInvalidString = errors.New("invalid string")

// Unpack is a function for unpacking a string with escape sequences
func Unpack(s string) (string, error) {
	sr := []rune(s)
	var s2 string
	var n int
	var backslash bool

	for i, item := range sr {
		if unicode.IsDigit(item) && i == 0 {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(item) && unicode.IsDigit(sr[i-1]) && sr[i-2] != '\\' {
			return "", ErrInvalidString
		}
		if item == '\\' && !backslash {
			backslash = true
			continue
		}
		if backslash && unicode.IsLetter(item) {
			return "", ErrInvalidString
		}
		if backslash {
			s2 += string(item)
			backslash = false
			continue
		}
		if unicode.IsDigit(item) {
			n = int(item - '0')
			if n == 0 {
				s2 = s2[:len(s2)-1]
				continue
			}
			for j := 0; j < n-1; j++ {
				s2 += string(sr[i-1])
			}
			continue
		}
		s2 += string(item)
	}

	return s2, nil
}

func main() {
	strToUnpack := "45"
	result, err := Unpack(strToUnpack)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	os.Exit(0)
}
