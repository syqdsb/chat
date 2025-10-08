package handlers

import (
	"chat-backend/models"
	"chat-backend/utils"
	"encoding/json"
	"net/http"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddFriendRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// 注册处理
func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.JSONError(w, "用户名、邮箱和密码不能为空", http.StatusBadRequest)
		return
	}

	// 创建用户
	user, err := models.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		utils.JSONError(w, "注册失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"message":  "注册成功",
		"token":    user.Token,
		"username": user.Username,
		"email":    user.Email,
	})
}

// 登录处理
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.Email == "" || req.Password == "" {
		utils.JSONError(w, "邮箱和密码不能为空", http.StatusBadRequest)
		return
	}

	// 用户登录
	user, err := models.LoginUser(req.Email, req.Password)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"message":  "登录成功",
		"token":    user.Token,
		"username": user.Username,
		"email":    user.Email,
	})
}

// 添加好友处理
func AddFriend(w http.ResponseWriter, r *http.Request) {
	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Token == "" {
		utils.JSONError(w, "用户名和 token 不能为空", http.StatusBadRequest)
		return
	}

	// 验证 token
	currentUser, err := models.GetUserByToken(req.Token)
	if err != nil {
		utils.JSONError(w, "无效的 token", http.StatusUnauthorized)
		return
	}

	// 获取要添加的好友
	friendUser, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.JSONError(w, "用户不存在", http.StatusNotFound)
		return
	}

	// 不能添加自己为好友
	if currentUser.ID == friendUser.ID {
		utils.JSONError(w, "不能添加自己为好友", http.StatusBadRequest)
		return
	}

	// 添加好友
	if err := models.AddFriend(currentUser.ID, friendUser.ID); err != nil {
		utils.JSONError(w, "添加好友失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"message":  "添加好友成功",
		"username": friendUser.Username,
	})
}