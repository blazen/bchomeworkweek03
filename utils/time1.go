package utils

import (
	"fmt"
	"time"
)

// 1. 定义自定义时间类型（基于time.Time）
type Time1 time.Time

// 2. 实现json.Marshaler接口，自定义序列化逻辑
func (ct Time1) MarshalJSON() ([]byte, error) {
	// 转换为time.Time类型
	t := time.Time(ct)
	// 定义目标格式（Go的时间模板固定为2006-01-02 15:04:05）
	format := t.Format("2006-01-02 15:04:05")
	// 拼接JSON字符串（需加引号，否则JSON格式错误）
	return []byte(fmt.Sprintf("\"%s\"", format)), nil
}

// 可选：实现UnmarshalJSON，支持反序列化（前端传字符串转CustomTime）
func (ct *Time1) UnmarshalJSON(data []byte) error {
	// 去掉字符串两端的引号
	str := string(data)
	str = str[1 : len(str)-1]
	// 解析为time.Time
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	// 赋值给自定义类型
	*ct = Time1(t)
	return nil
}
