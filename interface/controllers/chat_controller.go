package controllers

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/internal/domain/repositories"
	"clean-architecture-gochat/internal/usecases/chat"
	ws "clean-architecture-gochat/internal/usecases/websocket"

	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

type ChatController struct {
	privateChatService chat.Service
	groupChatService   chat.GroupChatService
	connectionService  ws.ConnectionService
	messageService     repositories.MessageService
}

func NewChatController(
	privateChatService chat.Service,
	groupChatService chat.GroupChatService,
	connectionService ws.ConnectionService,
	messageService repositories.MessageService,
) *ChatController {
	return &ChatController{
		privateChatService: privateChatService,
		groupChatService:   groupChatService,
		connectionService:  connectionService,
		messageService:     messageService,
	}
}

// WebSocket 處理
func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	// 獲取用戶ID
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	// 升級 HTTP 連接為 WebSocket
	upgrader := gorilla.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// 使用 connectionService 處理 WebSocket 連接
	if err := cc.connectionService.Connect(c.Request.Context(), conn, uint(userID)); err != nil {
		log.Printf("Failed to handle WebSocket connection: %v", err)
		conn.Close()
		return
	}

	// 連接成功，返回 200 狀態碼
	c.Status(http.StatusOK)
}

// 發送私人訊息
func (cc *ChatController) SendPrivateMessage(c *gin.Context) {
	var req struct {
		FromUserID uint   `json:"fromUserId"`
		ToUserID   uint   `json:"toUserId"`
		Content    string `json:"content"`
		Type       int    `json:"type"`
		Media      int    `json:"media"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的請求格式"})
		return
	}

	message := &entities.Message{
		UserId:    req.FromUserID,
		TargetId:  req.ToUserID,
		Content:   req.Content,
		Type:      entities.MessageType(req.Type),
		Media:     entities.MediaType(req.Media),
		CreatedAt: time.Now(),
	}

	if err := cc.messageService.SendPrivateMessage(c.Request.Context(), message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "訊息發送成功"})
}

// 獲取私人聊天歷史
func (cc *ChatController) GetPrivateHistory(c *gin.Context) {
	fromUserID, err := strconv.ParseUint(c.Query("fromUserId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的發送者ID"})
		return
	}

	toUserID, err := strconv.ParseUint(c.Query("toUserId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的接收者ID"})
		return
	}

	messages, err := cc.messageService.GetPrivateHistory(c.Request.Context(), uint(fromUserID), uint(toUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "獲取聊天歷史成功", "data": messages})
}

// 創建群組
func (cc *ChatController) CreateGroup(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     int    `json:"type" binding:"required"`
		Desc     string `json:"desc" binding:"required"`
		Size     int    `json:"size" binding:"required"`
		JoinType int    `json:"join_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}

	// 獲取用戶ID
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "userId is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid userId"})
		return
	}

	group, err := cc.groupChatService.CreateGroup(
		c.Request.Context(),
		req.Name,
		uint(userID),
		req.Type,
		req.Desc,
		req.Size,
		req.JoinType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// 更新群組
func (cc *ChatController) UpdateGroup(c *gin.Context) {
	var group entities.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}

	if err := cc.groupChatService.UpdateGroup(c.Request.Context(), &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// 刪除群組
func (cc *ChatController) DeleteGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "group id is required"})
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid group id"})
		return
	}

	if err := cc.groupChatService.DeleteGroup(c.Request.Context(), uint(groupID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// 獲取群組信息
func (cc *ChatController) GetGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "group id is required"})
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid group id"})
		return
	}

	group, err := cc.groupChatService.GetGroup(c.Request.Context(), uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// 獲取群組列表
func (cc *ChatController) GetGroups(c *gin.Context) {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "userId is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid userId"})
		return
	}

	groups, err := cc.groupChatService.GetUserGroups(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": groups})
}

// 獲取群組成員
func (cc *ChatController) GetGroupMembers(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "group id is required"})
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid group id"})
		return
	}

	members, err := cc.groupChatService.GetGroupMembers(c.Request.Context(), uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": members})
}

