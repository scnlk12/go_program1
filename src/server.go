package main

import (
	"fmt"
	"net"
	"sync"
)

// 创建一个Server类
type Server struct {
	Ip string
	Port int

	// 增加在线用户map
	OnlineMap map[string]*User
	// 加锁
	mapLock sync.RWMutex

	// 消息广播的channel
	Message chan string
}


// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就会发送给所有的在线user
func (server *Server)ListenMessager()  {
	for {
		// 监听channel
		msg := <-server.Message
		// 分发给所有的在线用户
		// 加锁
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// 广播消息的方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

// handler
func (server *Server) Handler(conn net.Conn) {
	// ... 当前链接的业务
	// fmt.Println("链接建立成功")

	// 用户上线 将用户加入到onlineMap中
	user := NewUser(conn)

	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线消息
	server.BroadCast(user, "is online.\n")

	// 当前handler阻塞

}

// 启动服务器的接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动监听message的goroutine
	go server.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// do handler
		go server.Handler(conn)
	}

}