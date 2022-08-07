package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{Name: userAddr, Addr: userAddr, C: make(chan string), conn: conn, server: server}

	go user.ListenMessage()
	return user
}

func (user *User) ListenMessage() {
	for {
		message := <-user.C
		_, err := user.conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("listen message err: ", err)
			return
		}
	}
}

// user online
func (user *User) online() {
	// The user is online, add the user to the OnlineMap
	user.server.MapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.MapLock.Unlock()

	// Broadcast the message that the user is online
	user.handleMessage("is online")
}

// user offline
func (user *User) offline() {
	// The user is offline, delete the user from the OnlineMap
	user.server.MapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.MapLock.Unlock()

	// Broadcast the message that the user is online
	user.handleMessage("is offline")
}

func (user *User) sendMsg(msg string) {
	_, err := user.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("send message err: ", err)
		return
	}
}

// user handles message
func (user *User) handleMessage(msg string) {
	// query all online users
	if msg == "who" {
		user.server.MapLock.Lock()
		for _, u := range user.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s is online\n", u.Addr, u.Name)
			user.sendMsg(onlineMsg)
		}
		user.server.MapLock.Unlock()
		return
	}
	sendMsg := fmt.Sprintf("[%s]%s %s", user.Addr, user.Name, msg)

	user.server.Message <- sendMsg
}
