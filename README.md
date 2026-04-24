## 描述
基于Golang实现的简单TCP即时通信服务端和客户端，提供公聊，私聊，查询在线客户端，修改客户端用户名功能

```sh
运行服务端
go build -o server main.go server.go user.go
./server

运行客户端
go build -o client client.go
./client

客户端输入
1：公聊
2：私聊，需要先输入私聊用户名，再输入聊天内容
3：修改当前用户名
0：退出程序
```
