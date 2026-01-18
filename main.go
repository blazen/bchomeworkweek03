package main

import (
	"log"

	"gin-examples/project/config"
	"gin-examples/project/models"
	"gin-examples/project/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db := config.NewDB("blog.db")

	/**
		自动迁移
	- 自动创建表（如果不存在）
	- 添加新列（如果结构体有新字段）
	- 创建索引（根据标签）
	- **不会删除已存在的列**
	- **不会修改现有数据**
	- **不会删除索引**
	- 追加 uniqueIndex 时，若是数据中存在重复的数据，那么启动失败。
	*/
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 定义路由
	r := router.SetupRouter(cfg, db)

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
