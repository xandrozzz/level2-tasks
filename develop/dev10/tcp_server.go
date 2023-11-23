package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:4040")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn)

	buf := make([]byte, 32)
	for {
		_, err := conn.Write([]byte("data from server, awaiting response from client\n"))
		if err != nil {
			log.Fatalln(err)
		}

		readLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		clientData := string(buf)
		fmt.Println(clientData)

		_, err = conn.Write(append([]byte("got from client:"), buf[:readLen]...))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
