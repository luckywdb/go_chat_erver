package actor

import (
	"fmt"
	"time"
)

type Actor interface {
	Start() // 启动一个actor
	Stop()  // 停止一个actor

	HandleCast(request interface{}, state interface{}) (newState interface{}, err error)                     // 处理异步消息
	HandleCall(request interface{}, state interface{}) (result interface{}, newState interface{}, err error) // 处理同步消息

	GetMailBox() chan interface{} // 获取自己的邮箱
}

// 一般性actor
type NormalActor struct {
	Name    string           // actor名字
	MailBox chan interface{} // actor的邮箱
	State   interface{}      // 当前actor的状态
}

// 同步消息结构体
type SyncMsg struct {
	Request interface{}      // 请求
	From    chan interface{} // 来源
}

// 异步消息结构体
type AsyncMsg struct {
	Request interface{} // 请求
}

// 发送一步消息
func Cast(receiver chan interface{}, request interface{}) {
	receiver <- request
}

// 发送同步消息, 默认5秒超时
func CallDefault(sender chan interface{}, receiver chan interface{}, request interface{}) interface{} {
	return CallTimeout(sender, receiver, request, 5)
}

// 发送同步消息, 自定义超时时间
func CallTimeout(sender chan interface{}, receiver chan interface{},
	request interface{}, timeout time.Duration) (result interface{}) {
	Cast(receiver, request) // 异步发送消息给接受者
	select {
	case result = <-sender: // 等待返回
		return result
	case <-time.After(time.Second * timeout): // 设置超时
		return "timeout"
	}
}

// 循环处理消息
func actorLoop(actor Actor, state interface{}) {
	mailBox := actor.GetMailBox()
	for {
		msgInterface := <-mailBox
		switch msg := msgInterface.(type) {
		case SyncMsg: // 同步处理
			rsp, newState, err := actor.HandleCall(msg.Request, state)
			if err != nil {
				fmt.Println(err)
				return
			}
			msg.From <- rsp
			state = newState
		case AsyncMsg: // 异步处理
			newState, err := actor.HandleCast(msg.Request, state)
			if err != nil {
				fmt.Println(err)
				return
			}
			state = newState
		}
	}
}
