package main

import (
	"fmt"
	"slices"
	"strings"
)

// функция для поиска анаграмм
func searchAnagrams(input *[]string) *map[string]*[]string {
	anagrams := make(map[string]*[]string) // создание предварительной map для сортировки анаграмм
	for _, word := range *input {
		anagram := []rune(strings.ToLower(word)) // приведение слова к строчным буквам и конвертация в слайс рун
		slices.Sort(anagram)                     // сортировка слайса рун
		sortedWord := string(anagram)            // конвертация слайса обратно в строку
		// проверка на существование анаграммы в map
		if _, exists := anagrams[sortedWord]; !exists {
			newWordSlice := []string{word}       // создание нового слайса слов
			anagrams[sortedWord] = &newWordSlice // добавление слайса в map
		} else {
			*anagrams[sortedWord] = append(*anagrams[sortedWord], word) // добавление слова в слайс внутри map
		}
	}
	result := make(map[string]*[]string) // создание результирующей map
	for _, val := range anagrams {
		slices.Sort(*val)           // сортировка слайса слов внутри map
		*val = slices.Compact(*val) // очистка повторяющихся слов
		if len(*val) != 1 {
			result[(*val)[0]] = val // добавление слайса в результирующую map, если слайс длиннее одного элемента
		}
	}
	return &result // возврат ссылки на результирующую map
}

func main() {
	// слайс для проверки функции
	words := []string{
		"пятак",
		"пятак",
		"тяпка",
		"слиток",
		"термика",
		"столик",
		"абоба",
		"пятка",
		"материк",
		"метрика",
		"абоба",
		"листок",
		"мошкара",
		"ромашка",
		"слиток",
	}
	res := searchAnagrams(&words) // вызов функции
	// вывод результатов
	for k, v := range *res {
		fmt.Println(k, *v)
	}
}
