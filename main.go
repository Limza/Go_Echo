// https://jacking75.github.io/go_network_tcp_echo_server/
package main

import (
	"echoserver/logger"
	"log"
)

// PORT 는 접속하는 서버의 포트 번호
const PORT = "7777"

// Log 는 main에서 log를 나타낼때 사용할 변수
var Log *logger.Logger

func init() {
	var err error
	Log, err = logger.NewLogger(logger.Day, "./logs", "server")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server := NewServer(PORT)
	if err := server.ListenAndServe(); err != nil {
		Log.Error(err)
	}
}
