package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-examples/project/config"
	"gin-examples/project/handlers"
	"gin-examples/project/middleware"
	"gin-examples/project/services"
	"gin-examples/project/utils"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 初始化服务：用户
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, []byte(cfg.JWT.Secret))

	postService := services.NewPostService(db)
	postHandler := handlers.NewPostHandler(postService)

	commentService := services.NewCommentService(db)
	commentHandler := handlers.NewCommentHandler(commentService)

	// ...

	// var json = jsoniter.Config{
	// 	TimeFormat: "2006-01-02 15:04:05", // 全局时间格式
	// 	EscapeHTML: false,
	// }.Froze()

	// // 替换Gin默认的JSON编码器
	// gin.DefaultJSONSerializer = &gin.JSONSerializer{
	// 	JSONEncoder: json,
	// }
	// // 替换表单绑定的JSON编码器
	// binding.Encoder = json

	// 创建 Gin 引擎
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)
		public.GET("/users/sta", userHandler.StatisticPostAuditStatus)

		public.GET("/posts", postHandler.ListPostAll)
		public.GET("/posts/:id", postHandler.GetPostById)
		public.POST("/posts/condition", postHandler.ListPostByCondition)
		public.GET("/posts/comment/number/max", postHandler.GetPostByMaxCommentNumber)

		public.GET("/comments/:postId", commentHandler.ListCommentByPostId)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	log.Printf("Router auth jwt secret:%s\n", cfg.JWT.Secret)
	protected.Use(middleware.Auth([]byte(cfg.JWT.Secret)))
	{
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)

		protected.POST("/posts/me", postHandler.CreatePost)
		protected.GET("/posts/me", postHandler.ListPost)
		protected.PUT("/posts/me", postHandler.UpdatePost)
		protected.DELETE("/posts/me/:id", postHandler.DeletePost)

		protected.POST("/comments", commentHandler.CreateComment)
		protected.DELETE("/comments/me/:postId/:id", commentHandler.DeleteComment)
	}

	return r
}
