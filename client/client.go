// Package client 는 서버에서 클라이언트를 나타내는 패키지 이다
package client

import "net"

// Client 는 client 패키지를 나타내는 구조체
type Client struct {
	Conn net.Conn
}
