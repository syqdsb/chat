package config

import (
	"log"
)

type Config struct {
	ServerPort string
	DBPath     string
	JWTSecret  string
}

var AppConfig *Config

func Init() {
	AppConfig = &Config{
		ServerPort: "8080",
		DBPath:     "./chat.db",
		JWTSecret:  "your-super-secret-key-change-in-production", // 生产环境请修改
	}
	log.Println("配置初始化完成")
}