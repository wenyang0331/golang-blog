package routes

import (
	"homework_blog/config"
	"homework_blog/database"
	"homework_blog/handlers"
	"homework_blog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置路由
func SetupRoutes(config *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 使用全局中间件（无需认证）
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 创建控制器实例
	userHandler := handlers.NewUserHandler(database.NewUserOperator(db), []byte(config.JWT.Secret))
	postHandler := handlers.NewPostHandler(database.NewPostOperator(db), []byte(config.JWT.Secret))
	commentHandler := handlers.NewCommentHandler(db, []byte(config.JWT.Secret))

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		public := api.Group("/public")
		{
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)

			public.GET("/posts", postHandler.GetPosts)
			public.GET("/posts/:id", postHandler.GetPost)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.Auth([]byte(config.JWT.Secret)))
		{
			// 用户信息
			authenticated.GET("/profile", userHandler.GetProfile)

			// 文章相关路由
			posts := authenticated.Group("/posts")
			{
				posts.POST("", postHandler.CreatePost)
				posts.PUT("/:id", postHandler.UpdatePost)
				posts.DELETE("/:id", postHandler.DeletePost)
			}

			// 评论相关路由
			comments := authenticated.Group("/posts/:post_id/comments")
			{
				comments.POST("", commentHandler.CreateComment)
			}
		}

		// 评论公开路由（单独分组避免路由冲突）
		comments := api.Group("/comments")
		{
			comments.GET("/post/:post_id", commentHandler.GetComments)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	return r
}
