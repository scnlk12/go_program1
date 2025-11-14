package main

import (
	"flag"
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

// 执行顺序 pkg -> const -> var -> init -> main

// 定义全局变量
var serverIp string
var serverPort int

// 优先于main执行
func init() {
	// 定义和注册命令行参数
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器的ip地址")
	flag.IntVar(&serverPort, "port", 8888, "服务器的端口")
}

func main() {

	// 命令行解析 解析命令行参数
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}

	fmt.Println("连接服务器成功")
}