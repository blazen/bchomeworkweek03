package test

import (
	"gin-examples/project/config"
	"gin-examples/project/models"
	"gin-examples/project/services"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostService_CreatePost(t *testing.T) {
	// 创建测试数据库
	db := setupTestDB(t)

	// 测试完毕后，清空数据库
	defer config.CleanupDB(db)

	log.Print("*****************************")
	log.Print("文章测试 START")
	log.Print("*****************************")

	// 初始化测试数据
	user, err := setupTestServicePostData(db)

	// 本次测试
	postService := services.NewPostService(db)
	req := models.CreatePostRequest{
		Title:   "hello",
		Content: "hello world",
	}
	post, err := postService.CreatePost(user.ID, req)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Title, post.Title)
	assert.Equal(t, req.Content, post.Content)
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// 创建测试数据库连接
	db := config.NewDB("blog_test.db")
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})

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
	return db
}

func setupTestServicePostData(db *gorm.DB) (*models.User, error) {
	// 测试文章，预先创建初始测试账户，显式设置 ID 为 1 和 2，确保每次测试都使用相同的 ID
	userService := services.NewUserService(db)
	req := models.CreateUserRequest{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
	}
	user, err := userService.CreateUser(req)
	if err != nil {
		log.Fatalf("Failed to create test user:admin: %v", err)
	}
	log.Printf("User admin create success userId:%d, username=%s", user.ID, user.Username)
	return user, nil
}
