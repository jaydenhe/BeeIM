package s2c

import (
	"net"
	"log"
	"github.com/jaydenhe/BeeIM/uitl/uniq"
)

const (
	MAX_SESSIONS = 10000
)

type SessionTable map[TypeSessionID]*Session


type Server struct {
	listener net.Listener
	sessions SessionTable
	pending  chan net.Conn
	quiting  chan net.Conn
	incoming chan Packet
	outgoing chan Packet
	tokens   chan byte
}

func (self *Server) generateToken() {
	self.tokens <- 0
}

func (self *Server) takeToken() {
	<-self.tokens
}

func CreateServer() *Server {
	server := &Server{
		sessions:  make(SessionTable, MAX_SESSIONS),
		tokens:   make(chan byte, MAX_SESSIONS),
		pending:  make(chan net.Conn),
		quiting:  make(chan net.Conn),
		incoming: make(chan Packet),
		outgoing: make(chan Packet),
	}
	server.listen()
	return server
}

func (self *Server) Start(connString string) {

	listener, err := net.Listen("tcp", connString)

	if err != nil {
		log.Println(err)
	}

	self.listener = listener

	log.Printf("Server %p starts\n", self)

	// filling the tokens
	for i := 0; i < MAX_SESSIONS; i++ {
		self.generateToken()
	}

	for {
		conn, err := self.listener.Accept()

		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("A new connection %v kicks\n", conn)

		self.takeToken()
		self.pending <- conn
	}
}

func (self *Server) Stop() {
	self.listener.Close()
}


func (self *Server) listen() {

	go func() {
		for {
			select {
			case conn := <-self.pending:
				self.join(conn)
			case conn := <-self.quiting:
				self.leave(conn)
			}
		}
	}()
}

func (self *Server) join(conn net.Conn) {

	session := CreateSession(conn)
	id := TypeSessionID(uniq.GetUniq())
	session.SetID(id)
	self.sessions[id] = session

	go func() {
		<-session.quiting
		delete(self.sessions, session.GetID())
		self.quiting <- session.GetConn()
	}()

}

func (self *Server) leave(conn net.Conn) {
	if conn != nil {
		conn.Close()
	}
	self.generateToken()
}
