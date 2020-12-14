package room

import (
	"net"
	"time"
)

// 定义一个聊天室对象
// id 聊天室id
// title 聊天室名字
// creator 创建人（昵称）
// time 创建时间
// member 成员名
type ChatRoom struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Creator string   `json:"creator"`
	Time    int64    `json:"time,string"`
	Member  []string `json:",string"`
}

// 创建聊天室
func CreatorRoom(id string, title string, creator string) ChatRoom {

	return ChatRoom{
		id,
		title,
		creator,
		time.Now().Unix(),
		[]string{creator},
	}
}

// 聊天室添加成员
func (cr *ChatRoom) AddMember(newMember string) bool {
	notExist := true
	for _, v := range cr.Member {
		if v == newMember {
			notExist = false // 成员 已经存在 置为false
			break
		}
	}
	//如果成员不存在， 添加到成员列表中
	if notExist {
		cr.Member = append(cr.Member, newMember)
	}
	return notExist
}

// 聊天室删除成员
func (cr *ChatRoom) DeleteMember(name string) {
	for i, v := range cr.Member {
		if v == name {
			// 找到聊天室id所在位置
			cr.Member = append(cr.Member[:i], cr.Member[i+1:]...)
			break
		}
	}
}

// 聊天室广播消息
// name : 谁说话
// msg : 消息
func (cr ChatRoom) Broadcast(name string, msg string, conns map[string]net.Conn) {
	for _, who := range cr.Member {
		if con, ok := conns[who]; ok {
			con.Write([]byte("say:" + name + "_" + msg))
		}
	}
}
