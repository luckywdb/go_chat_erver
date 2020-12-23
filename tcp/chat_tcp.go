package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go_chat_server/actor"
	"go_chat_server/global"
	"go_chat_server/port"
	"io"
	"net"
)

var routes = make(map[string]global.Handler) // {"消息名字":Handler}
//type Factory func() interface{}
//
//var factories = make(map[string]Factory)

var rsp = make(map[string]string)

func Register(name string, rspName string, handler global.Handler) {
	routes[name] = handler
	rsp[name] = rspName
}

// 等待客户端连接
func WaitConnect(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error = ", err)
			return
		}
		fmt.Println(conn.RemoteAddr(), "connect success ...")
		go handleConn(conn)
	}
}

// 处理tcp链接消息
// 一次处理分三个步骤
// 1. 首先从消息中拿出 4个字节的，数据 head； head中存储的是，本次客户端给我的实际数据的长度
// 2. 根据从head中拿出来的长度， 在从消息中拿出对应长度的数据出来，对实际的数据做处理
// 3. 处理数据后，返还客户端对应的消息
// 一直循环做步骤1，2 ，3
func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := UnmarshalMsg(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 消息直接丢给 role actor 处理
		handler := routes[msg.MsgName]
		if err = process(conn, msg.Data, handler); err != nil {
			fmt.Println(err)
			return
		}
	}
}

type RequestMsg interface {
	GetName() string
}

func process(conn net.Conn, data []byte, handler global.Handler) error {
	var msg proto.Message
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}
	msg, ok := msg.(*port.MsgLogin)
	name := msg.(RequestMsg).GetName()
	if ok {
		roleActor := actor.NewRoleActor(name) // 创建一个role actor
		global.ActorsChan[name] = roleActor.GetMailBox()
		roleActor.Start()
	}
	roleChan := global.ActorsChan[name]
	actor.Cast(roleChan, actor.AsyncMsg{
		Request: global.Request{
			Msg:     msg,
			Conn:    conn,
			Handler: handler,
			RspName: rsp[name],
		},
	})
	return nil
}

// 反序列化消息
func UnmarshalMsg(conn net.Conn) (*port.Msg, error) {
	headByte := make([]byte, global.HeadLength) // 头部

	if _, err := io.ReadFull(conn, headByte); err != nil {
		return nil, err
	}

	bytesBuffer := bytes.NewBuffer(headByte)
	var bodyLength int32
	// 从头部， 解析出body的长度
	if err := binary.Read(bytesBuffer, binary.BigEndian, &bodyLength); err != nil {
		return nil, err
	}
	bodyByte := make([]byte, bodyLength) // body

	if _, err := io.ReadFull(conn, bodyByte); err != nil {
		return nil, err
	}
	msg := &port.Msg{}

	if err := proto.Unmarshal(bodyByte, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func MarshalMsg(msgName string, rsp proto.Message, conn net.Conn) error {
	data, err := proto.Marshal(rsp)
	if err != nil {
		return err
	}
	msgReturn := &port.MsgReturn{
		MsgName: msgName,
		Data:    data,
	}
	msgReturnBytes, err := proto.Marshal(msgReturn)
	if err != nil {
		return err
	}
	bytesBuffer := bytes.NewBuffer([]byte{})
	// int32 占4个byte
	if err = binary.Write(bytesBuffer, binary.BigEndian, int32(len(msgReturnBytes))); err != nil {
		return err
	}

	if _, err = conn.Write(append(bytesBuffer.Bytes(), msgReturnBytes...)); err != nil {
		return err
	}
	return nil
}
