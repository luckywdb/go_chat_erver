package port

//
//import (
//	"github.com/golang/protobuf/proto"
//	"go_chat_server/actor"
//	"go_chat_server/global"
//	"go_chat_server/tcp"
//	"net"
//)
//
//// 用与和客户端通信
//
//// 登陆
//// 参数：name （用户名）
//// 返回值： ok|error
//func Login(conn net.Conn, data []byte) error {
//	loginMsg := &MsgLogin{}
//
//	if err := proto.Unmarshal(data, loginMsg); err != nil {
//		return err
//	}
//	name := loginMsg.Name
//	roleActor := actor.NewRoleActor(name) // 创建一个role actor
//	global.ActorsChan[name] = roleActor.GetMailBox()
//	roleActor.Start()
//
//	actor.Cast(roleActor.GetMailBox(), actor.AsyncMsg{
//		Request: global.Login{
//			LoginMsg: loginMsg,
//			Conn:     conn,
//		},
//	})
//	return nil
//}
//
//// 登出
//// 参数：name （用户名）
//// 返回值： ok|error
//func Logout(conn net.Conn, data []byte) error {
//	logoutMsg := &MsgLogout{}
//	if err := proto.Unmarshal(data, logoutMsg); err != nil {
//		return err
//	}
//	name := logoutMsg.Name
//	if roleActorChan, ok := global.ActorsChan[name]; ok {
//		actor.Cast(roleActorChan, actor.AsyncMsg{
//			Request: global.Logout{
//				LogoutMsg: logoutMsg,
//				Conn:      conn,
//			},
//		})
//	}
//	return nil
//}
//
//// 创建聊天室
//// 参数: name (创建者)
//// 返回值：聊天室id
//func Create(conn net.Conn, data []byte) error {
//	createMsg := &MsgCreate{}
//	if err := proto.Unmarshal(data, createMsg); err != nil {
//		return err
//	}
//	castMsg(createMsg.Name, global.CreateRoom{
//		CreateMsg: createMsg,
//		Conn:      conn,
//	})
//	return nil
//}
//
//// 加入聊天室
//// 参数：name （加入者）& id （聊天室id）
//// 返回值：ok
//func Join(conn net.Conn, data []byte) error {
//	joinMsg := &MsgJoin{}
//	if err := proto.Unmarshal(data, joinMsg); err != nil {
//		return err
//	}
//	castMsg(joinMsg.Name, global.JoinRoom{
//		JoinMsg: joinMsg,
//		Conn:    conn,
//	})
//	return nil
//}
//
//// 发消息
//// 参数：name （消息发送者） & id（聊天室id）& msg（消息内容）
//// 返回值： ok
//func Say(conn net.Conn, data []byte) error {
//	sayMsg := &MsgSay{}
//	if err := proto.Unmarshal(data, sayMsg); err != nil {
//		return err
//	}
//	castMsg(sayMsg.Name, global.Say{
//		SayMsg: sayMsg,
//		Conn:   conn,
//	})
//	return nil
//}
//
////  退出
//// 参数： name （退出者）& id (聊天室id)
//// 返回值： ok
//func Quit(conn net.Conn, data []byte) error {
//	quitMsg := &MsgQuit{}
//	if err := proto.Unmarshal(data, quitMsg); err != nil {
//		return err
//	}
//	castMsg(quitMsg.Name, global.Quit{
//		QuitMsg: quitMsg,
//		Conn:    conn,
//	})
//	return nil
//}
//
//type RequestMsg interface {
//	GetName() string
//}
//
//func Process(conn net.Conn, data []byte, factory tcp.Factory, handler tcp.Handler) error {
//	msg := factory()
//	if err := proto.Unmarshal(data, msg); err != nil {
//		return err
//	}
//	if msg, ok := msg.(*MsgLogin); ok {
//		startRoleActor(msg.Name)
//	}
//	castMsg(msg.(RequestMsg).GetName(), global.Request{
//		Msg:     msg,
//		Handler: handler,
//		Conn:    conn,
//	})
//	return nil
//}
//
//func castMsg(actorName string, request interface{}) {
//	roleActorChan := global.ActorsChan[actorName]
//	actor.Cast(roleActorChan, actor.AsyncMsg{
//		Request: request,
//	})
//}
//
//func startRoleActor(name string) {
//	roleActor := actor.NewRoleActor(name) // 创建一个role actor
//	global.ActorsChan[name] = roleActor.GetMailBox()
//	roleActor.Start()
//}
