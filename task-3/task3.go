package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:50;not null;unique"`
	Email     string `gorm:"size:100;not null;unique"`
	Password  string `gorm:"size:100;not null"`
	PostCount int    `gorm:"default:0"` // 用户发布的文章数量
	Posts     []Post `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Post 文章模型
type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"size:200;not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        uint      `gorm:"index;not null"`
	User          User      `gorm:"foreignKey:UserID"`
	CommentCount  int       `gorm:"default:0"` // 文章的评论数量
	CommentStatus int       `gorm:"default:0"` // 评论状态：0-未评论，1-已评论
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// Comment 评论模型
type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	PostID    uint   `gorm:"index;not null"`
	Post      Post   `gorm:"foreignKey:PostID"`
	UserID    uint   `gorm:"index;not null"` // 评论者ID
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// AfterCreate 在创建文章后更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章数量
	tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1))
	return
}

// AfterCreate 在创建评论后更新文章的评论数量和状态
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	// 更新文章的评论数量和状态
	tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(map[string]interface{}{
		"comment_count":  gorm.Expr("comment_count + ?", 1),
		"comment_status": 1, // 已评论
	})
	return
}

// AfterDelete 在删除评论后检查文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var count int64
	// 先更新评论数量
	tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	// 检查文章是否还有其他评论
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)
	// 如果没有评论，更新评论状态为未评论
	if count == 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", 0)
	}
	return
}

func main() {
	// 连接MySQL数据库
	dsn := "root:root@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	fmt.Println("成功连接到MySQL数据库")

	//自动迁移模型到数据库表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}
	fmt.Println("数据库表创建成功")

	// 创建测试数据
	createTestData(db)

	// 查询某个用户发布的所有文章及其对应的评论信息
	queryUserPostsWithComments(db, 15)

	// 查询评论数量最多的文章信息
	queryMostCommentedPost(db)
}

// 创建测试数据
func createTestData(db *gorm.DB) {
	// 清空表数据
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	// 创建用户
	users := []User{
		{Username: "alice", Email: "alice@example.com", Password: "password123"},
		{Username: "bob", Email: "bob@example.com", Password: "password456"},
	}
	for i := range users {
		db.Create(&users[i])
	}

	// 创建文章
	posts := []Post{
		{Title: "Go语言入门", Content: "Go是一种简单高效的编程语言...", UserID: users[0].ID},
		{Title: "Gorm使用教程", Content: "Gorm是Go语言中流行的ORM库...", UserID: users[0].ID},
		{Title: "Web开发实践", Content: "使用Go进行Web开发的最佳实践...", UserID: users[1].ID},
	}
	for i := range posts {
		db.Create(&posts[i])
	}

	// 创建评论
	comments := []Comment{
		{Content: "非常有用的教程！", PostID: posts[0].ID, UserID: users[1].ID},
		{Content: "我学到了很多，谢谢分享。", PostID: posts[0].ID, UserID: users[1].ID},
		{Content: "这篇文章解决了我的问题。", PostID: posts[1].ID, UserID: users[1].ID},
		{Content: "期待更多的教程。", PostID: posts[2].ID, UserID: users[0].ID},
	}
	for i := range comments {
		db.Create(&comments[i])
	}

	fmt.Println("测试数据创建成功")
}

// 查询某个用户发布的所有文章及其对应的评论信息
func queryUserPostsWithComments(db *gorm.DB, userID uint) {
	var user User

	// 使用Preload预加载关联数据
	result := db.Preload("Posts.Comments.User").First(&user, userID)
	if result.Error != nil {
		fmt.Println("查询用户失败:", result.Error)
		return
	}

	fmt.Printf("\n用户 %s 的文章及评论信息:\n", user.Username)
	for _, post := range user.Posts {
		fmt.Printf("文章: %s (评论数: %d, 评论状态: %d)\n", post.Title, len(post.Comments), post.CommentStatus)
		for i, comment := range post.Comments {
			fmt.Printf("  评论 %d: %s (by %s)\n", i+1, comment.Content, comment.User.Username)
		}
		fmt.Println()
	}
}

// 查询评论数量最多的文章信息
func queryMostCommentedPost(db *gorm.DB) {
	var post Post

	// 使用Order排序并限制结果数量
	result := db.Preload("Comments").Preload("User").Order("comment_count DESC").First(&post)
	if result.Error != nil {
		fmt.Println("查询文章失败:", result.Error)
		return
	}

	fmt.Println("\n评论数量最多的文章:")
	fmt.Printf("标题: %s\n", post.Title)
	fmt.Printf("作者: %s\n", post.User.Username)
	fmt.Printf("评论数: %d\n", post.CommentCount)
	fmt.Printf("评论状态: %d (0-未评论, 1-已评论)\n", post.CommentStatus)
	fmt.Println("评论列表:")
	for i, comment := range post.Comments {
		fmt.Printf("  评论 %d: %s\n", i+1, comment.Content)
	}
}
