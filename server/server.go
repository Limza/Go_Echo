// Package server 는 서버 작업에 관련된 작업들을 처리하기 위해 만들었다
package server

import (
	"echoserver/client"
	"log"
)

type message string

// Server 는 server 패키지를 나타내는 구조체
type Server struct {
	port         string
	clients      []*client.Client
	addClient    chan *client.Client
	removeClient chan *client.Client
	broadcast    chan message
}

// NewServer 는 새로운 서버객체를 만들어 반환한다
func NewServer(port string) *Server {
	return &Server{
		port:         port,
		clients:      make([]*client.Client, 0),
		addClient:    make(chan *client.Client),
		removeClient: make(chan *client.Client),
		broadcast:    make(chan message),
	}
}

// serve 는 클라이언트의 추가 or 삭제, 브로드 케스팅을 제공
func (s *Server) serve() {
	for {
		select {
		case c := <-s.addClient:
			log.Println("Join: " + c.Conn.RemoteAddr().String())
			s.clients = append(s.clients, c)
		}
	}
}
