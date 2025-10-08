package main

import (
	"chat-backend/config"
	"chat-backend/database"
	"chat-backend/routes"
	"log"
	"net/http"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	// 设置路由
	router := routes.SetupRouter()

	// 启动服务器
	port := config.AppConfig.ServerPort
	log.Printf("服务器启动在端口: %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}