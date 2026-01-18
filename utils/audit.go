package utils

import (
	"fmt"
	"strings"
)

// 通过结构体实现子包
type AUDIT struct{}

var Audit = &AUDIT{}

// 验证输入是否合法
func (*AUDIT) Check(input string) bool {
	if strings.Contains(input, "金融") {
		fmt.Printf("输入非法 input=%s \n", input)
		return false
	}
	fmt.Printf("输入合法 input=%s \n", input)
	return true
}
