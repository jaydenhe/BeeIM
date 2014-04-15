package s2c

import (
	"log"
	"net"
)

func Run() {
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

	ch := make(chan Frame,2)

	go StartAgent(ch, conn)

	for {
		frame,err:=ReadAFrame(conn)

		if err != nil {
			break
		}
		ch <- *frame
	}
	close(ch)
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error:%v", err)
	}
}











