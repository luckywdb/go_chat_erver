package handlers

import (
	"github.com/golang/protobuf/proto"
	"go_chat_server/actor"
	"go_chat_server/db"
	"go_chat_server/global"
	"go_chat_server/lib"
	"go_chat_server/port"
	"go_chat_server/room"
)

func roomCreate(request global.Request) (proto.Message, error) {
	rsp, id, err := func(request global.Request) (string, string, error) {
		params := request.Msg.(*port.MsgCreate)
		roomId := lib.GetId()
		chatRoom := room.CreatorRoom(roomId, params.Title, params.Name)
		rol, err := db.GetRole(params.Name)
		if err == db.ErrRedisNil {
			return "no role", "", err
		} else if err != nil {
			return "get role error", "", err
		}
		rol.AddChatRoomId(roomId)

		if err = db.UpdateRole(rol); err != nil {
			return "update role error", "", err
		}

		if err = db.UpdateChatRoom(chatRoom); err != nil {
			return "update chat room error", "", err
		}
		roomActor := actor.NewRoomActor(roomId)            // 新建一个room actor
		global.ActorsChan[roomId] = roomActor.GetMailBox() // 存储room actor
		go roomActor.Start()                               // 启动room actor

		return "ok", roomId, nil
	}(request)
	return &port.MsgCreateReturn{
		Msg: rsp,
		Id:  id,
	}, err

}
