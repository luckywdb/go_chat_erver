package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_chat_server/role"
	"go_chat_server/room"
	"net"
	"time"
)

var connects = make(map[string]net.Conn)

// redis 连接池
var RedisClient *redis.Pool

func init() {
	RedisClient = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			if _, err = c.Do("SELECT", 0); err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, nil
		},
		//DialContext:     nil,
		//TestOnBorrow:    nil,
		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: 1,
		//最大的激活连接数，表示同时最多有N个连接
		MaxActive: 10,
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 180 * time.Second,
		//Wait:            false,
		//MaxConnLifetime: 0,
	}
}

// 存储role对象
func UpdateRole(r role.Role) error {
	c := RedisClient.Get()
	defer c.Close()
	data, err := json.Marshal(&r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 存储数据
	if _, err = c.Do("SET", r.Name, string(data)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 获取role对象
var ErrRedisNil = errors.New("redigo: nil returned")

func GetRole(name string) (role.Role, error) {
	c := RedisClient.Get()
	defer c.Close()
	r := role.Role{}
	data, err := redis.String(c.Do("GET", name))
	if err != nil {
		fmt.Println(err)
		return r, err
	}
	if err = json.Unmarshal([]byte(data), &r); err != nil {
		fmt.Println(err)
		return r, err
	}
	return r, nil
}

// 删除role对象
func DeleteRole(name string) error {
	c := RedisClient.Get()
	defer c.Close()
	if _, err := c.Do("DEL", name); err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

// 存储 chat_room 对象
func UpdateChatRoom(cr room.ChatRoom) error {
	c := RedisClient.Get()
	defer c.Close()
	chatRoom, err := json.Marshal(&cr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 存储数据
	if _, err = c.Do("SET", cr.Id, string(chatRoom)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 获取chat_room对象
func GetChatRoom(id string) (room.ChatRoom, error) {
	c := RedisClient.Get()
	defer c.Close()
	cr := room.ChatRoom{}
	data, err := redis.String(c.Do("GET", id))
	if err != nil { // 当为redigo: nil returned 这个错误时表示，redis中没有数据
		fmt.Println(err)
		return cr, nil
	}

	if err = json.Unmarshal([]byte(data), &cr); err != nil {
		fmt.Println(err)
		return cr, err
	}
	return cr, nil
}

// 删除 chat_room对象
func DeleteChatRoom(id string) error {
	c := RedisClient.Get()
	defer c.Close()
	if _, err := c.Do("DEL", id); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 获取连接
func GetConnects(name string) net.Conn {
	return connects[name]
}

// 获取全部连接
func GetAllConnects() map[string]net.Conn {
	return connects
}

// 更新连接
func UpdateConnects(name string, con net.Conn, f func(string, net.Conn, map[string]net.Conn)) {
	f(name, con, connects)
}

// 新增连接
func AddConnects(name string, con net.Conn, connects map[string]net.Conn) {
	connects[name] = con
}

// 删除连接
func DeleteConnects(name string, _ net.Conn, connects map[string]net.Conn) {
	delete(connects, name)
}

// 存储聊天信息
//func UpdateMessage(key int64, value string) error {
//	c := RedisClient.Get()
//	defer c.Close()
//	// 存储数据
//	if _, err := c.Do("SET", key, value); err != nil {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}

//
