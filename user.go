package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	con  net.Conn

	server *Server
}

func NewUser(con net.Conn, server *Server) *User {
	userAddr := con.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		con:    con,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.con.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "on line")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "off line")
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]:" + user.Name + ":is on line\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("this username already used")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[this.Name] = this
			this.server.mapLock.Unlock()
			this.SendMsg("rename success:" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("message format error, please use \"to|username|message\"")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("this user do not exit")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("there is no message to send, please use \"to|username|message\"")
			return
		}
		remoteUser.SendMsg(this.Name + " send you a message:" + content)
	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) SendMsg(msg string) {
	this.con.Write([]byte(msg))
}
