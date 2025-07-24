package models

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user,omitempty"`
	PostID  uint   `json:"post_id"`
	Post    Post   `json:"post,omitempty" gorm:"foreignKey:PostID"`
}

// CommentInput 评论输入
type CommentInput struct {
	Content string `json:"content" binding:"required"`
}
