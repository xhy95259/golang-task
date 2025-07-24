package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xhy/blog-api/config"
	"github.com/xhy/blog-api/models"
	"gorm.io/gorm"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var input models.PostInput

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建文章
	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "文章创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: "文章创建成功",
		Data:    post,
	})
}

// GetPosts 获取所有文章
func GetPosts(c *gin.Context) {
	var posts []models.Post

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize

	// 查询文章列表
	if err := config.DB.Preload("User").Order("created_at desc").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文章列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "获取文章列表成功",
		Data:    posts,
	})
}

// GetPost 获取单个文章
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// 查询文章
	if err := config.DB.Preload("User").Preload("Comments").Preload("Comments.User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    http.StatusNotFound,
				Message: "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文章失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "获取文章成功",
		Data:    post,
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var input models.PostInput
	var post models.Post

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 查询文章
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    http.StatusNotFound,
				Message: "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文章失败: " + err.Error(),
		})
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    http.StatusForbidden,
			Message: "没有权限更新此文章",
		})
		return
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 更新文章
	post.Title = input.Title
	post.Content = input.Content

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "更新文章失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "更新文章成功",
		Data:    post,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	// 查询文章
	if err := config.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    http.StatusNotFound,
				Message: "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取文章失败: " + err.Error(),
		})
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    http.StatusForbidden,
			Message: "没有权限删除此文章",
		})
		return
	}

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "删除文章失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "删除文章成功",
	})
} 