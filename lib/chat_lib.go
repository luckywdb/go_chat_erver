package lib

// id 计数
var idNum int

// Id 创建
func GetId() string {
	idNum++
	return string(rune(idNum))
}
