package main

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

// 结构体保存client和msg
type Msg struct {
	ws  *websocket.Conn
	msg string
}

var clients = make(map[*websocket.Conn]bool) // 连接映射
var broadcast = make(chan Msg)               // 广播通道

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello API, HTTPS!")
}

func wsHandler(ws *websocket.Conn) {
	fmt.Println("WebSocket connection established")
	// 将新连接添加到连接列表
	clients[ws] = true

	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			// 移除连接
			delete(clients, ws)
			break
		}
		// 将接收到的消息发送到广播通道
		broadcast <- Msg{ws, msg}
	}
}

func handleMessages() {
	for {
		// 从广播通道接收消息
		msg := <-broadcast
		// 将消息发送给所有连接的客户端
		for client := range clients {
			if client == msg.ws {
				continue
			}
			err := websocket.Message.Send(client, msg.msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// 加载SSL证书和私钥
	certFile := "cert.pem"
	keyFile := "key.pem"

	// 添加静态文件服务
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/", http.StripPrefix("/", fs))

	// websocket服务
	http.Handle("/ws", websocket.Handler(wsHandler))

	// API服务
	http.HandleFunc("/api", handler)

	// 启动处理广播消息的协程
	go handleMessages()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting local IP addresses:", err)
	} else {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Println("Local IP address:", ipnet.IP.String())
					// 访问的URL
					fmt.Println("Access URL:", "https://"+ipnet.IP.String()+":443")
				}
			}
		}
	}

	// 创建HTTPS服务器
	err = http.ListenAndServeTLS(":443", certFile, keyFile, nil)
	if err != nil {
		fmt.Println("Error starting HTTPS server:", err)
	}
}
