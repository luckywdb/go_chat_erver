package handlers

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go_chat_server/actor"
	"go_chat_server/db"
	"go_chat_server/global"
	"go_chat_server/port"
	"go_chat_server/role"
	"go_chat_server/tcp"
	"time"
)

func init() {
	//tcp.Register("MsgLogin", func() interface{} {
	//	return &port.MsgLogin{}
	//}, login)
	//tcp.Register("MsgLogout", func() interface{} {
	//	return &port.MsgLogout{}
	//}, logout)
	tcp.Register("MsgLogin", "MsgLoginReturn", login)
	tcp.Register("MsgLogout", "MsgLogoutReturn", logout)
	tcp.Register("MsgCreate", "MsgCreateReturn", create)
	tcp.Register("MsgJoin", "MsgJoinReturn", join)
	tcp.Register("MsgSay", "MsgSayReturn", say)
	tcp.Register("MsgQuit", "MsgQuitReturn", quit)

}

// 登录
func login(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgLogin)
		name := params.Name
		r, err := db.GetRole(name)
		if err == db.ErrRedisNil { // 没有 role 需创建
			r = role.CreateRole(name) // 创建用户
		} else if err != nil {
			return "get role error ", err
		}
		// 2. 用户登陆
		r.LoginTime = time.Now().Unix()
		// 更新role对象
		if err = db.UpdateRole(r); err != nil {
			return "update role error", err
		}
		fmt.Println(request.Conn.RemoteAddr(), name, "login ok ...")
		db.UpdateConnects(name, request.Conn, db.AddConnects)
		return "ok", err
	}(request)
	return &port.MsgLoginReturn{
		Msg: rsp,
	}, err
}

// 登出
func logout(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgLogout)
		name := params.Name
		if db.GetConnects(name) != nil { // 在线
			rol, err := db.GetRole(name)
			if err == db.ErrRedisNil { // 没有 role 需创建
				return "no role", err
			} else if err != nil {
				return "get role error", err
			}
			for _, chatRoomId := range rol.ChatRooms { // 遍历用户所有的聊天室

				if _, err := db.GetChatRoom(chatRoomId); err != nil {
					return "get chat room error", err
				}
				//  广播
				roleBroadcast(request, name, chatRoomId, []byte("offline"))
			}
			db.UpdateConnects(name, request.Conn, db.DeleteConnects)
			fmt.Println(request.Conn.RemoteAddr(), name, "logout ok ...")
			return "logout", nil
		} else {
			return "not login", nil
		}
	}(request)
	return &port.MsgLogoutReturn{
		Msg: rsp,
	}, err
}

// 创建房间
func create(request global.Request) (proto.Message, error) {
	params := request.Msg.(*port.MsgCreate)
	name := params.Name
	if _, err := db.GetRole(name); err == db.ErrRedisNil {
		return &port.MsgCreateReturn{
			Msg: "no role",
			Id:  "",
		}, err
	} else if err != nil {
		return &port.MsgCreateReturn{
			Msg: "get role error",
			Id:  "",
		}, err
	}
	roomManagerActorChan := global.ActorsChan[global.RoomManagerActor]
	temporaryChan := make(chan interface{})
	defer close(temporaryChan)
	request.Handler = roomCreate
	rsp := actor.CallDefault(temporaryChan, roomManagerActorChan, actor.SyncMsg{
		Request: request,
		From:    temporaryChan,
	})
	rsp, ok := rsp.(proto.Message)
	if !ok {
		return &port.MsgCreateReturn{
			Msg: "timeout",
			Id:  "",
		}, nil
	}
	return proto.MessageV1(rsp), nil
}

// 加入房间
func join(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgJoin)
		if _, err := db.GetRole(params.Name); err == db.ErrRedisNil {
			return "no role", err
		} else if err != nil {
			return "get role error", err
		}
		id := params.Id
		chatRoom, err := db.GetChatRoom(id)
		if err == db.ErrRedisNil {
			return "no chat room", err
		} else if err != nil {
			return "get chat room error", err
		}
		if func() bool {
			for _, v := range chatRoom.Member {
				if v == params.Name {
					return false
				}
			}
			return true
		}() {
			//  加入聊天室
			chatRoom.AddMember(params.Name)

			if err := db.UpdateChatRoom(chatRoom); err != nil {
				return "join chat room fail", err
			}
			//  广播
			roleBroadcast(request, params.Name, id, []byte("join"))
			return "ok", nil
		} else {
			return "joined ", nil // 已经加入过了
		}
	}(request)
	return &port.MsgJoinReturn{
		Msg: rsp,
	}, err

}

//  广播
func roleBroadcast(request global.Request, name string, roomId string, msg []byte) {
	if roomChan, ok := global.ActorsChan[roomId]; ok {
		actor.Cast(roomChan, actor.AsyncMsg{
			Request: global.Request{
				Msg: port.MsgSay{
					Name: name,
					Id:   roomId,
					Msg:  msg,
				},
				Handler: broadcast,
				Conn:    request.Conn,
			},
		})
	}
}

// 说话 发消息
func say(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgSay)

		if _, err := db.GetChatRoom(params.Id); err == db.ErrRedisNil {
			return "no chat room", err
		} else if err != nil {
			return "get chat room error", err
		}

		if _, err := db.GetRole(params.Name); err == db.ErrRedisNil {
			return "no role", err
		} else if err != nil {
			return "get role error", err
		}
		// 向所有在聊天室的在线人员广播消息
		roleBroadcast(request, params.Name, params.Id, params.Msg)
		return "ok", nil
	}(request)
	return &port.MsgSayReturn{
		Msg: rsp,
	}, err
}

// 退出聊天室
func quit(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgQuit)
		chatRoom, err := db.GetChatRoom(params.Id)
		if err == db.ErrRedisNil {
			return "no chat room", err
		} else if err != nil {
			return "get chat room error", err
		}
		rol, err := db.GetRole(params.Name)
		if err == db.ErrRedisNil {
			return "no role", err
		} else if err != nil {
			return "get role error", err
		}

		//  退出聊天室后的操作
		chatRoom.DeleteMember(params.Name)
		rol.DeleteChatRoomId(params.Id)

		if err = db.UpdateChatRoom(chatRoom); err != nil {
			return "update chat room fail", err
		}

		if err = db.UpdateRole(rol); err != nil {
			return "update role fail", err
		}
		roleBroadcast(request, params.Name, params.Id, []byte("quit"))
		return "ok", nil
	}(request)
	return &port.MsgQuitReturn{
		Msg: rsp,
	}, err
}
