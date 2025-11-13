package main

import (
	"net"
	"strings"
)

// 创建User类
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 新建User方法
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息就直接发送给对端客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}

// 用户的上线功能
func (user *User) Online() {
	// 用户上线 将用户加入到onlineMap中
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "is online.\n")
}

// 用户的下线功能
func (user *User) Offline() {
	// 用户下线 将用户从onlineMap中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "is offline.\n")
}

// 给当前user对应的客户端发消息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息
func (user *User) DoMessage(msg string) {
	// 根据用户输入信息，执行不同的操作
	if strings.TrimSpace(msg) == "who" {
		// 加锁
		user.server.mapLock.Lock()
		// 遍历在线用户map
		for _, u := range user.server.OnlineMap {
			msg := "[" + u.Addr + "]" + u.Name + ": is online.\n"
			user.SendMsg(msg)
		}
		user.server.mapLock.Unlock()
	} else if len(strings.TrimSpace(msg)) > 7 && strings.TrimSpace(msg)[:7] == "rename|" {
		// 首先获取更新后的名称
		newName := strings.Split(strings.TrimSpace(msg), "|")[1]
		// 查询更改后的名称是否已被他人使用
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg("Current name is already in use and cannot be changed!")
		} else {
			// 加锁
			user.server.mapLock.Lock()
			// 删除OnlineMap中的当前用户
			delete(user.server.OnlineMap, user.Name)
			// 更新当前用户的名称
			user.Name = newName
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()
			user.SendMsg("Your name has been changed to " + newName + "!\n")
		}
	} else {
		user.server.BroadCast(user, msg)
	}
}
