package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"time"
)

// redis 连接池
var RedisClient *redis.Pool

func init() {
	RedisClient = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			_, err = c.Do("SELECT", 0)
			if err != nil {
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

type Role struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	aa   []string
}

func main() {
	//c := RedisClient.Get()
	//defer c.Close()
	//data, err := redis.String(c.Do("GET", "abc"))
	//if err != nil {
	//	fmt.Println("error:", err)
	//	return
	//}
	//fmt.Println("11111", data)

	m := make(map[string]net.Conn)
	fmt.Println(m["a"] == nil)
}
