package main

import (
	"homework_blog/config"
	"homework_blog/database"
	"homework_blog/routes"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err) //读取配置异常,直接退出程序
	}

	//初始化数据库
	db := database.InitDatabase()
	if db == nil {
		log.Fatal("Error initializing database")
	}
	//设置路由并启动HTTP服务
	r := routes.SetupRoutes(config, db)
	addr := config.Server.Host + ":" + config.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
