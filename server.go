package main

import "net"

// Server 는 서버를 나타내는 구조체
type Server struct {
	port         string
	clients      []*Client
	addClient    chan *Client
	removeClient chan *Client
	broadcast    chan string
}

// NewServer 는 새로운 서버객체를 만들어 반환한다
func NewServer(port string) *Server {
	return &Server{
		port:         port,
		clients:      make([]*Client, 0),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		broadcast:    make(chan string),
	}
}

// ListenAndServe 는 서버의 리슨과 서비스를 시작한다
func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return l.Close()
	}
	defer l.Close()

	go s.serve()

	for {
		conn, err := l.Accept()
		if err != nil {
			Log.ErrorF("Error: %s\n", err)
		}
		c := NewClient(conn, s)

		s.addClient <- c

		go c.Read()
		go c.Write()
	}
}

// serve 는 클라이언트의 추가 or 삭제, 브로드 케스팅을 제공
func (s *Server) serve() {
	for {
		select {
		case c := <-s.addClient:
			Log.Info("Join: " + c.conn.RemoteAddr().String())
			s.clients = append(s.clients, c)
		case c := <-s.removeClient:
			for i := range s.clients {
				if s.clients[i] == c {
					s.clients = append(s.clients[:i], s.clients[i+1:]...)
					Log.Info("Leave: " + c.conn.RemoteAddr().String())
					break
				}
			}
		case m := <-s.broadcast:
			for _, c := range s.clients {
				Log.InfoF("Broadcast (%s): \"s\"\n", c.conn.RemoteAddr(), m)
				c.channel <- m
			}
		}
	}
}
