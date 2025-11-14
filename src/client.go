package main

import (
	"fmt"
	"net"
)

// 定义Client类
type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
}

// 创建Client对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client {
		ServerIp: serverIp,
		ServerPort: serverPort,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net Dial err: ", err)
		return nil
	}

	client.conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}

	fmt.Println("连接服务器成功")
}