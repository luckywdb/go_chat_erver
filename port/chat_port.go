package port

import (
	"chat/db"
	"chat/lib"
	"chat/role"
	"chat/room"
	"fmt"
	"net"
	"time"
)

// 用与和客户端通信

// 登陆
// 参数：name （用户名）
// 返回值： ok|error
func Login(conn net.Conn, args map[string]string) string {
	// TODO 1. 检查用户是否创建， 没有用户则创建
	name := args["name"]
	r := role.CreateRole(name)
	// 2. 用户登陆
	r.SetLoginTime(time.Now().Unix())
	db.UpdateRole(r) // 更新role对象
	fmt.Println(conn.RemoteAddr(), args["name"], "login ok ...")
	return "ok"
}

// 创建聊天室
// 参数: name (创建者)
// 返回值：聊天室id
func Create(conn net.Conn, args map[string]string) string {
	id := lib.GetId()
	title := args["title"]
	name := args["name"]
	chatRoom := room.CreatorRoom(id, title, name) // 创建新到聊天室
	chatRoom.AddMember(name)                      // 把创建者添加到聊天室成员列表
	db.UpdateChatRoom(chatRoom)

	r := db.GetRole(name)
	r.AddChatRoomId(id)
	db.UpdateRole(r)

	fmt.Println(conn.RemoteAddr(), args["name"], "creat a chat room ...")
	return id
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

	return "ok"
}
