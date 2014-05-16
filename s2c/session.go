package s2c

import (
	"bufio"
	"log"
	"net"
)

type Session struct {
	conn     net.Conn
	incoming chan Frame
	outgoing chan Frame
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan net.Conn
	name     string
}

func (self *Session) GetName() string {
	return self.name
}

func (self *Session) SetName(name string) {
	self.name = name
}

func (self *Session) GetIncoming() string {
	return <-self.incoming
}

func (self *Session) PutOutgoing(message string) {
	self.outgoing <- message
}

func CreateSession(conn net.Conn) *Session {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	Session := &Session{
		conn:     conn,
		incoming: make(chan Frame),
		outgoing: make(chan Frame),
		quiting:  make(chan net.Conn),
		reader:   reader,
		writer:   writer,
	}
	Session.Listen()
	return Session
}

func (self *Session) Listen() {
	go self.Read()
	go self.Write()
}

func (self *Session) quit() {
	self.quiting <- self.conn
}

func (self *Session) Read() {
	for {
		if line, _, err := self.reader.ReadLine(); err == nil {
			self.incoming <- string(line)
		} else {
			log.Printf("Read error: %s\n", err)
			self.quit()
			return
		}
	}

}

func (self *Session) Write() {
	for data := range self.outgoing {
		if _, err := self.writer.WriteString(data + "\n"); err != nil {
			self.quit()
			return
		}

		if err := self.writer.Flush(); err != nil {
			log.Printf("Write error: %s\n", err)
			self.quit()
			return
		}
	}

}

func (self *Session) Close() {
	self.conn.Close()
}
