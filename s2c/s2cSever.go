package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func main() {
	log.Println("Starting the server")
	//listen
	service := ":1114"

	log.Println("Service:", service)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Println("s2cServer OK!")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		log.Println("accept :", conn.RemoteAddr())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 10)

	go StartAgent(ch, conn)

	for {
		n, err := io.ReadFull(conn, header)

		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receiving header:", err)
			break
		}

		//read data
		size := binary.BigEndian.Uint16(header)
		data := make([]byte, size)

		n, err := io.ReadFull(conn.data)

		if err != nil {
			log.Println("error receiving msg :", err)
			break
		}
		ch <- data
	}
	close(ch)
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error:%v", err)
	}
}
