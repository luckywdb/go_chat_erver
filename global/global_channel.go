package global

import (
	"github.com/golang/protobuf/proto"
	"net"
)

const ChannelSize = 1024 // actor mailbox缓冲大小

// 消息头部的长度
const HeadLength = 4

const RoomManagerActor = "RoomManagerActor"

// 路由
type Handler func(Request) (proto.Message, error)

var ActorsChan map[string]chan interface{} // 全局变量用于存储 go 的channel
// 请求
type Request struct {
	Msg     interface{}
	Handler Handler
	Conn    net.Conn
	RspName string
}

// 返回
type Response struct {
	MsgName string
	Msg     interface{}
}
