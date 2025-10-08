package database

import (
	"chat-backend/config"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", config.AppConfig.DBPath)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		token TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建好友关系表
	createFriendsTable := `
	CREATE TABLE IF NOT EXISTS friends (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		friend_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, friend_id)
	);`

	// 创建消息表
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		is_read INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatal("创建用户表失败:", err)
	}

	if _, err := DB.Exec(createFriendsTable); err != nil {
		log.Fatal("创建好友表失败:", err)
	}

	if _, err := DB.Exec(createMessagesTable); err != nil {
		log.Fatal("创建消息表失败:", err)
	}

	log.Println("数据库初始化完成")
}