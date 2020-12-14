package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	CONNECT = iota // 刚刚建立连接
	LOGIN          // 登陆但是不再聊天室中
	JOIN           // 在聊天室中
)

var state = 0 // 状态码
func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8088")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	///预先准备消息缓冲区
	buffer := make([]byte, 1024)
	//准备命令行标准输入
	reader := bufio.NewReader(os.Stdin)
	for {
		lineBytes, _, _ := reader.ReadLine()
		conn.Write(lineBytes)
		msgLen, err := conn.Read(buffer)
		if msgLen == 0 || err != nil {
			continue
		}
		result := string(buffer[0:msgLen])
		fmt.Println(result)
		if result == "logout" {
			break
		}
	}
}
