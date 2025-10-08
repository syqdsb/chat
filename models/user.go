package models

import (
	"chat-backend/database"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Token    string `json:"token"`
}

type Message struct {
	ID         int    `json:"id"`
	FromUserID int    `json:"from_user_id"`
	ToUserID   int    `json:"to_user_id"`
	Content    string `json:"content"`
	IsRead     int    `json:"is_read"`
	CreatedAt  string `json:"created_at"`
}

// 创建用户
func CreateUser(username, email, password string) (*User, error) {
	hashedPassword := hashPassword(password)
	token := generateToken(email)

	result, err := database.DB.Exec(
		"INSERT INTO users (username, email, password, token) VALUES (?, ?, ?, ?)",
		username, email, hashedPassword, token,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &User{
		ID:       int(id),
		Username: username,
		Email:    email,
		Token:    token,
	}, nil
}

// 用户登录
func LoginUser(email, password string) (*User, error) {
	hashedPassword := hashPassword(password)
	var user User

	err := database.DB.QueryRow(
		"SELECT id, username, email, token FROM users WHERE email = ? AND password = ?",
		email, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Token)

	if err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	return &user, nil
}

// 通过 token 获取用户
func GetUserByToken(token string) (*User, error) {
	var user User
	err := database.DB.QueryRow(
		"SELECT id, username, email, token FROM users WHERE token = ?",
		token,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Token)

	if err != nil {
		return nil, errors.New("无效的 token")
	}
	return &user, nil
}

// 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := database.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// 添加好友
func AddFriend(userID, friendID int) error {
	// 检查是否已经是好友
	var count int
	err := database.DB.QueryRow(
		"SELECT COUNT(*) FROM friends WHERE user_id = ? AND friend_id = ?",
		userID, friendID,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("已经是好友了")
	}

	// 添加双向好友关系
	_, err = database.DB.Exec(
		"INSERT INTO friends (user_id, friend_id) VALUES (?, ?)",
		userID, friendID,
	)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		"INSERT INTO friends (user_id, friend_id) VALUES (?, ?)",
		friendID, userID,
	)
	return err
}

// 保存消息
func SaveMessage(fromUserID, toUserID int, content string) error {
	_, err := database.DB.Exec(
		"INSERT INTO messages (from_user_id, to_user_id, content) VALUES (?, ?, ?)",
		fromUserID, toUserID, content,
	)
	return err
}

// 获取未读消息
func GetUnreadMessages(userID int) ([]Message, error) {
	rows, err := database.DB.Query(
		"SELECT id, from_user_id, to_user_id, content, is_read, created_at FROM messages WHERE to_user_id = ? AND is_read = 0 ORDER BY created_at ASC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.FromUserID, &msg.ToUserID, &msg.Content, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// 标记消息为已读
func MarkMessagesAsRead(messageIDs []int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	for _, msgID := range messageIDs {
		_, err := database.DB.Exec("UPDATE messages SET is_read = 1 WHERE id = ?", msgID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 密码哈希
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// 生成 token
func generateToken(email string) string {
	data := email + "-token-" + string(rune(len(email)))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}