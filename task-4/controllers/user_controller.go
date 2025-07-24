package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhy/blog-api/config"
	"github.com/xhy/blog-api/models"
	"github.com/xhy/blog-api/utils"
)

// Register 用户注册
func Register(c *gin.Context) {
	var input models.UserRegisterInput

	// 绑定请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "邮箱已存在",
		})
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "密码加密失败",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: input.Username,
		Password: hashedPassword,
		Email:    input.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "用户创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: "用户注册成功",
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input models.UserLoginInput

	// 绑定请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "令牌生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data: models.TokenResponse{
			Token: token,
		},
	})
} 