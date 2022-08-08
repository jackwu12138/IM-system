package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Post int

	// online user map
	OnlineMap map[string]*User
	MapLock   sync.Mutex

	// message broadcast channel
	Message chan string
}

func NewServer(IP string, post int) *Server {
	return &Server{IP: IP, Post: post, OnlineMap: make(map[string]*User), Message: make(chan string)}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		// Send message to all online users
		s.MapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.MapLock.Unlock()
	}
}

func (s *Server) handler(conn net.Conn) {
	fmt.Println("Connection established successfully")

	user := NewUser(conn, s)
	user.online()
	isLive := make(chan bool)

	// Receive messages from users
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err: ", err)
				return
			}
			// Get messages sent by users
			msg := string(buf[:n-1])
			// handle message
			user.handleMessage(msg)
			isLive <- true
		}
	}()

	for {
		select {

		case <-isLive:
			// The current user is active and the timer should be reset
		case <-time.After(time.Minute * 10):
			// user is time out, kick out the user
			user.sendMsg("you got kicked out")

			close(user.C)
			err := conn.Close()
			if err != nil {
				fmt.Println("close connection err: ", err)
				return
			}
			return
		}
	}
}

func (s *Server) Start() {
	fmt.Println("server start")
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Post))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("close net.Listen err: ", err)
			os.Exit(-1)
		}
	}(listener)

	// Start a goroutine that listens for Messages
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}
		// do handler
		go s.handler(conn)
	}
}
