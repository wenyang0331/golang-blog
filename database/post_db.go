package database

import (
	"homework_blog/models"
	"homework_blog/utils"
	"net/http"

	"gorm.io/gorm"
)

type PostDb struct {
	Db *gorm.DB
}

func NewPostOperator(db *gorm.DB) *PostDb {
	return &PostDb{Db: db}
}

func (postDb *PostDb) CreatePost(req models.CreatePostRequest, userID any) (*models.Post, error) {
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := postDb.Db.Create(&post).Error; err != nil {
		return nil, utils.NewAppError(409, "Failed to create post")
	}

	// 预加载用户信息
	postDb.Db.Preload("User").First(&post, post.ID)

	return &post, nil
}

func (postDb *PostDb) GetPosts(page int, pageSize int, offset int) (models.ResponsePosts, error) {
	var posts []models.Post
	var responsePosts models.ResponsePosts
	// 查询文章列表，预加载用户信息
	if err := postDb.Db.Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return responsePosts, utils.NewAppError(409, "Failed to get posts")
	}

	// 获取总数
	var total int64
	postDb.Db.Model(&models.Post{}).Count(&total)

	responsePosts.Posts = posts
	responsePosts.Total = total
	responsePosts.Page = page
	responsePosts.PageSize = pageSize
	return responsePosts, nil
}

func (postDb *PostDb) GetPost(postID uint) (models.Post, error) {
	var post models.Post
	if err := postDb.Db.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		return post, utils.NewAppError(http.StatusNotFound, "Post not found")
	}
	return post, nil
}
