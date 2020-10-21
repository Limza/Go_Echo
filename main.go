package main

import (
	"echoserver/logger"
	"log"
)

// PORT 는 접속하는 서버의 포트 번호
const PORT = 7777

func init() {
	if _, err := logger.NewLogger(logger.Day, "./logs", "server"); err != nil {
		log.Fatal(err)
	}
}

func main() {

}
