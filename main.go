package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	out, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}

// 递归遍历目录并收集文件
func collectFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// 递归调用 collectFiles 函数
			subFiles, err := collectFiles(dir + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, dir+"/"+entry.Name())
		}
	}
	return files, nil
}

func listVideosHandler(w http.ResponseWriter, r *http.Request) {
	// 调用辅助函数递归收集文件
	videoFiles, err := collectFiles("videos")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error reading videos directory",
			"list":    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"list":    videoFiles,
	})
}

func main() {
	// 默认端口号
	port := 443

	// 解析命令行参数
	flag.IntVar(&port, "port", port, "Port number to listen on")
	flag.Parse()

	flag.Usage()

	// 加载SSL证书和私钥
	certFile := "cert.pem"
	keyFile := "key.pem"

	// 添加静态文件服务
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/", http.StripPrefix("/", fs))

	// 视频文件服务
	video_fs := http.FileServer(http.Dir("videos"))
	http.Handle("/videos/", http.StripPrefix("/videos/", video_fs))

	// websocket服务
	http.Handle("/ws", websocket.Handler(wsHandler))

	// API服务
	http.HandleFunc("/api", handler)

	// 文件上传服务
	http.HandleFunc("/api/upload", uploadHandler)

	// 列出视频文件服务
	http.HandleFunc("/api/videos", listVideosHandler)

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
					fmt.Printf("Access URL: https://%s:%d", ipnet.IP.String(), port)
				}
			}
		}
	}

	// 创建HTTPS服务器
	// 创建HTTPS服务器并监听指定端口
	addr := fmt.Sprintf(":%d", port)
	err = http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	if err != nil {
		fmt.Println("Error starting HTTPS server:", err)
	}
}
