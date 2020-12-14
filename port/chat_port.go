package port

import (
	"chat/db"
	"chat/lib"
	"chat/role"
	"chat/room"
	"fmt"
	"net"
	"reflect"
	"time"
)

// 用与和客户端通信

// 登陆
// 参数：name （用户名）
// 返回值： ok|error
func Login(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		<-dbState // 获取到db操作权限
		result := func() string {
			name := args["name"]
			r, err := db.GetRole(name)
			if err != nil {
				fmt.Println(err)
				return "get role error"
			}
			if reflect.DeepEqual(r, role.Role{}) { // 没有用户
				r = role.CreateRole(name) // 创建用户
			}
			// 2. 用户登陆
			r.LoginTime = time.Now().Unix()
			err = db.UpdateRole(r) // 更新role对象
			if err != nil {
				fmt.Println(err)
				return "update role error"
			}
			fmt.Println(conn.RemoteAddr(), args["name"], "login ok ...")
			db.UpdateConnects(name, conn, db.AddConnects)
			return "ok"
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true
	}()
}

// 登出
// 参数：name （用户名）
// 返回值： ok|error
func Logout(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		<-dbState // 获取到db操作权限
		result := func() string {
			name := args["name"]
			if db.GetConnects(name) != nil { // 在线
				r, err := db.GetRole(name)
				if err != nil {
					fmt.Println(err)
					return "get role error"
				}
				if reflect.DeepEqual(r, role.Role{}) {
					fmt.Println(err)
					return "no role"
				}
				//  向拥有的所有的聊天室通知，当前用户退出
				go func() {
					for _, chatRoomId := range r.ChatRooms { // 遍历用户所有的聊天室
						chatRoom, err := db.GetChatRoom(chatRoomId)
						if err != nil || reflect.DeepEqual(chatRoom, room.ChatRoom{}) {
							fmt.Println(err)
							continue
						}
						chatRoom.Broadcast(name, name+" logout ", db.GetAllConnects()) // 向聊天室中的成员广播
						db.UpdateConnects(name, conn, db.DeleteConnects)
					}
				}()
				fmt.Println(conn.RemoteAddr(), args["name"], "logout ok ...")
				return "logout"
			} else {
				return "not login"
			}
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true // 释放db操作权限
	}()
}

// 创建聊天室
// 参数: name (创建者)
// 返回值：聊天室id
func Create(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		<-dbState // 获取到db操作权限
		result := func() string {
			name := args["name"]
			if db.GetConnects(name) != nil {
				id := lib.GetId()
				title := args["title"]
				name := args["name"]
				r, err := db.GetRole(name)
				if err != nil {
					fmt.Println(err)
					return "get role error"
				}
				if reflect.DeepEqual(r, role.Role{}) {
					fmt.Println(err)
					return "no role"
				}
				chatRoom := room.CreatorRoom(id, title, name) // 创建新到聊天室
				chatRoom.AddMember(name)                      // 把创建者添加到聊天室成员列表

				if r.AddChatRoomId(id) {
					err := db.UpdateChatRoomRole(chatRoom, r)
					if err != nil {
						fmt.Println(err)
						return "create chat room fail1"
					}
				} else {
					return "create chat room fail2"
				}
				fmt.Println(conn.RemoteAddr(), args["name"], "creat a chat room id:", id)
				return id
			} else {
				return "not login"
			}
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true // 释放db操作权限
	}()
}

// 加入聊天室
// 参数：name （加入者）& id （聊天室id）
// 返回值：ok
func Join(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		fmt.Println(dbState)
		<-dbState // 获取到db操作权限
		fmt.Println("1111111111")
		result := func() string {
			name := args["name"]
			if db.GetConnects(name) != nil {
				id := args["id"]
				chatRoom, err := db.GetChatRoom(id)
				if err != nil {
					fmt.Println(err)
					return "get chat room error"
				}
				if reflect.DeepEqual(chatRoom, room.ChatRoom{}) {
					return "no chat room"
				}
				r, err1 := db.GetRole(name)
				if err1 != nil {
					fmt.Println(err1)
					return "get role error"
				}
				if reflect.DeepEqual(r, role.Role{}) {
					return "no role"
				}
				if chatRoom.AddMember(name) && r.AddChatRoomId(id) {
					err := db.UpdateChatRoomRole(chatRoom, r)
					if err != nil {
						fmt.Println(err)
						return "join chat room fail1"
					}
				} else {
					return "join chat room fail2"
				}
				// 广播消息， 通知聊天室内的所有人，有新的成员加入
				chatRoom.Broadcast(name, name+"join chat room", db.GetAllConnects())
				fmt.Println(conn.RemoteAddr(), args["name"], "join chat room ", args["id"], " ...")
				return "ok"
			} else {
				return "not login"
			}
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true // 释放db操作权限
	}()

}

// 发消息
// 参数：name （消息发送者） & id（聊天室id）& msg（消息内容）
// 返回值： ok
func Say(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		<-dbState // 获取到db操作权限
		result := func() string {
			name := args["name"]
			if db.GetConnects(name) != nil {
				id := args["id"]
				chatRoom, err := db.GetChatRoom(id)
				if err != nil {
					fmt.Println(err)
					return "get chat room error"
				}
				if reflect.DeepEqual(chatRoom, room.ChatRoom{}) {
					return "no chat room"
				}
				r, err1 := db.GetRole(name)
				if err1 != nil {
					fmt.Println(err1)
					return "get role error"
				}
				if reflect.DeepEqual(r, role.Role{}) {
					return "no role"
				}
				msg := args["msg"]
				//  向所有在聊天室的在线人员广播消息
				chatRoom.Broadcast(name, msg, db.GetAllConnects())

				fmt.Println(conn.RemoteAddr(), args["name"], "in chat room ", args["id"], " say ", args["msg"], " ...")
				return "ok"
			} else {
				return "not login"
			}
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true // 释放db操作权限
	}()

}

//  退出
// 参数： name （退出者）& id (聊天室id)
// 返回值： ok
func Quit(dbState chan bool, conn net.Conn, args map[string]string) {
	go func() {
		<-dbState // 获取到db操作权限
		result := func() string {
			name := args["name"]
			if db.GetConnects(name) != nil {
				id := args["id"]
				chatRoom, err := db.GetChatRoom(id)
				if err != nil {
					return "get chat room error"
				}
				if reflect.DeepEqual(chatRoom, room.ChatRoom{}) {
					return "no chat room"
				}
				r, err1 := db.GetRole(name)
				if err1 != nil {
					return "get role error"
				}
				if reflect.DeepEqual(r, role.Role{}) {
					return "no role"
				}
				r.DeleteChatRoomId(id)

				chatRoom.DeleteMember(name)
				err = db.UpdateChatRoomRole(chatRoom, r)
				if err != nil {
					return "update role error"
				}
				//  向所有在聊天室的在线人员广播消息
				chatRoom.Broadcast(name, name+" quit chat room", db.GetAllConnects())
				fmt.Println(conn.RemoteAddr(), args["name"], "quit chat room", args["id"], " ...")

				return "ok"
			} else {
				return "not login"
			}
		}()
		conn.Write([]byte("result:" + result))
		dbState <- true // 释放db操作权限
	}()

}
