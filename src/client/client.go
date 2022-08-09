package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       -1,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIP, client.ServerPort))
	if err != nil {
		fmt.Println("connection failure ...")
		return nil
	}
	client.conn = conn

	return client
}

func (c *Client) Menu() bool {
	var n int
	fmt.Println("1. chatroom")
	fmt.Println("2. private chat")
	fmt.Println("3. modify username")
	fmt.Println("0. exit")

	_, err := fmt.Scanf("%d", &n)
	if err != nil {
		fmt.Println("error input")
		return false
	}
	if n > 3 || n < 0 {
		fmt.Println("error input")
		return false
	}
	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for !c.Menu() {
		}
		switch c.flag {
		case 1:
			break
		case 2:
			break
		case 3:
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set the ip address of the server")
	flag.IntVar(&serverPort, "port", 8888, "Set the server port number")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Failed to start client [" + fmt.Sprintf("%s:%d", serverIp, serverPort) + "]")
		return
	}
	fmt.Println("Start the client successfully")

	select {}
}
