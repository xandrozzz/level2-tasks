package main

import (
	"fmt"
	"reflect"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	cases := make([]reflect.SelectCase, len(channels)) // создание слайса reflect.SelectCase с количеством кейсов, равным количеству каналов
	// добавление кейсов для всех каналов в слайс
	for i, ch := range channels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	reflect.Select(cases)                // вызов операции select
	newChannel := make(chan interface{}) // создание нового канала
	close(newChannel)                    // закрытие канала
	return newChannel                    // возврат закрытого канала
}

func main() {
	// создание функции для закрытия канала в горутине после определенного времени
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now() // получение времени начала отсчета
	// вызов функции or с блокировкой потока
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start)) // вывод времени, прошедшего с начала отсчета

}
