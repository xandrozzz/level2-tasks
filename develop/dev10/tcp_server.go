package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:4040") // открываем слушающий сокет
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			continue
		}
		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
	}
}

func handleClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn) // закрываем сокет при выходе из функции

	buf := make([]byte, 32) // буфер для чтения клиентских данных
	for {
		_, err := conn.Write([]byte("data from server, awaiting response from client\n"))
		if err != nil {
			log.Fatalln(err)
		} // пишем в сокет

		readLen, err := conn.Read(buf) // читаем из сокета
		if err != nil {
			fmt.Println(err)
			break
		}
		clientData := string(buf)
		fmt.Println(clientData)

		_, err = conn.Write(append([]byte("got from client:"), buf[:readLen]...))
		if err != nil {
			log.Fatalln(err)
		} // пишем в сокет
	}
}
