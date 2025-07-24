package models

import (
	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(200);not null" json:"title"`
	Content  string    `gorm:"type:text;not null" json:"content"`
	UserID   uint      `json:"user_id"`
	User     User      `json:"user,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// PostInput 文章输入
type PostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
