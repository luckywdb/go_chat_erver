package message

import "time"

type Message struct {
	Time int64  `json:"time"`
	Who  string `json:"who"`
	Msg  string `json:"msg"`
}

func CreateMessage(who string, msg string) Message {
	return Message{
		time.Now().UnixNano(),
		who,
		msg,
	}
}
