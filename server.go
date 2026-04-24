package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]:%s:%s", user.Addr, user.Name, msg)
	this.Message <- sendMsg

}
func (this *Server) Handler(con net.Conn) {
	// fmt.Println("链接建立成功")

	user := NewUser(con, this)
	user.Online()
	isAlive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := con.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err.Error() != "EOF" {
				fmt.Println("con.Read error:", err)
				return
			}

			msg := string(buf[:n])
			user.DoMessage(msg)
			isAlive <- true
		}
	}()

	for {
		select {
		case <-isAlive:
		case <-time.After(time.Second * 600):
			user.SendMsg("you are kicked and offline")
			close(user.C)
			user.con.Close()
			return
		}
	}
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer listener.Close()
	go this.ListenMessager()
	fmt.Println("server start at:", this.Ip, this.Port)
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept error:", err)
			continue
		}

		go this.Handler(con)
	}
}
