package actor

import (
	"go_chat_server/global"
	"go_chat_server/tcp"
)

type RoleActor struct {
	NormalActor NormalActor
}

// 创建一个roleActor实例
func NewRoleActor(name string) RoleActor {
	return RoleActor{
		NormalActor: NormalActor{
			Name:    name,
			MailBox: make(chan interface{}, global.ChannelSize),
		},
	}
}

// 获取自己的邮箱
func (ra RoleActor) GetMailBox() chan interface{} {
	return ra.NormalActor.MailBox
}

// 开启actor
func (ra RoleActor) Start() {
	go actorLoop(ra, nil)
}

// 关闭 actor
func (ra RoleActor) Stop() {
	close(ra.NormalActor.MailBox)
}

// 处理异步消息
func (ra RoleActor) HandleCast(request interface{}, state interface{}) (newState interface{}, err error) {
	switch request := request.(type) {
	case global.Request:
		rsp, err := request.Handler(request)
		if err != nil {
			return state, err
		}
		err = tcp.MarshalMsg(request.RspName, rsp, request.Conn)
		return state, err
	}
	return state, nil
}

// 处理同步消息
func (ra RoleActor) HandleCall(request interface{},
	state interface{}) (result interface{}, newState interface{}, err error) {
	switch _ := request.(type) {

	}
	return nil, state, nil
}
