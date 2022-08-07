package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	IP   string
	Post int
}

func NewServer(IP string, post int) *Server {
	return &Server{IP: IP, Post: post}
}

func (s *Server) handler(conn net.Conn) {
	fmt.Println(conn)
	fmt.Println("建立连接成功")
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
