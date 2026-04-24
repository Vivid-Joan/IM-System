package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func (client *Client) menu() bool {
	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. Update username")
	fmt.Println("0. Exit")

	var flag int
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	}

	fmt.Println(">>>Please input valid number")
	return false
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>Please input username:")

	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("UpdateName conn.Write error:", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	fmt.Println(">>>Please input chat content, exit by inputting 'exit'")
	var chatMsg string
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("PublicChat conn.Write error:", err)
				return
			}
		}
		chatMsg = ""
		fmt.Println(">>>Please input chat content, exit by inputting 'exit'")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("SelectUsers conn.Write error:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	client.SelectUsers()
	fmt.Println(">>>Please input chat object, exit by inputting 'exit'")

	var remoteName string
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>Please input chat content, exit by inputting 'exit'")
		var chatMsg string
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("PrivateChat conn.Write error:", err)
					return
				}
			}
			chatMsg = ""
			fmt.Println(">>>Please input chat content, exit by inputting 'exit'")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		remoteName = ""
		fmt.Println(">>>Please input chat object, exit by inputting 'exit'")
		fmt.Scanln(&remoteName)
	}
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			// fmt.Println(">>>Public chat")
			client.PublicChat()
		case 2:
			// fmt.Println(">>>Private chat")
			client.PrivateChat()
		case 3:
			// fmt.Println(">>>Update username")
			client.UpdateName()
		}
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP address, default is 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8080, "Server port, default is 8080")
}
func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>connect to server error")
		return
	}

	go client.DealResponse()
	fmt.Println(">>>connect to server success")

	client.Run()
}
