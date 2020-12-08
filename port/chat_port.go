package port

import (
	"fmt"
	"net"
)

// 用与和客户端通信

// 登陆
// 参数：name （用户名）
// 返回值： ok
func Login(conn net.Conn, args map[string]string) string {
	fmt.Println(conn.RemoteAddr(), args["name"], "login ...")
	return "ok"
}

// 创建聊天室
// 参数: name (创建者)
// 返回值：聊天室id
func Create(conn net.Conn, args map[string]string) string {
	fmt.Println(conn.RemoteAddr(), args["name"], "creat a chat room ...")
	return "123"
}

// 加入聊天室
// 参数：name （加入者）& id （聊天室id）
// 返回值：聊天室中所有人的name
func Join(conn net.Conn, args map[string]string) string {
	fmt.Println(conn.RemoteAddr(), args["name"], "join chat room ", args["id"], " ...")
	return "张三李四王二麻子"
}

// 发消息
// 参数：name （消息发送者） & id（聊天室id）& msg（消息内容）
// 返回值： ok
func Say(conn net.Conn, args map[string]string) string {
	fmt.Println(conn.RemoteAddr(), args["name"], "in chat room ", args["id"], " say ", args["msg"], " ...")
	return "ok"
}

//  退出
// 参数： name （退出者）& id (聊天室id)
// 返回值： ok
func Quit(conn net.Conn, args map[string]string) string {
	fmt.Println(conn.RemoteAddr(), args["name"], "quit chat room", args["id"], " ...")
}
