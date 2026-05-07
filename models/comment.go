package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论表
type Comment struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Content   string         `json:"content" gorm:"type:text"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	//关联关系
	User User `json:"user,omitempty" gorm:"foreignkey:UserID"` //作用:关联后可以直接comment.User.Name,避免二次查找
	Post Post `json:"post,omitempty" gorm:"foreignkey:PostID"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}
