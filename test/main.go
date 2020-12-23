package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go_chat_server/port"
	"net"
	"reflect"
)

func main() {
}
func testSlice() {
	fmt.Println([]int{})
}
func test() {
	msgLogin := &port.MsgLogin{
		Name: "abc",
	}

	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	defer conn.Close()
	if err != nil {
		println(err)
	}

	data, err := proto.Marshal(msgLogin)
	if err != nil {
		println(err)
	}
	msg := &port.Msg{
		MsgName: "MsgLogin",
		Data:    data,
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		println(err)
	}
	_, err = conn.Write(msgBytes)
	if err != nil {
		println(err)

	}
}

type user struct {
	name string
}

func testReflect() {
	//u := new(user)
	u := user{
		name: "abc",
	}
	fmt.Println(u)
	fmt.Println("typeOf", reflect.TypeOf(u))
	fmt.Println("kind", reflect.TypeOf(u).Kind() == reflect.Ptr)
	fmt.Println("element", reflect.TypeOf(u).Elem().Name())
}
