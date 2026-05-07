package database

import (
	"homework_blog/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() *gorm.DB {
	var err error

	// 连接MySQL数据库
	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect sqlite database: %v", err)
		return nil
	}

	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return nil
	}

	log.Println("Sqlite database connected and migrated successfully")
	return db
}

// GetDB 获取数据库连接实例
//func GetDB() *gorm.DB {
//	return DB
//}
