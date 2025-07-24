package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xhy/blog-api/config"
	"github.com/xhy/blog-api/models"
	"github.com/xhy/blog-api/routes"
	"github.com/xhy/blog-api/utils"
)

func main() {
	// 获取配置
	cfg := config.GetConfig()

	// 设置JWT配置
	utils.SetJWTSecret(cfg.JWT.Secret)
	utils.SetJWTDuration(cfg.JWT.ExpiresIn)

	// 初始化数据库连接
	config.InitDB()

	// 删除旧表（如果存在）
	//log.Println("删除旧表...")
	//config.DB.Exec("DROP TABLE IF EXISTS comments")
	//config.DB.Exec("DROP TABLE IF EXISTS posts")
	//config.DB.Exec("DROP TABLE IF EXISTS users")
	//log.Println("旧表删除完成")

	// 自动迁移数据库模型
	log.Println("开始数据库迁移...")
	err := config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	// 初始化示例数据
	//config.SeedData()

	// 创建Gin实例
	router := gin.Default()

	// 设置路由
	routes.SetupRoutes(router)

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%s", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}
