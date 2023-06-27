package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	header := http.Header{}
	header.Add("Key", "123456789")
	// 建立WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial("ws://13.250.39.11:8080/ws", header)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 循环接收服务器发送的消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", message)
	}
}
