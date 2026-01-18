package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// LocalTime 自定义时间类型（格式化为 YYYY-MM-DD HH:MM:SS）
type Time1 time.Time

// 常量：时间格式化模板
const timeFormat = "2006-01-02 15:04:05"

// MarshalJSON 实现 JSON 序列化接口（返回格式化后的字符串）
func (t Time1) MarshalJSON() ([]byte, error) {
	// 将自定义时间类型转为 time.Time
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("null"), nil // 空时间返回 null
	}
	// 格式化时间为字符串并包裹双引号（JSON 字符串格式）
	return json.Marshal(tt.Format(timeFormat))
}

// UnmarshalJSON 实现 JSON 反序列化接口（接收格式化字符串转为时间）
func (t *Time1) UnmarshalJSON(data []byte) error {
	// 空值处理
	if string(data) == "null" {
		*t = Time1(time.Time{})
		return nil
	}
	// 解析字符串为时间
	tt, err := time.Parse(`"`+timeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	*t = Time1(tt)
	return nil
}

// Value 实现 driver.Valuer 接口（数据库存储时的格式）
func (t Time1) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return tt.Format(timeFormat), nil
}

// Scan 实现 sql.Scanner 接口（从数据库读取时解析为 LocalTime）
func (t *Time1) Scan(v interface{}) error {
	if v == nil {
		*t = Time1(time.Time{})
		return nil
	}
	switch v := v.(type) {
	case time.Time:
		*t = Time1(v)
	case []byte:
		tt, err := time.Parse(timeFormat, string(v))
		if err != nil {
			return err
		}
		*t = Time1(tt)
	case string:
		tt, err := time.Parse(timeFormat, v)
		if err != nil {
			return err
		}
		*t = Time1(tt)
	default:
		return fmt.Errorf("不支持的时间类型：%T", v)
	}
	return nil
}

// String 实现 Stringer 接口（打印时返回格式化字符串）
func (t Time1) String() string {
	return time.Time(t).Format(timeFormat)
}
