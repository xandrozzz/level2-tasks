package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	// создание регулярного выражения для поиска чисел
	re, err := regexp.Compile("[0-9]+")
	if err != nil {
		return
	}

	// создание канала для отслеживания сигнала об остановке
	doneChannel := make(chan os.Signal, 1)
	signal.Notify(doneChannel, syscall.SIGTERM)

	if len(os.Args) >= 3 {
		// установка длительности таймаута по умолчанию
		timeoutLength := time.Duration(10) * time.Second
		// проверка ключа таймаута
		if slices.ContainsFunc(os.Args, containsTimeout) {
			// получение длины таймаута
			argIndex := slices.IndexFunc(os.Args, containsTimeout)
			if len(os.Args) == 4 {
				timeoutString := os.Args[argIndex]
				timeoutSplit := strings.Split(timeoutString, "=")
				timeoutLengthString := string(re.Find([]byte(timeoutSplit[1])))
				timeoutLengthInt, err := strconv.Atoi(timeoutLengthString)
				if err != nil {
					log.Fatalln(err)
				}
				// установка новой длины таймаута
				timeoutLength = time.Duration(timeoutLengthInt) * time.Second
			} else {
				log.Fatalln("invalid arguments given")
			}
		}

		serv := os.Args[1] + ":" + os.Args[2]                    // получение адреса и порта сервера
		conn, err := net.DialTimeout("tcp", serv, timeoutLength) // открытие TCP-подключения к серверу с таймаутом
		if err != nil {
			time.Sleep(timeoutLength)
			log.Fatalln(err)
		}
		go copyTo(os.Stdout, conn, doneChannel) // чтение из сокета в stdout
		copyTo(conn, os.Stdin, doneChannel)     // запись в сокет из stdin
	} else {
		log.Fatalln("invalid arguments given")
	}

}

// copyTo - функция для копирования данных
func copyTo(dst io.Writer, src io.Reader, doneChannel chan os.Signal) {
	// неблокирующая проверка на сигнал остановки
	select {
	case <-doneChannel:
		return
	default:
	}
	// копирование данных
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

// containsTimeout - функция для проверки нахождения строке таймаута
func containsTimeout(s string) bool {
	splitS := strings.Split(s, "=")
	if len(splitS) == 2 {
		if strings.Contains(splitS[0], "--timeout") {
			return true
		}
	}
	return false
}
