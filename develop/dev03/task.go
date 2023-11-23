package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// iSorter - интерфейс для сортировщика
type iSorter interface {
	sortSlice(*[]string)
	checkSlice(*[]string) bool
}

// sorter - основной класс сортировщика
type sorter struct {
	reverse              bool
	unique               bool
	removeTrailingSpaces bool
	higherColumn         int
	lowerColumn          int
}

// sortSlice - реализация метода sortSlice интерфейса iSorter классом sorter
func (a *sorter) sortSlice(s *[]string) {
	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if a.removeTrailingSpaces {
		for _, line := range *s {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range *s {
		splitLines = append(splitLines, strings.Fields(line))
	}
	// сортировка строк в алфавитном порядке на основании выбранных полей
	sort.Slice(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := a.lowerColumn; col >= a.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})

	// сборка полей обратно в строки
	for key, splitLine := range splitLines {
		(*s)[key] = strings.Join(splitLine, " ")
	}

	// удаление повторяющихся строк на основании поля unique
	if a.unique {
		*s = slices.Compact(*s)
	}

	// разворот слайса на основании поля reverse
	if a.reverse {
		slices.Reverse(*s)
	}
}

// checkSlice - реализация метода checkSlice интерфейса iSorter классом sorter
func (a *sorter) checkSlice(s *[]string) bool {
	// копирование слайса в новый
	var sCopy []string
	copy(sCopy, *s)

	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if a.removeTrailingSpaces {
		for _, line := range sCopy {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range sCopy {
		splitLines = append(splitLines, strings.Fields(line))
	}
	// проверка сортировки строк в алфавитном порядке на основании выбранных полей
	return sort.SliceIsSorted(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := a.lowerColumn; col >= a.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})

}

// alphabeticalSorter - конкретный класс сортировщика
type alphabeticalSorter struct {
	sorter
}

// sortSlice - реализация метода sortSlice интерфейса iSorter классом alphabeticalSorter
func (a *alphabeticalSorter) sortSlice(s *[]string) {
	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if a.removeTrailingSpaces {
		for _, line := range *s {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range *s {
		splitLines = append(splitLines, strings.Fields(line))
	}
	// сортировка строк в алфавитном порядке на основании выбранных полей
	sort.Slice(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := a.lowerColumn; col >= a.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})

	// сборка полей обратно в строки
	for key, splitLine := range splitLines {
		(*s)[key] = strings.Join(splitLine, " ")
	}

	// удаление повторяющихся строк на основании поля unique
	if a.unique {
		*s = slices.Compact(*s)
	}

	// разворот слайса на основании поля reverse
	if a.reverse {
		slices.Reverse(*s)
	}

}

// checkSlice - реализация метода checkSlice интерфейса iSorter классом alphabeticalSorter
func (a *alphabeticalSorter) checkSlice(s *[]string) bool {
	// копирование слайса в новый
	var sCopy []string
	copy(sCopy, *s)

	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if a.removeTrailingSpaces {
		for _, line := range sCopy {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range sCopy {
		splitLines = append(splitLines, strings.Fields(line))
	}
	// проверка сортировки строк в алфавитном порядке на основании выбранных полей
	return sort.SliceIsSorted(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := a.lowerColumn; col >= a.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})

}

// numericalSorter - конкретный класс сортировщика
type numericalSorter struct {
	sorter
}

// sortSlice - реализация метода sortSlice интерфейса iSorter классом numericalSorter
func (n *numericalSorter) sortSlice(s *[]string) {
	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if n.removeTrailingSpaces {
		for _, line := range *s {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля c конвертацией строк в числа
	splitLines := make([][]int, 0)
	for _, line := range *s {
		splitLine := strings.Fields(line)
		intSplitLine := make([]int, 0)
		for _, word := range splitLine {
			intWord, err := strconv.Atoi(word)
			if err != nil {
				log.Fatalln("file contains non-numerical values")
			}
			intSplitLine = append(intSplitLine, intWord)
		}
		splitLines = append(splitLines, intSplitLine)
	}
	// сортировка строк в числовом порядке на основании выбранных полей
	sort.Slice(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := n.lowerColumn; col >= n.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})
	// сборка полей с конвертацией типа обратно в строки
	for key, splitLine := range splitLines {
		stringLine := make([]string, 0)
		for _, intWord := range splitLine {
			word := strconv.Itoa(intWord)
			stringLine = append(stringLine, word)
		}
		(*s)[key] = strings.Join(stringLine, " ")
	}

	// удаление повторяющихся строк на основании поля unique
	if n.unique {
		*s = slices.Compact(*s)
	}

	// разворот слайса на основании поля reverse
	if n.reverse {
		slices.Reverse(*s)
	}
}

// checkSlice - реализация метода checkSlice интерфейса iSorter классом numericalSorter
func (n *numericalSorter) checkSlice(s *[]string) bool {
	// копирование слайса в новый
	var sCopy []string
	copy(sCopy, *s)

	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if n.removeTrailingSpaces {
		for _, line := range sCopy {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля c конвертацией строк в числа
	splitLines := make([][]int, 0)
	for _, line := range sCopy {
		splitLine := strings.Fields(line)
		intSplitLine := make([]int, 0)
		for _, word := range splitLine {
			intWord, err := strconv.Atoi(word)
			if err != nil {
				log.Fatalln("file contains non-numerical values")
			}
			intSplitLine = append(intSplitLine, intWord)
		}
		splitLines = append(splitLines, intSplitLine)
	}
	// проверка сортировки строк в числовом порядке на основании выбранных полей
	return sort.SliceIsSorted(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := n.lowerColumn; col >= n.higherColumn; col++ {
			lastCol = col
			if splitLines[i][col] != splitLines[j][col] {
				break
			}
		}
		return splitLines[i][lastCol] < splitLines[j][lastCol]
	})
}

// monthSorter - конкретный класс сортировщика
type monthSorter struct {
	sorter
}

// sortSlice - реализация метода sortSlice интерфейса iSorter классом monthSorter
func (m *monthSorter) sortSlice(s *[]string) {
	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if m.removeTrailingSpaces {
		for _, line := range *s {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля c конвертацией строк в тип time.Month
	monthLines := make([][]time.Month, 0)
	for _, line := range *s {
		splitLine := strings.Fields(line)
		monthWords := make([]time.Month, 0)
		for _, word := range splitLine {
			dateWord, err := time.Parse("January", word)
			if err != nil {
				log.Fatalln("file contains non-month values", err)
			}
			monthWords = append(monthWords, dateWord.Month())
		}
		monthLines = append(monthLines, monthWords)
	}
	// сортировка строк в порядке месяцев на основании выбранных полей
	sort.Slice(monthLines, func(i int, j int) bool {
		lastCol := 0
		for col := m.lowerColumn; col >= m.higherColumn; col++ {
			lastCol = col
			if monthLines[i][col] != monthLines[j][col] {
				break
			}
		}
		return monthLines[i][lastCol] < monthLines[j][lastCol]
	})
	// сборка полей с конвертацией типа обратно в строки
	for key, line := range monthLines {
		for _, monthWord := range line {
			word := monthWord.String()
			(*s)[key] = word
		}
	}

	// удаление повторяющихся строк на основании поля unique
	if m.unique {
		*s = slices.Compact(*s)
	}

	// разворот слайса на основании поля reverse
	if m.reverse {
		slices.Reverse(*s)
	}
}

// checkSlice - реализация метода checkSlice интерфейса iSorter классом monthSorter
func (m *monthSorter) checkSlice(s *[]string) bool {
	// копирование слайса в новый
	var sCopy []string
	copy(sCopy, *s)

	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if m.removeTrailingSpaces {
		for _, line := range sCopy {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля c конвертацией строк в тип time.Month
	monthLines := make([][]time.Month, 0)
	for _, line := range sCopy {
		splitLine := strings.Fields(line)
		monthWords := make([]time.Month, 0)
		for _, word := range splitLine {
			dateWord, err := time.Parse("January", word)
			if err != nil {
				log.Fatalln("file contains non-month values", err)
			}
			monthWords = append(monthWords, dateWord.Month())
		}
		monthLines = append(monthLines, monthWords)
	}
	// проверка сортировки строк в порядке месяцев на основании выбранных полей
	return sort.SliceIsSorted(monthLines, func(i int, j int) bool {
		lastCol := 0
		for col := m.lowerColumn; col >= m.higherColumn; col++ {
			lastCol = col
			if monthLines[i][col] != monthLines[j][col] {
				break
			}
		}
		return monthLines[i][lastCol] < monthLines[j][lastCol]
	})
}

// humanNumericalSorter - конкретный класс сортировщика
type humanNumericalSorter struct {
	sorter
}

// sortSlice - реализация метода sortSlice интерфейса iSorter классом humanNumericalSorter
func (h *humanNumericalSorter) sortSlice(s *[]string) {
	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if h.removeTrailingSpaces {
		for _, line := range *s {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range *s {
		splitLines = append(splitLines, strings.Split(line, " "))
	}
	re, err := regexp.Compile("[0-9]+") // создание regexp для поиска лидирующих чисел
	if err != nil {
		log.Fatalln("incorrect regexp", err)
	}
	// сортировка строк в числовом порядке относительно лидирующих чисел
	sort.Slice(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := h.lowerColumn; col >= h.higherColumn; col++ {
			lastCol = col
			// сравнение лидирующих чисел через regexp
			if string(re.Find([]byte(splitLines[i][col]))) != string(re.Find([]byte(splitLines[j][col]))) {
				break
			}
		}
		// конвертация лидирующих чисел в int
		val1, err := strconv.Atoi(string(re.Find([]byte(splitLines[i][lastCol]))))
		if err != nil {
			log.Fatalln(err)
		}
		val2, err := strconv.Atoi(string(re.Find([]byte(splitLines[j][lastCol]))))
		if err != nil {
			log.Fatalln(err)
		}
		// возврат результата сравнения
		return val1 < val2
	})

	for key, splitLine := range splitLines {
		(*s)[key] = strings.Join(splitLine, " ")
	}

	if h.unique {
		*s = slices.Compact(*s)
	}
	if h.reverse {
		slices.Reverse(*s)
	}
}

// checkSlice - реализация метода checkSlice интерфейса iSorter классом humanNumericalSorter
func (h *humanNumericalSorter) checkSlice(s *[]string) bool {
	// копирование слайса в новый
	var sCopy []string
	copy(sCopy, *s)

	// удаление хвостовых пробелов на основании поля removeTrailingSpaces
	if h.removeTrailingSpaces {
		for _, line := range sCopy {
			strings.TrimRight(line, " \n\r\t")
		}
	}
	// разделение строк на поля
	splitLines := make([][]string, 0)
	for _, line := range sCopy {
		splitLines = append(splitLines, strings.Fields(line))
	}
	re, err := regexp.Compile("[0-9]+") // создание regexp для поиска лидирующих чисел
	if err != nil {
		log.Fatalln("incorrect regexp", err)
	}
	// проверка сортировки строк в числовом порядке относительно лидирующих чисел
	return sort.SliceIsSorted(splitLines, func(i int, j int) bool {
		lastCol := 0
		for col := h.lowerColumn; col >= h.higherColumn; col++ {
			lastCol = col
			// сравнение лидирующих чисел через regexp
			if string(re.Find([]byte(splitLines[i][col]))) != string(re.Find([]byte(splitLines[j][col]))) {
				break
			}
		}
		// конвертация лидирующих чисел в int
		val1, err := strconv.Atoi(string(re.Find([]byte(splitLines[i][lastCol]))))
		if err != nil {
			log.Fatalln(err)
		}
		val2, err := strconv.Atoi(string(re.Find([]byte(splitLines[j][lastCol]))))
		if err != nil {
			log.Fatalln(err)
		}
		// возврат результата сравнения
		return val1 < val2
	})
}

// sortContext - класс контекст для сортировщиков
type sortContext struct {
	sort iSorter
}

// sortLines - метод для вызова внутреннего метода сортировщика sortLines
func (c *sortContext) sortLines(s *[]string) {
	c.sort.sortSlice(s)
}

// checkSortedLines - метод для вызова внутреннего метода сортировщика checkSortedLines
func (c *sortContext) checkSortedLines(s *[]string) bool {
	return c.sort.checkSlice(s)
}

// getColumns - функция для получения номеров столбцов из командной строки
func getColumns(argIndex, guaranteedColumns int) (lowerColumn, higherColumn int) {
	// получение аргумента, следующего за ключом
	cols := os.Args[argIndex+1]
	// проверка формата аргумента
	if strings.Count(cols, ",") != 1 {
		log.Fatalln("incorrect column indexes given")
	}
	// получение номеров столбцов
	splitCols := strings.Split(cols, ",")
	lowerCol, err := strconv.Atoi(splitCols[0])
	if err != nil {
		log.Fatalln("non-numerical colon index given:", err)
	}
	lowerCol-- // уменьшение номера первого столбца на 1 для получения корректного индекса
	higherCol, err := strconv.Atoi(splitCols[0])
	if err != nil {
		log.Fatalln("non-numerical colon index given:", err)
	}
	// сравнение номеров столбцов с количеством полных столбцов
	if higherCol > guaranteedColumns || lowerCol > guaranteedColumns {
		log.Fatalln("given column indexes exceed number of complete columns")
	}
	// сравнение номеров начального и конечного столбца
	if lowerCol >= higherCol {
		log.Fatalln("incorrect column indexes given")
	}
	return lowerCol, higherCol
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

	// проверка на ключ -r
	isReverse := false

	if slices.Contains(os.Args[2:], "-r") {
		isReverse = true
	}

	// проверка на ключ -u
	isUnique := false

	if slices.Contains(os.Args[2:], "-u") {
		isUnique = true
	}

	// проверка на ключ -b
	ignoreTrailingSpaces := false

	if slices.Contains(os.Args[2:], "-b") {
		ignoreTrailingSpaces = true
	}

	// получение количества полных столбцов
	guaranteedColumns := len(strings.Fields(lines[0]))

	for _, line := range lines {
		colsCount := len(strings.Fields(line))
		if colsCount < guaranteedColumns {
			guaranteedColumns = colsCount
		}
	}

	var sliceSorter *sortContext // объявление контекста сортировщика

	switch {
	// проверка на ключ -n
	case slices.Contains(os.Args[2:], "-n"):
		// проверка на ключ -k
		if slices.Contains(os.Args, "-k") {
			argIndex := slices.Index(os.Args, "-k")
			if len(os.Args) > argIndex+1 {
				lowerCol, higherCol := getColumns(argIndex, guaranteedColumns) // получение номеров столбцов

				// создание объекта сортировщика в числовом порядке с указанными столбцами
				sliceSorter = &sortContext{
					sort: &numericalSorter{
						sorter: sorter{
							unique:               isUnique,
							reverse:              isReverse,
							removeTrailingSpaces: ignoreTrailingSpaces,
							higherColumn:         higherCol,
							lowerColumn:          lowerCol,
						},
					},
				}

			} else {
				log.Fatalln("not enough arguments given after -k flag")
			}
		} else {
			// создание объекта сортировщика в числовом порядке для всех столбцов
			sliceSorter = &sortContext{
				sort: &numericalSorter{
					sorter: sorter{
						unique:               isUnique,
						reverse:              isReverse,
						removeTrailingSpaces: ignoreTrailingSpaces,
						higherColumn:         guaranteedColumns,
						lowerColumn:          0,
					},
				},
			}
		}
	// проверка на ключ -h
	case slices.Contains(os.Args[2:], "-h"):
		// проверка на ключ -k
		if slices.Contains(os.Args, "-k") {
			argIndex := slices.Index(os.Args, "-k")
			if len(os.Args) > argIndex+1 {
				lowerCol, higherCol := getColumns(argIndex, guaranteedColumns) // получение номеров столбцов

				// создание объекта сортировщика в числовом порядке на основании лидирующих чисел с указанными столбцами
				sliceSorter = &sortContext{
					sort: &humanNumericalSorter{
						sorter: sorter{
							unique:               isUnique,
							reverse:              isReverse,
							removeTrailingSpaces: ignoreTrailingSpaces,
							higherColumn:         higherCol,
							lowerColumn:          lowerCol,
						},
					},
				}

			} else {
				log.Fatalln("not enough arguments given after -k flag")
			}
		} else {
			// создание объекта сортировщика в числовом порядке на основании лидирующих чисел для всех столбцов
			sliceSorter = &sortContext{
				sort: &humanNumericalSorter{
					sorter: sorter{
						unique:               isUnique,
						reverse:              isReverse,
						removeTrailingSpaces: ignoreTrailingSpaces,
						higherColumn:         guaranteedColumns,
						lowerColumn:          0,
					},
				},
			}
		}
	// проверка на ключ -M
	case slices.Contains(os.Args[2:], "-M"):
		// проверка на ключ -k
		if slices.Contains(os.Args, "-k") {
			argIndex := slices.Index(os.Args, "-k")
			if len(os.Args) > argIndex+1 {
				lowerCol, higherCol := getColumns(argIndex, guaranteedColumns) // получение номеров столбцов

				// создание объекта сортировщика в порядке месяцев с указанными столбцами
				sliceSorter = &sortContext{
					sort: &monthSorter{
						sorter: sorter{
							unique:               isUnique,
							reverse:              isReverse,
							removeTrailingSpaces: ignoreTrailingSpaces,
							higherColumn:         higherCol,
							lowerColumn:          lowerCol,
						},
					},
				}

			} else {
				log.Fatalln("not enough arguments given after -k flag")
			}
		} else {
			// создание объекта сортировщика в порядке месяцев для всех столбцов
			sliceSorter = &sortContext{
				sort: &monthSorter{
					sorter: sorter{
						unique:               isUnique,
						reverse:              isReverse,
						removeTrailingSpaces: ignoreTrailingSpaces,
						higherColumn:         guaranteedColumns,
						lowerColumn:          0,
					},
				},
			}
		}
	default:
		// проверка на ключ -k
		if slices.Contains(os.Args, "-k") {
			argIndex := slices.Index(os.Args, "-k")
			if len(os.Args) > argIndex+1 {
				lowerCol, higherCol := getColumns(argIndex, guaranteedColumns) // получение номеров столбцов

				// создание объекта сортировщика в алфавитном порядке с указанными столбцами
				sliceSorter = &sortContext{
					sort: &alphabeticalSorter{
						sorter: sorter{
							unique:               isUnique,
							reverse:              isReverse,
							removeTrailingSpaces: ignoreTrailingSpaces,
							higherColumn:         higherCol,
							lowerColumn:          lowerCol,
						},
					},
				}

			} else {
				log.Fatalln("not enough arguments given after -k flag")
			}
		} else {
			// создание объекта сортировщика в алфавитном порядке для всех столбцов
			sliceSorter = &sortContext{
				sort: &alphabeticalSorter{
					sorter: sorter{
						unique:               isUnique,
						reverse:              isReverse,
						removeTrailingSpaces: ignoreTrailingSpaces,
						higherColumn:         guaranteedColumns,
						lowerColumn:          0,
					},
				},
			}
		}
	}

	// проверка на ключ -c
	if slices.Contains(os.Args[2:], "-c") {
		fmt.Println(sliceSorter.checkSortedLines(&lines)) // вызов метода для проверки сортировки
	} else {
		sliceSorter.sortLines(&lines) // вызов метода для сортировки
	}

	// очистка файла
	err = file.Truncate(0)
	if err != nil {
		log.Fatalln(err)
	}

	// установка оффсета файла на 0
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalln(err)
	}

	// запись строк в файл
	for _, line := range lines {
		_, err := file.WriteString(fmt.Sprintf("%s", line) + "\n")
		if err != nil {
			return
		}
	}
}
