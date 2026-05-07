package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章表
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text;not null"` //text类型用于存储长度可变且可能很大的字符串数据
	UserID    uint           `json:"user_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` //gorm.DeletedAt指定软删除

	//关联关系
	User     User      `json:"user,omitempty" gorm:"foreignkey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignkey:PostID"` //omitempty 空值时,该字段会被忽略
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type ResponsePosts struct {
	Posts    []Post `json:"posts"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
