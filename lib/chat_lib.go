package lib

// id 计数
var idNum = uint32(1000001)

// Id 创建
func GetId() string {
	idNum++
	return string(idNum)
}
