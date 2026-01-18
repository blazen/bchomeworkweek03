package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GencOrderNo() string {
	timestamp := time.Now().UnixMilli()
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9000) + 1000
	return fmt.Sprintf("%d%d", timestamp, randomNum)
}

// type Generator struct {
// 	machineID  uint64
// 	sequence   uint64
// 	lastMillis int64
// }

// var globalGenerator = &Generator{
// 	machineID: 1, // 示例：机器ID，实际可从配置读取
// }

// func (g *Generator) GencOrderNo() string {
// 	now := time.Now().UnixMilli()
// 	// 1. 原子递增序列
// 	atomic.AddUint64(&g.sequence, 1)

// 	// 2. 如果跨毫秒，重置序列
// 	if now != g.lastMillis {
// 		atomic.StoreUint64(&g.sequence, 1)
// 		atomic.StoreInt64(&g.lastMillis, now)
// 	}

// 	// 3. 拼接规则：时间戳(13位) + 机器ID(2位) + 序列(3位) = 18位
// 	return fmt.Sprintf("%d%02d%03d", now, g.machineID, g.sequence%1000)
// }
