package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// lineSeparator - класс разделителя строк
type lineSeparator struct {
	fields        []int
	delimiter     string
	separatedOnly bool
	lines         []string
}

// newLineSeparator - конструктор класса lineSeparator
func newLineSeparator(args []string, lines []string) *lineSeparator {

	argFieldsSplit := make([]int, 0) // создание массива для номеров выводимых полей

	// получение количества полных столбцов
	guaranteedColumns := len(strings.Fields(lines[0]))

	for _, line := range lines {
		colsCount := len(strings.Fields(line))
		if colsCount < guaranteedColumns {
			guaranteedColumns = colsCount
		}
	}

	// проверка флага -f
	if !slices.Contains(args, "-f") {
		// если флага нет, выводить все поля
		for i := 0; i < guaranteedColumns; i++ {
			argFieldsSplit = append(argFieldsSplit, i)
		}
	} else {
		// если флаг есть, получение аргумента после него
		argIndex := slices.Index(args, "-f")
		if len(args) > argIndex+1 {
			argFieldsString := args[argIndex+1]
			// если столбцы написаны через запятую
			if strings.Contains(argFieldsString, ",") {
				// добавление каждого столбца по отдельности
				for _, stringField := range strings.Split(argFieldsString, ",") {
					intField, err := strconv.Atoi(stringField)
					if err != nil {
						log.Fatalln(err)
					}
					argFieldsSplit = append(argFieldsSplit, intField)
				}
				// если столбцы написаны через дефис
			} else if strings.Contains(argFieldsString, "-") {
				splitFields := strings.Split(argFieldsString, "-")
				// проверка аргументов, написанных через дефис
				if len(splitFields) == 2 {
					// если оба аргумента присутствуют, вывод столбцов между ними
					if len(splitFields[0]) > 0 && len(splitFields[1]) > 0 {
						startField, err := strconv.Atoi(splitFields[0])
						if err != nil {
							log.Fatalln(err)
						}
						endField, err := strconv.Atoi(splitFields[1])
						if err != nil {
							log.Fatalln(err)
						}
						if startField > endField {
							log.Fatalln("invalid arguments given")
						}
						for i := startField - 1; i < endField; i++ {
							argFieldsSplit = append(argFieldsSplit, i)
						}
						// если только начальный аргумент присутствует, вывод столбцов после него
					} else if len(splitFields[0]) > 0 {
						startField, err := strconv.Atoi(splitFields[0])
						if err != nil {
							log.Fatalln(err)
						}
						endField := guaranteedColumns
						if startField > endField {
							log.Fatalln("invalid arguments given")
						}
						for i := startField - 1; i < endField; i++ {
							argFieldsSplit = append(argFieldsSplit, i)
						}
						// если только конечный аргумент присутствует, вывод столбцов после него
					} else if len(splitFields[1]) > 0 {
						startField := 1
						endField, err := strconv.Atoi(splitFields[1])
						if err != nil {
							log.Fatalln(err)
						}
						if startField > endField {
							log.Fatalln("invalid arguments given")
						}
						for i := startField - 1; i < endField; i++ {
							argFieldsSplit = append(argFieldsSplit, i)
						}
					} else {
						log.Fatalln("invalid arguments given")
					}
				} else {
					log.Fatalln("invalid arguments given")
				}
				// если запятых и дефисов нет, вывод конкретного столбца
			} else {
				singleField, err := strconv.Atoi(argFieldsString)
				if err != nil {
					log.Fatalln("invalid arguments given")
				}
				argFieldsSplit = append(argFieldsSplit, singleField)
			}

		} else {
			log.Fatalln("invalid arguments given")
		}
	}

	var sep string

	// проверка флага -d
	if slices.Contains(args, "-d") {
		// если флаг присутствует, установка указанного разделителя
		sepIndex := slices.Index(args, "-d")
		if len(args) > sepIndex+1 {
			// получение количества строк контекста
			sep = args[sepIndex+1]
		} else {
			log.Fatalln("invalid arguments given")
		}
		// по умолчанию установка разделителя на tab
	} else {
		sep = "\t"
	}

	// проверка флага -s
	separatedOnly := slices.Contains(args, "-s")

	// возврат ссылки на созданный объект
	return &lineSeparator{
		lines:         lines,
		fields:        argFieldsSplit,
		delimiter:     sep,
		separatedOnly: separatedOnly,
	}

}

// separateLines - метод класса lineSeparator для разделения строк на поля
func (l *lineSeparator) separateLines() *[][]string {
	separatedLines := make([][]string, 0) // создание слайса для разделенных строк
	// цикл для разделения строк
	for _, line := range l.lines {
		// если поле separatedOnly установлено на true и в строке нет разделителя, пропуск этой строки
		if l.separatedOnly {
			if !strings.Contains(line, l.delimiter) {
				continue
			}
		}
		// разделение строки
		separatedLine := make([]string, 0)
		splitLine := strings.Split(line, l.delimiter)
		for _, field := range l.fields {
			separatedLine = append(separatedLine, splitLine[field])
		}
		// добавление разделенной строки в слайс
		separatedLines = append(separatedLines, separatedLine)
	}
	return &separatedLines // возврат ссылки на слайс
}

func main() {
	filename := os.Args[1] // получение пути к файлу для обработки

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

	separator := newLineSeparator(os.Args[2:], lines) // создание объекта lineSeparator

	separatedLines := separator.separateLines() // вызов метода separateLines для разделения строк

	// вывод результатов
	for _, line := range *separatedLines {
		fmt.Println(strings.Join(line, "\t"))
	}
}
