package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

// 创建一个Server类
type Server struct {
	Ip   string
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
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就会发送给所有的在线user
func (server *Server) ListenMessager() {
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
	// 用户上线 将用户加入到onlineMap中
	user := NewUser(conn, server)

	user.Online()

	isLive := make(chan bool)

	// 接收用户输入信息，将消息进行广播
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// 关闭客户端请求时触发
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err:", err)
				return
			}

			// 提取用户消息 去除\n
			msg := string(buf[:n-1])

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息 代表当前用户是活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的 应该重置定时器
			// 不做任何事情 为了激活select 更新下面的计时器
			// 定时器被垃圾回收
		case <-time.After(time.Second * 10):
			// time.After()每次执行select时都会创建一个新的定时器
			// 超时强踢
			user.SendMsg("You have been kicked due to prolonged inactivity!\n")
			// 销毁用的资源
			close(user.C)
			// 关闭客户端连接
			conn.Close()

			// 关闭协程
			runtime.Goexit()
			// return 也可以
		}
	}
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
