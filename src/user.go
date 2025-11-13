package main

import "net"

// 创建User类
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 新建User方法
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息就直接发送给对端客户端
func (user *User) ListenMessage() {
	for {
		msg := <- user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
