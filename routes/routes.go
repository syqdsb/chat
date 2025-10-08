package routes

import (
	"chat-backend/handlers"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// API 路由
	mux.HandleFunc("/api/signup", handlers.Signup)
	mux.HandleFunc("/api/login", handlers.Login)
	mux.HandleFunc("/api/adduser", handlers.AddFriend)

	// WebSocket 路由
	mux.HandleFunc("/ws/chat", handlers.HandleWebSocket)

	return mux
}