// 添加群組成員
func (cc *ChatController) AddGroupMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "group id is required"})
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid group id"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}

	if err := cc.groupChatService.AddMember(c.Request.Context(), uint(groupID), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// 移除群組成員
func (cc *ChatController) RemoveGroupMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "group id is required"})
		return
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "invalid group id"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}

	if err := cc.groupChatService.RemoveMember(c.Request.Context(), uint(groupID), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// 發送群組消息
func (cc *ChatController) SendGroupMessage(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"userId"`
		RoomID  uint   `json:"roomId"`
		Content string `json:"content"`
		Type    int    `json:"type"`
		Media   int    `json:"media"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的請求格式"})
		return
	}

	message := &entities.Message{
		UserId:    req.UserID,
		RoomID:    req.RoomID,
		Content:   req.Content,
		Type:      entities.MessageType(req.Type),
		Media:     entities.MediaType(req.Media),
		CreatedAt: time.Now(),
	}

	if err := cc.messageService.SendGroupMessage(c.Request.Context(), message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "群組訊息發送成功"})
}

// 獲取群組聊天歷史
func (cc *ChatController) GetGroupHistory(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Query("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的群組ID"})
		return
	}

	messages, err := cc.messageService.GetGroupHistory(c.Request.Context(), uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "獲取群組聊天歷史成功", "data": messages})
}

// CreateCustomGroup 處理前端自定義格式的群組創建請求
func (cc *ChatController) CreateCustomGroup(c *gin.Context) {
	log.Println("接收到創建群組請求")

	// 自定義群組創建請求
	var req struct {
		OwnerID  uint   `json:"ownerId"`
		Icon     string `json:"icon"`
		Cate     string `json:"cate"`
		Name     string `json:"name"`
		Memo     string `json:"memo"`
		Size     string `json:"size"`
		JoinType string `json:"joinType"`
	}

	// 讀取請求內容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("讀取請求體失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "請求體讀取失敗"})
		return
	}

	// 打印原始請求數據
	log.Printf("創建群組請求原始數據: %s", string(body))

	// 重置請求體並解析 JSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("解析 JSON 失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "JSON 解析錯誤", "error": err.Error()})
		return
	}

	// 驗證必要字段
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "群名稱不能為空"})
		return
	}

	if req.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "群主 ID 不能為空"})
		return
	}

	// 轉換數據類型
	cate, err := strconv.Atoi(req.Cate)
	if err != nil {
		cate = 0 // 默認
	}

	size, err := strconv.Atoi(req.Size)
	if err != nil {
		size = 50 // 默認
	}

	joinType, err := strconv.Atoi(req.JoinType)
	if err != nil {
		joinType = 0 // 默認
	}

	// 創建群組
	log.Printf("創建群組參數: name=%s, ownerID=%d, cate=%d, memo=%s, size=%d, joinType=%d, icon=%s",
		req.Name, req.OwnerID, cate, req.Memo, size, joinType, req.Icon)

	group, err := cc.groupChatService.CreateGroup(
		c.Request.Context(),
		req.Name,
		req.OwnerID,
		cate,
		req.Memo,
		size,
		joinType,
	)
	if err != nil {
		log.Printf("創建群組失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	// 如果有圖標URL，則更新群組圖標
	if req.Icon != "" {
		group.Icon = req.Icon
		if err := cc.groupChatService.UpdateGroup(c.Request.Context(), group); err != nil {
			log.Printf("更新群組圖標失敗: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "創建群組成功", "data": group})
}

// JoinGroup 處理加入群組的請求
func (cc *ChatController) JoinGroup(c *gin.Context) {
	var req struct {
		ComID  string `json:"comId"`
		UserID uint   `json:"userId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "請求參數錯誤"})
		return
	}

	// 將 comId 轉換為 uint
	comID, err := strconv.ParseUint(req.ComID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的群組ID"})
		return
	}

	// 檢查用戶是否已經是群組成員
	isMember, err := cc.groupChatService.IsGroupMember(c.Request.Context(), uint(comID), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": "檢查群組成員狀態失敗"})
		return
	}

	if isMember {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "您已經是該群組的成員"})
		return
	}

	// 添加用戶到群組
	if err := cc.groupChatService.AddMember(c.Request.Context(), uint(comID), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": "加入群組失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "加入群組成功"})
}

// GetRecentMessages 獲取最近的訊息列表
func (cc *ChatController) GetRecentMessages(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的用戶ID"})
		return
	}

	messages, err := cc.messageService.GetRecentMessages(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "獲取最近訊息成功", "data": messages})
}
