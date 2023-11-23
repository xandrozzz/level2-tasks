package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
)

func main() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org") // получение точного времени с ntp сервера
	// обработка ошибки
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Current time:", time.String()) // вывод результата
}
