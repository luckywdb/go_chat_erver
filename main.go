package main

import (
	_ "chat/db"
	"chat/port"
	"fmt"
	"net"
	"strings"
	"sync"
)

//1.创建聊天室
//2.加入聊天室
//3.发送、广播消息
//4.获取历史聊天记录

// 数据库标记 用来控制操作数据库的 goroutine 同一时间只有一个
var dbState chan bool

var wg sync.WaitGroup

func main() {
	dbState := make(chan bool)
	// 开一个tcp监听
	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("lister error = ", err)
		return
	}
	fmt.Println("chat server start ...")
	// 在函数结束时，关闭监听
	defer listener.Close()

	wg.Add(1)
	// 等待客户端连接，连接成功后，等待下一个连接
	go loop(listener, dbState)
	dbState <- true

	wg.Wait()
}

func loop(listener net.Listener, dbState chan bool) {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("accept error = ", err)
		return
	}

	fmt.Println(conn.RemoteAddr(), "connect success ...")

	go handleConn(conn, dbState)
	// 等待下一个客户端连接
	loop(listener, dbState)
}

func handleConn(conn net.Conn, dbState chan bool) {
	defer conn.Close()
	for {
		data := make([]byte, 512)

		msgLen, err := conn.Read(data)
		if msgLen == 0 || err != nil {
			continue
		}
		msgs := strings.Split(string(data[0:msgLen]), "?")
		fmt.Println("msgs: ", msgs)
		// 解析参数
		if len(msgs) > 1 {
			args := encode(msgs[1])
			switch msgs[0] {
			// 登陆 : login?name(用户名)="aaa"
			case "login":
				if checkArgs(args, []string{"name"}) {
					port.Login(dbState, conn, args)
				} else {
					conn.Write([]byte("result:args error"))
				}
				// 登出 : login?name(用户名)="aaa"
			case "logout":
				if checkArgs(args, []string{"name"}) {
					port.Logout(dbState, conn, args)
					break
				} else {
					conn.Write([]byte("result:args error"))
				}
			// 创建聊天室 : create?name(创建者)="aaa"&title(聊天室名字)="；bbb"
			case "create":
				if checkArgs(args, []string{"name", "title"}) {
					// 处理创建聊天室
					port.Create(dbState, conn, args)
				} else {
					conn.Write([]byte("result:args error"))
				}

			// 加入聊天室 : join?name(加入者)="aaa"&id(聊天室id)=123
			case "join":
				if checkArgs(args, []string{"name", "id"}) {
					port.Join(dbState, conn, args)
				} else {
					conn.Write([]byte("result:args error"))
				}

			// 发送消息 : say?name(发送者)="aaa"&id(聊天室id)=123&msg(消息)="abc 132,sdf"
			case "say":
				if checkArgs(args, []string{"name", "id", "msg"}) {
					port.Say(dbState, conn, args)
				} else {
					conn.Write([]byte("result:args error"))
				}

			// 退出 : quit?name(退出者)="aaa"&id(聊天室id)=123
			case "quit":
				if checkArgs(args, []string{"name", "id"}) {
					port.Quit(dbState, conn, args)
				} else {
					conn.Write([]byte("result:args error"))
				}
			}
		}
	}
}

func encode(msg string) map[string]string {
	result := make(map[string]string)
	values := strings.Split(msg, "&")
	for _, v := range values {
		keyValues := strings.Split(v, "=")
		if len(keyValues) == 2 {
			result[keyValues[0]] = keyValues[1]
		}
	}
	fmt.Println("encode:", result)
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
