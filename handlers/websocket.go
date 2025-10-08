package handlers

import (
	"chat-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// 存储所有在线用户的连接
var (
	clients   = make(map[int]*websocket.Conn) // userID -> connection
	clientsMu sync.RWMutex
)

type WSMessage struct {
	Type     string `json:"type"`      // 消息类型: auth, send, get
	Token    string `json:"token"`     // 用户 token
	ToUser   string `json:"to_user"`   // 接收者用户名
	Content  string `json:"content"`   // 消息内容
	Messages []models.Message `json:"messages,omitempty"` // 消息列表
}

type WSResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WebSocket 聊天处理
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	defer conn.Close()

	var currentUser *models.User
	var userID int

	// 处理消息循环
	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}

		// 处理鉴权
		if msg.Type == "auth" {
			user, err := models.GetUserByToken(msg.Token)
			if err != nil {
				sendWSResponse(conn, false, "鉴权失败: "+err.Error(), nil)
				continue
			}

			currentUser = user
			userID = user.ID

			// 保存连接
			clientsMu.Lock()
			clients[userID] = conn
			clientsMu.Unlock()

			sendWSResponse(conn, true, "鉴权成功", map[string]interface{}{
				"user_id":  user.ID,
				"username": user.Username,
			})

			log.Printf("用户 %s (ID: %d) 已连接", user.Username, user.ID)
			continue
		}

		// 未鉴权不允许其他操作
		if currentUser == nil {
			sendWSResponse(conn, false, "请先进行鉴权", nil)
			continue
		}

		// 处理发送消息
		if msg.Type == "send" {
			if msg.ToUser == "" || msg.Content == "" {
				sendWSResponse(conn, false, "接收者用户名和消息内容不能为空", nil)
				continue
			}

			// 获取接收者信息
			toUser, err := models.GetUserByUsername(msg.ToUser)
			if err != nil {
				sendWSResponse(conn, false, "接收者不存在", nil)
				continue
			}

			// 保存消息到数据库
			if err := models.SaveMessage(currentUser.ID, toUser.ID, msg.Content); err != nil {
				sendWSResponse(conn, false, "消息保存失败", nil)
				continue
			}

			// 如果接收者在线，实时推送消息
			clientsMu.RLock()
			toConn, online := clients[toUser.ID]
			clientsMu.RUnlock()

			if online {
				notification := WSResponse{
					Success: true,
					Message: "新消息",
					Data: map[string]interface{}{
						"from_username": currentUser.Username,
						"from_user_id":  currentUser.ID,
						"content":       msg.Content,
					},
				}
				toConn.WriteJSON(notification)
			}

			sendWSResponse(conn, true, "消息发送成功", nil)
			log.Printf("用户 %s 向 %s 发送消息", currentUser.Username, toUser.Username)
			continue
		}

		// 处理获取未读消息
		if msg.Type == "get" {
			messages, err := models.GetUnreadMessages(currentUser.ID)
			if err != nil {
				sendWSResponse(conn, false, "获取消息失败", nil)
				continue
			}

			// 标记消息为已读
			if len(messages) > 0 {
				var msgIDs []int
				for _, m := range messages {
					msgIDs = append(msgIDs, m.ID)
				}
				models.MarkMessagesAsRead(msgIDs)
			}

			sendWSResponse(conn, true, "获取消息成功", map[string]interface{}{
				"count":    len(messages),
				"messages": messages,
			})
			continue
		}

		sendWSResponse(conn, false, "未知的消息类型", nil)
	}

	// 断开连接时清理
	if userID > 0 {
		clientsMu.Lock()
		delete(clients, userID)
		clientsMu.Unlock()
		log.Printf("用户 ID %d 已断开连接", userID)
	}
}

// 发送 WebSocket 响应
func sendWSResponse(conn *websocket.Conn, success bool, message string, data interface{}) {
	response := WSResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Println("发送响应失败:", err)
	}
}