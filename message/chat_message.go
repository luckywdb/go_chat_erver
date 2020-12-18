package message

type Message struct {
	Time int64  `json:"time"`
	Who  string `json:"who"`
	Msg  []byte `json:"msg"`
}
