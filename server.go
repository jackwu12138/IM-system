package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
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

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s %s", user.Addr, user.Name, msg)

	s.Message <- sendMsg
}

func (s *Server) handler(conn net.Conn) {
	fmt.Println("Connection established successfully")

	user := NewUser(conn)

	// The user is online, add the user to the OnlineMap
	s.MapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.MapLock.Unlock()

	// Broadcast the message that the user is online
	s.BroadCast(user, "is online")
	// Receive messages from users
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, " is leaved")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err: ", err)
				return
			}
			// Get messages sent by users
			msg := string(buf[:n-1])
			// broadcast the message
			s.BroadCast(user, msg)
		}
	}()

	select {}
}

func (s *Server) start() {
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
