package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Message struct {
	From    string
	Payload []byte
}
type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message),
	}
}

func (s *Server) Start() error {
	ln, error := net.Listen("tcp", s.listenAddr)
	if error != nil {
		return error
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	<-s.quitch
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error")
			continue
		}
		fmt.Println("new connection to the server", conn.RemoteAddr())
		go s.readLoop(conn)
	}

}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error")
			break
		}
		msg := buf[:n]
		s.msgch <- Message{
			From:    conn.RemoteAddr().String(),
			Payload: msg,
		}
		conn.Write([]byte("thank u"))
	}

}

func (s *Server) quitServer() {
	time.Sleep(10 * time.Second)
	close(s.quitch)
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("received a message from %s:%s\n", msg.From, msg.Payload)
		}
	}()

	log.Fatal(server.Start())

}
