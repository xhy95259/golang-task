package config

import (
	"log"

	"github.com/xhy/blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希处理（在这里复制一份避免循环导入）
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// SeedData 初始化一些示例数据
func SeedData() {
	// 检查是否已有用户数据
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	// 如果没有用户数据，则创建示例用户
	if userCount == 0 {
		log.Println("添加示例用户数据...")

		// 创建管理员用户
		adminPassword, _ := hashPassword("admin123")
		admin := models.User{
			Username: "admin",
			Password: adminPassword,
			Email:    "admin@example.com",
		}
		DB.Create(&admin)

		// 创建普通用户
		userPassword, _ := hashPassword("user123")
		user := models.User{
			Username: "user",
			Password: userPassword,
			Email:    "user@example.com",
		}
		DB.Create(&user)

		log.Println("示例用户创建成功")
	}

	// 检查是否已有文章数据
	var postCount int64
	DB.Model(&models.Post{}).Count(&postCount)

	// 如果没有文章数据，则创建示例文章
	if postCount == 0 {
		log.Println("添加示例文章数据...")

		// 管理员发布的文章
		post1 := models.Post{
			Title:   "欢迎使用个人博客系统",
			Content: "这是一个基于Go语言、Gin框架和GORM库开发的个人博客系统。您可以使用它来发布文章、评论等。",
			UserID:  1, // admin用户ID
		}
		DB.Create(&post1)

		// 普通用户发布的文章
		post2 := models.Post{
			Title:   "Go语言学习笔记",
			Content: "Go是一种开源编程语言，它使构建简单、可靠和高效的软件变得容易。本文将分享一些Go语言的基础知识和学习心得。",
			UserID:  2, // user用户ID
		}
		DB.Create(&post2)

		log.Println("示例文章创建成功")
	}

	// 检查是否已有评论数据
	var commentCount int64
	DB.Model(&models.Comment{}).Count(&commentCount)

	// 如果没有评论数据，则创建示例评论
	if commentCount == 0 {
		log.Println("添加示例评论数据...")

		// 用户对管理员文章的评论
		comment1 := models.Comment{
			Content: "这个博客系统非常好用！",
			UserID:  2, // user用户ID
			PostID:  1, // 管理员的文章ID
		}
		DB.Create(&comment1)

		// 管理员对用户文章的评论
		comment2 := models.Comment{
			Content: "写得很好，继续加油！",
			UserID:  1, // admin用户ID
			PostID:  2, // 用户的文章ID
		}
		DB.Create(&comment2)

		log.Println("示例评论创建成功")
	}

	log.Println("数据初始化完成")
}
