package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhy/blog-api/config"
	"github.com/xhy/blog-api/models"
	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var input models.CommentInput
	postID := c.Param("id")
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

	// 查询文章是否存在
	if err := config.DB.First(&post, postID).Error; err != nil {
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

	// 绑定请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: input.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "评论创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: "评论创建成功",
		Data:    comment,
	})
}

// GetComments 获取文章的所有评论
func GetComments(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post
	var comments []models.Comment

	// 查询文章是否存在
	if err := config.DB.First(&post, postID).Error; err != nil {
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

	// 查询评论列表
	if err := config.DB.Where("post_id = ?", postID).Preload("User").Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取评论列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "获取评论列表成功",
		Data:    comments,
	})
} 