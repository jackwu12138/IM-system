package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{Name: userAddr, Addr: userAddr, C: make(chan string), conn: conn}

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
