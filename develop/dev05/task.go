package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[2] // получение пути к файлу для обработки

	// открытие выбранного файла для чтения и записи
	file, err := os.OpenFile(filename, os.O_RDWR, os.FileMode(0755))
	if err != nil {
		log.Fatalln("file opening error, file:", filename, err)
	}

	// закрытие файла в defer, чтобы избежать утечки
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file) // создание сканера для чтения файла

	scanner.Split(bufio.ScanLines) // установка функции разделения для чтения по строкам

	lines := make([]string, 0) // создание слайса строк

	// чтение строк из файла
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// проверка опциональных флагов
	count := slices.Contains(os.Args[3:], "-c")

	ignoreCase := slices.Contains(os.Args[3:], "-i")

	invert := slices.Contains(os.Args[3:], "-v")

	fixed := slices.Contains(os.Args[3:], "-F")

	numbers := slices.Contains(os.Args[3:], "-n")

	// получение строки для поиска
	searchQuery := os.Args[1]
	// проверка флага -i
	if fixed {
		searchQuery = strings.ToLower(searchQuery) // приведение строки для поиска к нижнему регистру
	}

	re, err := regexp.Compile(searchQuery) // создание регулярного выражения для поиска строк
	if err != nil {
		log.Fatalln(err)
	}

	// переключение по основным флагам
	switch {
	// проверка на флаг -C
	case slices.Contains(os.Args[3:], "-C"):
		linesToPrint := make([]int, 0)          // создание слайса для номеров строк файла, которые нужно вывести
		lineStrings := make(map[int]string)     // создание map для хранения содержимого строк
		argIndex := slices.Index(os.Args, "-C") // получение позиции флага -C
		// проверка количества аргументов
		if len(os.Args) > argIndex+1 {
			// получение количества строк контекста
			argParam := os.Args[argIndex+1]
			linesContext, err := strconv.Atoi(argParam)
			if err != nil {
				log.Fatalln(err)
			}
			// если запрос инвертирован, заполнение слайса и map всеми строками файла
			if invert {
				for i := 0; i < len(lines); i++ {
					linesToPrint = append(linesToPrint, i)
					lineStrings[i] = lines[i]
				}
			}
			// цикл для поиска по строкам
			for lineNum, line := range lines {
				// если был установлен флаг -i, приведение строки к нижнему регистру
				if ignoreCase {
					line = strings.ToLower(line)
				}
				var matches [][]byte
				// переключение проверки строки относительно флага -F
				if fixed {
					// проверка совпадения строки с запросом поиска
					if line == searchQuery {
						matches = [][]byte{[]byte("yes")}
					}
				} else {
					matches = re.FindAll([]byte(line), -1) // проверка совпадения строки с регулярным выражением
				}
				// если строка подходит под условие
				if len(matches) > 0 {
					// если был введен флаг -i, постепенно удаляются подходящие строки из слайса и map
					if invert {
						// цикл для всех строк в контексте
						for num := max(0, lineNum-linesContext); num < min(len(lines)-1, lineNum+linesContext)+1; num++ {
							delete(lineStrings, num) // удаление из map
							// удаление из слайса
							index := slices.Index(linesToPrint, num)
							if index != -1 {
								linesToPrint = slices.Delete(linesToPrint, index, index+1)
							}
						}
						// если флага -i не было, заполнения слайса и map подходящими строками
					} else {
						// цикл для всех строк в контексте
						for num := max(0, lineNum-linesContext); num < min(len(lines)-1, lineNum+linesContext)+1; num++ {
							linesToPrint = append(linesToPrint, num) // добавление номера строки в слайс
							lineStrings[num] = lines[num]            // добавление строки в map
						}
					}
				}

			}
			slices.Sort(linesToPrint)                   // сортировка слайса с номерами строк
			linesToPrint = slices.Compact(linesToPrint) // удаление повторяющихся номеров строк
			if count {
				fmt.Println(len(linesToPrint)) // если был введен флаг -c, вывод количества подходящих строк
			} else {
				// если не был введен флаг -c вывод подходящих строк
				for num, line := range linesToPrint {
					if numbers {
						fmt.Print(strconv.Itoa(num) + ":") // если был введен флаг -n, вывод номера строки перед самой строкой
					}
					fmt.Println(lines[line])
				}
			}
		} else {
			log.Fatalln("incorrect parameters given")
		}
	// проверка на флаг -A
	case slices.Contains(os.Args[3:], "-A"):
		linesToPrint := make([]int, 0)          // создание слайса для номеров строк файла, которые нужно вывести
		lineStrings := make(map[int]string)     // создание map для хранения содержимого строк
		argIndex := slices.Index(os.Args, "-C") // получение позиции флага -C
		// проверка количества аргументов
		if len(os.Args) > argIndex+1 {
			// получение количества строк контекста
			argParam := os.Args[argIndex+1]
			linesContext, err := strconv.Atoi(argParam)
			if err != nil {
				log.Fatalln(err)
			}
			// если запрос инвертирован, заполнение слайса и map всеми строками файла
			if invert {
				for i := 0; i < len(lines); i++ {
					linesToPrint = append(linesToPrint, i)
					lineStrings[i] = lines[i]
				}
			}
			// цикл для поиска по строкам
			for lineNum, line := range lines {
				// если был установлен флаг -i, приведение строки к нижнему регистру
				if ignoreCase {
					line = strings.ToLower(line)
				}
				var matches [][]byte
				// переключение проверки строки относительно флага -F
				if fixed {
					// проверка совпадения строки с запросом поиска
					if line == searchQuery {
						matches = [][]byte{[]byte("yes")}
					}
				} else {
					matches = re.FindAll([]byte(line), -1) // проверка совпадения строки с регулярным выражением
				}
				// если строка подходит под условие
				if len(matches) > 0 {
					// если был введен флаг -i, постепенно удаляются подходящие строки из слайса и map
					if invert {
						// цикл для всех строк после найденной, включая ее саму
						for num := lineNum; num < min(len(lines)-1, lineNum+linesContext)+1; num++ {
							delete(lineStrings, num) // удаление из map
							// удаление из слайса
							index := slices.Index(linesToPrint, num)
							if index != -1 {
								linesToPrint = slices.Delete(linesToPrint, index, index+1)
							}
						}
						// если флага -i не было, заполнения слайса и map подходящими строками
					} else {
						// цикл для всех строк после найденной, включая ее саму
						for num := lineNum; num < min(len(lines)-1, lineNum+linesContext)+1; num++ {
							linesToPrint = append(linesToPrint, num) // добавление номера строки в слайс
							lineStrings[num] = lines[num]            // добавление строки в map
						}
					}
				}

			}
			slices.Sort(linesToPrint)                   // сортировка слайса с номерами строк
			linesToPrint = slices.Compact(linesToPrint) // удаление повторяющихся номеров строк
			if count {
				fmt.Println(len(linesToPrint)) // если был введен флаг -c, вывод количества подходящих строк
			} else {
				// если не был введен флаг -c вывод подходящих строк
				for num, line := range linesToPrint {
					if numbers {
						fmt.Print(strconv.Itoa(num) + ":") // если был введен флаг -n, вывод номера строки перед самой строкой
					}
					fmt.Println(lines[line])
				}
			}
		} else {
			log.Fatalln("incorrect parameters given")
		}
	// проверка на флаг -B
	case slices.Contains(os.Args[3:], "-B"):
		linesToPrint := make([]int, 0)          // создание слайса для номеров строк файла, которые нужно вывести
		lineStrings := make(map[int]string)     // создание map для хранения содержимого строк
		argIndex := slices.Index(os.Args, "-C") // получение позиции флага -C
		// проверка количества аргументов
		if len(os.Args) > argIndex+1 {
			// получение количества строк контекста
			argParam := os.Args[argIndex+1]
			linesContext, err := strconv.Atoi(argParam)
			if err != nil {
				log.Fatalln(err)
			}
			// если запрос инвертирован, заполнение слайса и map всеми строками файла
			if invert {
				for i := 0; i < len(lines); i++ {
					linesToPrint = append(linesToPrint, i)
					lineStrings[i] = lines[i]
				}
			}
			// цикл для поиска по строкам
			for lineNum, line := range lines {
				// если был установлен флаг -i, приведение строки к нижнему регистру
				if ignoreCase {
					line = strings.ToLower(line)
				}
				var matches [][]byte
				// переключение проверки строки относительно флага -F
				if fixed {
					// проверка совпадения строки с запросом поиска
					if line == searchQuery {
						matches = [][]byte{[]byte("yes")}
					}
				} else {
					matches = re.FindAll([]byte(line), -1) // проверка совпадения строки с регулярным выражением
				}
				// если строка подходит под условие
				if len(matches) > 0 {
					// если был введен флаг -i, постепенно удаляются подходящие строки из слайса и map
					if invert {
						// цикл для всех строк перед найденной, включая ее саму
						for num := max(0, lineNum-linesContext); num < lineNum+1; num++ {
							delete(lineStrings, num) // удаление из map
							// удаление из слайса
							index := slices.Index(linesToPrint, num)
							if index != -1 {
								linesToPrint = slices.Delete(linesToPrint, index, index+1)
							}
						}
						// если флага -i не было, заполнения слайса и map подходящими строками
					} else {
						// цикл для всех строк перед найденной, включая ее саму
						for num := max(0, lineNum-linesContext); num < lineNum+1; num++ {
							linesToPrint = append(linesToPrint, num) // добавление номера строки в слайс
							lineStrings[num] = lines[num]            // добавление строки в map
						}
					}
				}

			}
			slices.Sort(linesToPrint)                   // сортировка слайса с номерами строк
			linesToPrint = slices.Compact(linesToPrint) // удаление повторяющихся номеров строк
			if count {
				fmt.Println(len(linesToPrint)) // если был введен флаг -c, вывод количества подходящих строк
			} else {
				// если не был введен флаг -c вывод подходящих строк
				for num, line := range linesToPrint {
					if numbers {
						fmt.Print(strconv.Itoa(num) + ":") // если был введен флаг -n, вывод номера строки перед самой строкой
					}
					fmt.Println(lines[line])
				}
			}
		} else {
			log.Fatalln("incorrect parameters given")
		}
	// стандартный случай, если флагов -A -B -C нет
	default:
		linesToPrint := make([]int, 0)      // создание слайса для номеров строк файла, которые нужно вывести
		lineStrings := make(map[int]string) // создание map для хранения содержимого строк
		// если запрос инвертирован, заполнение слайса и map всеми строками файла
		if invert {
			for i := 0; i < len(lines); i++ {
				linesToPrint = append(linesToPrint, i)
				lineStrings[i] = lines[i]
			}
		}
		// цикл для поиска по строкам
		for lineNum, line := range lines {
			// если был установлен флаг -i, приведение строки к нижнему регистру
			if ignoreCase {
				line = strings.ToLower(line)
			}
			var matches [][]byte
			// переключение проверки строки относительно флага -F
			if fixed {
				// проверка совпадения строки с запросом поиска
				if line == searchQuery {
					matches = [][]byte{[]byte("yes")}
				}
			} else {
				matches = re.FindAll([]byte(line), -1) // проверка совпадения строки с регулярным выражением
			}
			// если строка подходит под условие
			if len(matches) > 0 {
				// если был введен флаг -i, постепенно удаляются подходящие строки из слайса и map
				if invert {
					delete(lineStrings, lineNum) // удаление из map
					// удаление из слайса
					index := slices.Index(linesToPrint, lineNum)
					if index != -1 {
						linesToPrint = slices.Delete(linesToPrint, index, index+1)
					}
					// если флага -i не было, заполнения слайса и map подходящими строками
				} else {
					linesToPrint = append(linesToPrint, lineNum) // добавление номера строки в слайс
					lineStrings[lineNum] = lines[lineNum]        // добавление строки в map
				}
			}

		}
		slices.Sort(linesToPrint)                   // сортировка слайса с номерами строк
		linesToPrint = slices.Compact(linesToPrint) // удаление повторяющихся номеров строк
		if count {
			fmt.Println(len(linesToPrint)) // если был введен флаг -c, вывод количества подходящих строк
		} else {
			// если не был введен флаг -c вывод подходящих строк
			for num, line := range linesToPrint {
				if numbers {
					fmt.Print(strconv.Itoa(num) + ":") // если был введен флаг -n, вывод номера строки перед самой строкой
				}
				fmt.Println(lines[line])
			}
		}
	}

}
