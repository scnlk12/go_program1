package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// 定义Client类
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	input      int
}

// 创建Client对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		input:      999,
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

// 显示菜单
func (client *Client) menu() bool {
	var input int

	fmt.Println("1. BroadCast...")
	fmt.Println("2. Private Chat...")
	fmt.Println("3. Change Name...")
	fmt.Println("0. Exit")

	// 接收用户输入
	fmt.Scanln(&input)

	if input >= 0 && input <= 3 {
		client.input = input
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) BroadCast() {
	// 提示用户输入内容
	fmt.Println("请输入消息内容, exit表示退出")
	// 获取用户输入
	var msg string
	fmt.Scanln(&msg)
	// 循环判断 用户输入exit表示退出
	for msg != "exit" {
		// 判断内容是否为空
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err: ", err)
				break
			}
		}
		msg = ""
		fmt.Println("请输入消息内容, exit表示退出")
		fmt.Scanln(&msg)
	}
}

// 处理Server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

// Run
func (client *Client) Run() {
	for client.input != 0 {
		for client.menu() != true {
		}

		// 根据不同模式处理不同操作
		switch client.input {
		case 1:
			client.BroadCast()
			break
		case 2:
			fmt.Println("私聊...")
			break
		case 3:
			client.UpdateName()
		}
	}
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

	// 单独开启一个goroutine处理server的回执消息
	go client.DealResponse()

	fmt.Println("连接服务器成功")
	client.Run()
}
