package room

import "time"

// 定义一个聊天室对象
// id 聊天室id
// title 聊天室名字
// creator 创建人（昵称）
// time 创建时间
// member 成员名
type ChatRoom struct {
	id      string
	title   string
	creator string
	time    int64
	member  []string
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

// 设置聊天室标题
func (cr *ChatRoom) SetTitle(newTitle string) {
	cr.title = newTitle
}

// 聊天室添加成员
func (cr *ChatRoom) AddMember(newMember string) bool {
	notExist := true
	for _, v := range cr.member {
		if v == newMember {
			notExist = false // 成员 已经存在 置为false
			break
		}
	}
	//如果成员不存在， 添加到成员列表中
	if notExist {
		cr.member = append(cr.member, newMember)
	}
	return notExist
}

// 聊天室删除成员
func (cr *ChatRoom) DeleteMember(name string) {
	for i, v := range cr.member {
		if v == name {
			// 找到聊天室id所在位置
			cr.member = append(cr.member[:i], cr.member[i+1:]...)
			break
		}
	}
}

// 聊天室广播消息
func (cr ChatRoom) broadcast() {

}
