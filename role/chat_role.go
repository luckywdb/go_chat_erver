package role

import (
	"time"
)

// 聊天用户对象

// name 用户名
// createTime 创建时间
// loginTime 登陆时间
// chatRooms 拥有的聊天室id
type Role struct {
	Name       string   `json:"name"`
	CreateTime int64    `json:"createTime,string"`
	LoginTime  int64    `json:"loginTime,string"`
	ChatRooms  []string `json:"chatRooms"`
}

// 创建用户
func CreateRole(name string) Role {
	now := time.Now().Unix()
	return Role{name, now, now, []string{}}
}

// 添加聊天室id
// 成功返回 true， 不成功返回 false
func (r *Role) AddChatRoomId(id string) bool {
	notExist := true
	for _, v := range r.ChatRooms {
		if v == id {
			notExist = false // id 已经存在 置为false
			break
		}
	}
	//如果id不存在， 添加到id列表中
	if notExist {
		r.ChatRooms = append(r.ChatRooms, id)
	}
	return notExist
}

// 删除聊天室id
func (r *Role) DeleteChatRoomId(id string) {
	for i, v := range r.ChatRooms {
		if v == id {
			// 找到聊天室id所在位置
			r.ChatRooms = append(r.ChatRooms[:i], r.ChatRooms[i+1:]...)
			break
		}
	}
}
