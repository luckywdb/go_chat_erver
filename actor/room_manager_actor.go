package actor

import (
	"go_chat_server/global"
)

type RoomManagerActor struct {
	NormalActor NormalActor
}

// 创建一个roleActor实例
func NewRoomManagerActor(name string) RoomManagerActor {
	return RoomManagerActor{
		NormalActor: NormalActor{
			Name:    name,
			MailBox: make(chan interface{}, global.ChannelSize),
		},
	}
}

// 获取自己的邮箱
func (rma RoomManagerActor) GetMailBox() chan interface{} {
	return rma.NormalActor.MailBox
}

// 开启actor
func (rma RoomManagerActor) Start() {
	go actorLoop(rma, nil)
}

// 关闭 actor
func (rma RoomManagerActor) Stop() {
	close(rma.NormalActor.MailBox)
}

// 处理异步消息
func (rma RoomManagerActor) HandleCast(request interface{}, state interface{}) (newState interface{}, err error) {
	switch _ := request.(type) {

	}
	return state, nil
}

// 处理同步消息
func (rma RoomManagerActor) HandleCall(request interface{},
	state interface{}) (result interface{}, newState interface{}, err error) {

	switch request := request.(type) {
	case global.Request:
		rsp, err := request.Handler(request)
		return rsp, state, err
	}
	return nil, state, nil
}
