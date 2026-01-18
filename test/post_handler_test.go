package test

import (
	"bytes"
	"encoding/json"
	"gin-examples/project/config"
	"gin-examples/project/router"
	"log"

	// "gin-examples/project/main"
	"gin-examples/project/models"
	"gin-examples/project/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 测试前初始化：设置Gin为测试模式（禁用颜色、日志精简）
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // 关键：测试模式避免输出干扰
	m.Run()
}

func TestUserHandler_Register(t *testing.T) {
	// 加载配置
	cfg := config.Load()
	// 创建测试数据库
	db := setupTestDB(t)
	router := setupTestHandlerRouter(cfg, db)

	// 测试完毕后，清空数据库
	defer config.CleanupDB(db)

	log.Print("*****************************")
	log.Print("API用户测试 START")
	log.Print("*****************************")

	req := models.CreateUserRequest{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, 200, w.Code)

	var response utils.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, response.Code)
}

func setupTestHandlerRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 定义路由
	r := router.SetupRouter(cfg, db)

	// 启动服务器
	// addr := cfg.Server.Host + ":" + cfg.Server.Port
	// log.Printf("Server starting on %s", addr)
	// if err := r.Run(addr); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
	return r
}
