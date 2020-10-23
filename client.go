package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client 는 클라이언트를 나타내는 구조체
type Client struct {
	conn    net.Conn
	server  *Server
	channel chan string
	done    chan struct{}
}

// NewClient 는 새로운 쿨라이언트 객체를 만들어 반환
func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		conn:    conn,
		server:  server,
		channel: make(chan string),
		done:    make(chan struct{}),
	}
}

// Read 는 클라이언트 로부터 읽은 데이터를 처리한다
func (c *Client) Read() {
	scanner := bufio.NewScanner(c.conn)
LOOP:
	for scanner.Scan() {
		select {
		case <-c.done:
			break LOOP
		default:
			m := scanner.Text()
			Log.InfoF("Receive (%s): \"%s\"\n", c.conn.RemoteAddr(), m)
			c.server.broadcast <- m
		}
	}

	if err := scanner.Err(); err != nil {
		Log.ErrorF("Error (read): %s\n", err)
	}
	c.server.removeClient <- c
	c.done <- struct{}{}
}

// Write 는 클라이언트에 쓰기 작업을 처리한다
func (c *Client) Write() {
LOOP:
	for {
		select {
		case <-c.done:
			break LOOP
		case m := <-c.channel:
			Log.InfoF("Send (%s): \"%s\"\n", c.conn.RemoteAddr(), m)
			_, err := fmt.Fprintln(c.conn, m)
			if err != nil {
				Log.ErrorF("Error (write): %s\n", err)
				break LOOP
			}
		}
	}

	c.server.removeClient <- c
	c.done <- struct{}{}
}
