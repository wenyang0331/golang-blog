package database

import (
	"errors"
	"homework_blog/models"
	"homework_blog/utils"

	"gorm.io/gorm"
)

type UserDb struct {
	db *gorm.DB
}

func NewUserOperator(db *gorm.DB) *UserDb {
	return &UserDb{db: db}
}

func (userDb *UserDb) CreateUser(req models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := userDb.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, utils.NewAppError(409, "Username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.NewAppError(500, "Database error checking username")
	}

	// 检查邮箱是否已存在
	if err := userDb.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, utils.NewAppError(409, "Email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.NewAppError(500, "Database error checking email")
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 密码会在BeforeCreate钩子中自动加密
	}

	if err := userDb.db.Create(&user).Error; err != nil {
		return nil, utils.NewAppError(500, "Failed to create user")
	}

	// 生成JWT token
	//token, err := utils.GenerateToken(jwtSecret, user.ID, user.Username)
	//if err != nil {
	//	return nil, utils.NewAppError(500, "Failed to generate token")
	//}

	return &user, nil
}

// 检查用户是否存在
func (userDb *UserDb) CheckUser(req models.LoginRequest) (*models.User, error) {
	// 查找用户
	var user models.User
	if err := userDb.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(401, "Invalid username or password")
		}
		return nil, utils.NewAppError(500, "Database error")
	}
	return &user, nil
}

// GetProfile 获取用户信息
func (userDb *UserDb) GetProfile(userID any) (*models.User, error) {
	var user models.User
	if err := userDb.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, "Database error")
	}
	return &user, nil
}
