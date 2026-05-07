package handlers

import (
	"homework_blog/database"
	"homework_blog/models"
	"homework_blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
 * 文章处理器 处理文章相关的请求
 */
type PostHandler struct {
	pstDb     *database.PostDb
	jwtSecret []byte
}

func NewPostHandler(pstDb *database.PostDb, jwtSecret []byte) *PostHandler {
	return &PostHandler{
		pstDb:     pstDb,
		jwtSecret: jwtSecret}
}

// CreatePost 创建文章
func (ph *PostHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	post, err := ph.pstDb.CreatePost(req, userID.(uint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, post)
}

// GetPosts 获取文章列表
func (ph *PostHandler) GetPosts(c *gin.Context) {
	//分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	result, err := ph.pstDb.GetPosts(page, pageSize, offset)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, result)
}

func (ph *PostHandler) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}
	post, err := ph.pstDb.GetPost(uint(postID))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, post)
}

// UpdatePost 更新文章
func (ph *PostHandler) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var post models.Post
	if err := ph.pstDb.Db.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only update your own posts")
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content

	if err := ph.pstDb.Db.Save(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to update post")
		return
	}

	// 预加载用户信息
	ph.pstDb.Db.Preload("User").First(&post, post.ID)

	utils.Success(c, post)
}

// DeletePost 删除文章
func (ph *PostHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var post models.Post
	if err := ph.pstDb.Db.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only delete your own posts")
		return
	}

	// 删除文章（软删除）
	if err := ph.pstDb.Db.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete post")
		return
	}

	utils.Success(c, gin.H{"message": "Post deleted successfully"})
}
