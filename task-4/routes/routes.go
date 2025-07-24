package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xhy/blog-api/controllers"
	"github.com/xhy/blog-api/middleware"
)

// SetupRoutes 配置路由
func SetupRoutes(router *gin.Engine) {
	// 中间件
	router.Use(middleware.LoggerMiddleware())

	// 公开路由
	public := router.Group("/api")
	{
		// 用户认证
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		// 文章相关
		public.GET("/posts", controllers.GetPosts)
		public.GET("/posts/:id", controllers.GetPost)
		public.GET("/posts/:id/comments", controllers.GetComments)
	}

	// 需要认证的路由
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 文章相关
		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		// 评论相关
		protected.POST("/posts/:id/comments", controllers.CreateComment)
	}
} 