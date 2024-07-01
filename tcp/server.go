package tcpServer

import (
	"fmt"
	"net"
)

/*Server type */

type Message struct {
	From    string
	Payload []byte
}

type Server struct {
	ListenAddr string
	Ln         net.Listener
	Quitch     chan struct{}
	Msgs       chan Message
}

/*constructor inintialization */

func NewServer(lstnAddr string) *Server {
	return &Server{
		ListenAddr: lstnAddr,
		Quitch:     make(chan struct{}),
		Msgs:       make(chan Message, 10),
	}
}

/* using start method to boot up the server for tcp communication over the network */
func (s *Server) Start() error {

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("error occured while starting up the server : %v", err)
	}

	defer ln.Close()
	s.Ln = ln

	go s.connectionLoop()

	<-s.Quitch

	close(s.Quitch)

	return nil
}

/* as the name suggests the upcoming or downcoming function contineously listens for  the incoming connections requests */

func (s *Server) connectionLoop() {
	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			fmt.Printf("error while inintialization of connection listener : %v", err)
			continue
		}

		fmt.Println("new connection to the server :", conn.RemoteAddr())

		go s.readLoop(conn)
	}

}

/* readloop inintialization for reading messages from the connection */

func (s *Server) readLoop(conn net.Conn) error {
	defer conn.Close()
	buf := make([]byte, 2029)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("error reading from connection : %v", err)
		}
		msg := buf[:n]

		s.Msgs <- Message{
			From:    conn.RemoteAddr().String(),
			Payload: []byte(msg),
		}

		conn.Write([]byte("message sent!\n"))

	}

}
