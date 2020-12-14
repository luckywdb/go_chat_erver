package lib

import (
	"strconv"
)

// id 计数
var idNum = 0

// Id 创建
func GetId() string {
	idNum++
	return strconv.Itoa(idNum)
}
