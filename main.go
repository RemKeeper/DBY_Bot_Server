package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Response struct {
	Message string `json:"message"`
	RetCode int    `json:"retcode"`
}

var MessageChan = make(chan string, 50)
var ConnectKey string

func main() {
	file, err := os.ReadFile("./ConnectKey")
	if err != nil {
		fmt.Println("未找到ConnectKey文件,请在ConnectKey中配置一个连接密钥")
		_, _ = os.Create("ConnectKey")
		return
	}
	if len(file) > 8 {
		ConnectKey = string(file)
	} else {
		fmt.Println("密钥过短，请设置8位以上密钥")
		return
	}
	go http.ListenAndServe(":80", CallBack())
	http.ListenAndServe(":8080", WebSocket())
}

func CallBack() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		all, err := io.ReadAll(request.Body)
		if err != nil {
			return
		}
		fmt.Println(string(all))

		response := Response{
			Message: "",
			RetCode: 0,
		}
		marshal, err := json.Marshal(response)
		if err != nil {
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(marshal)
		MessageChan <- string(all)
	})
	return mux
}

func WebSocket() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Key") != ConnectKey {
			log.Println("认证失败")
			conn, err := upgrader.Upgrade(writer, request, nil)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer conn.Close()
			// 发送错误消息并立即关闭连接
			conn.WriteMessage(websocket.TextMessage, []byte("认证失败"))
			conn.Close()
			return
		}
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer conn.Close()
		for {
			// 每隔1秒发送当前时间的字符串给客户端
			err := conn.WriteMessage(websocket.TextMessage, []byte(<-MessageChan))
			if err != nil {
				log.Println(err)
				return
			}
		}
	})
	return mux
}
