package handlers

import (
	"homework_blog/database"
	"homework_blog/models"
	"homework_blog/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userDb    *database.UserDb
	jwtSecret []byte
}

func NewUserHandler(userDb *database.UserDb, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userDb:    userDb,
		jwtSecret: jwtSecret}
}

// Register 用户注册
func (uh *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	user, err := uh.userDb.CreateUser(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(uh.jwtSecret, user.ID, user.Username)
	if err != nil {
		utils.Error(c, 500, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login 用户登录
func (uh *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 查找用户
	user, err := uh.userDb.CheckUser(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.Error(c, 401, "Invalid username or password")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(uh.jwtSecret, user.ID, user.Username)
	if err != nil {
		utils.Error(c, 500, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile 获取用户信息
func (uh *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, 401, "User not authenticated")
		return
	}

	user, err := uh.userDb.GetProfile(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, user)
}
