/*
author:jaydenhe
 */
package s2c

import (
	"bufio"
	"log"
	"net"
)

/*
SessionID 64bit

|uid (32bit)|extendId (32bit)|

SessionID = uid<<32 + extendId
 */
type TypeSessionID uint64

type Session struct {
	conn     net.Conn
	incoming chan Packet
	outgoing chan Packet
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan byte
	name     string
	id       TypeSessionID
}

func (self *Session) GetName() string {
	return self.name
}

func (self *Session) SetName(name string) {
	self.name = name
}

func (self *Session) GetID() TypeSessionID {
	return self.id
}

func (self *Session) SetID(id TypeSessionID) {
	self.id = id
}

func (self *Session) GetConn() net.Conn {
	return self.conn
}

func (self *Session) GetIncoming() Packet {
	return <-self.incoming
}

func (self *Session) PutOutgoing(packet Packet) {
	self.outgoing <- packet
}

func CreateSession(conn net.Conn) *Session {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	Session := &Session{
		conn:     conn,
		incoming: make(chan Packet),
		outgoing: make(chan Packet),
		quiting:  make(chan byte),
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
	self.quiting <- 0
}

func (self *Session) Read() {
	for {
		//		if line, _, err := self.reader.ReadLine(); err == nil {
		//			self.incoming <- string(line)
		//		} else {
		//			log.Printf("Read error: %s\n", err)
		//			self.quit()
		//			return
		//		}

		//		if packet,err :=
	}
	log.Println("Read()")

}

func (self *Session) Write() {
	/*	for data := range self.outgoing {
			if _, err := self.writer.WriteString(data + "\n"); err != nil {
				self.quit()
				return
			}

			if err := self.writer.Flush(); err != nil {
				log.Printf("Write error: %s\n", err)
				self.quit()
				return
			}
		}*/

}
