package role

import "time"

// 聊天用户对象

// name 用户名
// createTime 创建时间
// loginTime 登陆时间
// chatRooms 拥有的聊天室id
type Role struct {
	name       string
	createTime int64
	loginTime  int64
	chatRooms  []string
}

// 创建用户
func CreateRole(name string) Role {
	now := time.Now().Unix()
	return Role{name, now, now, []string{}}
}

// 设置登陆时间
func (r *Role) SetLoginTime(time int64) {
	r.loginTime = time
}

// 添加聊天室id
// 成功返回 true， 不成功返回 false
func (r *Role) AddChatRoomId(id string) bool {
	notExist := true
	for _, v := range r.chatRooms {
		if v == id {
			notExist = false // id 已经存在 置为false
			break
		}
	}
	//如果id不存在， 添加到id列表中
	if notExist {
		r.chatRooms = append(r.chatRooms, id)
	}
	return notExist
}

// 删除聊天室id
func (r *Role) DeleteChatRoomId(id string) {
	for i, v := range r.chatRooms {
		if v == id {
			// 找到聊天室id所在位置
			r.chatRooms = append(r.chatRooms[:i], r.chatRooms[i+1:]...)
			break
		}
	}
}
