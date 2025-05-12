package controllers

import (
	"clean-architecture-gochat/internal/domain/entities"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (ic *IndexController) GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("web/views/index.html", "web/views/chat/head.html")
	if err != nil {
		log.Fatal(err)
	}
	ind.Execute(c.Writer, "index")
}

func (ic *IndexController) ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("web/views/user/register.html", "web/views/chat/head.html")
	if err != nil {
		log.Fatal(err)
	}
	ind.Execute(c.Writer, "register")
}

func (ic *IndexController) ToChat(c *gin.Context) {
	// 獲取並驗證參數
	userIdStr := c.Query("userId")
	token := c.Query("token")

	// 記錄接收到的參數
	log.Printf("ToChat 接收到的參數: userId=%s, token=%s", userIdStr, token)

	// 驗證參數
	if userIdStr == "" || token == "" {
		log.Printf("參數缺失: userId=%s, token=%s", userIdStr, token)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 轉換 userId
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Printf("userId 轉換失敗: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 加載模板
	ind, err := template.ParseFiles(
		"web/views/chat/index.html",
		"web/views/chat/head.html",
		"web/views/chat/foot.html",
		"web/views/chat/tabmenu.html",
		"web/views/chat/concat.html",
		"web/views/chat/group.html",
		"web/views/chat/profile.html",
		"web/views/chat/createcom.html",
		"web/views/chat/userinfo.html",
		"web/views/chat/main.html")
	if err != nil {
		log.Printf("加載模板失敗: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 創建用戶數據
	user := entities.User{
		ID:       uint(userId),
		Identity: token,
	}

	// 記錄即將渲染的數據
	log.Printf("準備渲染聊天頁面，用戶數據: %+v", user)

	// 渲染模板
	if err := ind.Execute(c.Writer, user); err != nil {
		log.Printf("渲染模板失敗: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}
}

//func (ic *IndexController) SearchFriends(c *gin.Context) {
//	// 獲取用戶ID
//	userIdStr := c.Query("userId")
//	if userIdStr == "" {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    -1,
//			"message": "用戶ID不能為空",
//		})
//		return
//	}
//
//	// 轉換用戶ID為整數
//	userId, err := strconv.Atoi(userIdStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    -1,
//			"message": "無效的用戶ID",
//		})
//		return
//	}
//
//	// 從數據庫中獲取好友列表
//	// TODO: 這裡需要實現獲取好友列表的具體邏輯
//	// 暫時返回一個示例響應
//	c.JSON(http.StatusOK, gin.H{
//		"code":    0,
//		"message": "獲取好友列表成功",
//		"rows": []gin.H{
//			{
//				"id":     1,
//				"name":   "好友1",
//				"avatar": "/web/asset/images/avatar.png",
//			},
//			{
//				"id":     2,
//				"name":   "好友2",
//				"avatar": "/web/asset/images/avatar.png",
//			},
//		},
//	})
//}
