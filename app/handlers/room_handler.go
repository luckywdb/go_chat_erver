package handlers

import (
	"github.com/golang/protobuf/proto"
	"go_chat_server/db"
	"go_chat_server/global"
	"go_chat_server/port"
	"go_chat_server/tcp"
)

func broadcast(request global.Request) (proto.Message, error) {
	rsp, err := func(request global.Request) (string, error) {
		params := request.Msg.(*port.MsgSay)
		room, err := db.GetChatRoom(params.Id)
		if err != nil {
			return "get chat room fail", nil
		}
		msg := params.Name + string(params.Msg)
		for _, name := range room.Member {
			conn := db.GetConnects(name)
			_ = tcp.MarshalMsg(request.RspName, &port.MsgSayReturn{
				Msg: msg,
			}, conn)
		}
		return "ok", nil
	}(request)
	return &port.MsgSayReturn{
		Msg: rsp,
	}, err
}
