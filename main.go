package main

import (
	"chat/port"
	"fmt"
	"net"
	"strings"
)

//1.创建聊天室
//2.加入聊天室
//3.发送、广播消息
//4.获取历史聊天记录

// 存储 用户名和对应的socket连接
var conns = make(chan map[string]net.Conn)

func main() {
	// 开一个tcp监听
	listener, err := net.Listen("tcp", "127.0.0.1：8088")
	if err != nil {
		fmt.Println("lister error = ", err)
		return
	}
	// 在函数结束时，关闭监听
	defer listener.Close()

	// 等待客户端连接，连接成功后，等待下一个连接
	loop(listener)

	fmt.Println("chat server start ...")
}

func loop(listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("accept error = ", err)
		return
	}

	fmt.Println(conn.RemoteAddr(), "connect success ...")

	go handleConn(conn)

	// 等待下一个客户端连接
	loop(listener)
}

func handleConn(conn net.Conn) {
	for {
		data := make([]byte, 512)

		msgLen, err := conn.Read(data)
		if msgLen == 0 || err != nil {
			continue
		}

		msgs := strings.Split(string(data[0:msgLen]), "?")
		// 解析参数
		args := encode(msgs[1])
		switch msgs[0] {
		// 登陆 : login?name(用户名)="aaa"
		case "login":
			if checkArgs(args, []string{"name"}) {
				conns1 := <-conns
				// 处理登陆
				result := port.Login(conn, args)
				conns1[args["name"]] = conn
				conns <- conns1
				conn.Write([]byte("result:" + result))
			} else {
				conn.Write([]byte("result:args error"))
			}
		// 创建聊天室 : create?name(创建者)="aaa"&title(聊天室名字)="；bbb"
		case "create":
			if checkArgs(args, []string{"name", "title"}) {
				conns1 := <-conns
				conns <- conns1
				// 检查是否已经登陆
				if _, ok := conns1["name"]; ok { // 已经登陆
					// 处理创建聊天室
					result := port.Create(conn, args)
					// 返回聊天室 id
					conn.Write([]byte("result:" + result))
				} else {
					conn.Write([]byte("result:not login"))
				}
			} else {
				conn.Write([]byte("result:args error"))
			}

		// 加入聊天室 : join?name(加入者)="aaa"&id(聊天室id)=123
		case "join":
			// 处理加入聊天室
			result := port.Join(conn, args)

		// 发送消息 : say?name(发送者)="aaa"&id(聊天室id)=123&msg(消息)="abc 132,sdf"
		case "say":

			// 处理发送消息
			result := port.Say(conn, args)
		// 退出 : quit?name(退出者)="aaa"&id(聊天室id)=123
		case "quit":
			result := port.Quit(conn, args)
		}

	}
}

func encode(msg string) map[string]string {
	result := make(map[string]string)
	values := strings.Split(msg, "&")
	for _, v := range values {
		keyValues := strings.Split(v, "=")
		result[keyValues[0]] = keyValues[1]
	}

	return result
}

// 参数检查
// 通过返回 true， 不通过返回false
func checkArgs(args map[string]string, argsName []string) bool {
	for _, name := range argsName {
		if _, ok := args[name]; ok {
			continue
		} else {
			return false
		}
	}
	return true
}
