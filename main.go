package main

import (
	"fmt"
	_ "go_chat_server/db"
	_ "go_chat_server/global"
	"go_chat_server/tcp"
	"net"
)

//1.创建聊天室
//2.加入聊天室
//3.发送、广播消息
//4.获取历史聊天记录

func main() {
	// 开一个tcp监听
	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("lister error = ", err)
		return
	}
	fmt.Println("chat server start ...")
	// 在函数结束时，关闭监听
	defer listener.Close()
	go tcp.WaitConnect(listener)

}
