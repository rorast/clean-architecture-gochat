package controllers

import (
	"clean-architecture-gochat/internal/domain/entities"
	"clean-architecture-gochat/pkg/response"
	"clean-architecture-gochat/pkg/utils"
	"clean-architecture-gochat/usecases/user"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserController struct {
	UserService user.Service
}

func NewUserController(us user.Service) *UserController {
	return &UserController{UserService: us}
}

func (uc *UserController) GetUserList(c *gin.Context) {
	users, err := uc.UserService.GetUserList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "列出用戶列表!", "data": users})
}

func (uc *UserController) CreateUser(c *gin.Context) {
	// 一次讀取並保留內容用於 debug
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("讀取請求體失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "請求體讀取失敗"})
		return
	}

	// 印出原始請求資料
	log.Printf("註冊請求原始數據: %s", string(body))

	// 定義註冊請求結構體
	var registerReq struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// 重設 Body 後再解析 JSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &registerReq); err != nil {
		log.Printf("JSON 解析錯誤: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "JSON 格式錯誤", "error": err.Error()})
		return
	}

	// 印出註冊請求資料
	log.Printf("接收到的註冊資料: %+v, name 長度: %d, password 長度: %d",
		registerReq, len(registerReq.Name), len(registerReq.Password))

	if registerReq.Name == "" || registerReq.Password == "" {
		log.Printf("用戶名或密碼為空: name=%s, password=%s", registerReq.Name, registerReq.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用戶名或密碼不能為空",
		})
		return
	}

	// 創建用戶實體
	user := entities.User{
		Name:     registerReq.Name,
		Password: registerReq.Password,
	}

	now := time.Now()
	user.LoginTime = now
	user.LogoutTime = now
	user.HeartbeatTime = now

	salt := fmt.Sprintf(viper.GetString("key.salt"), rand.Int31())
	log.Println("salt", salt)
	// 加密密碼
	user.Password = utils.MakePassword(user.Password, salt)
	//  將解碼的鹽值保存到 user 結構中
	user.Salt = salt
	// 生成用戶身份標識
	user.Identity = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int31())

	// 為空電子郵件生成唯一的佔位符
	if user.Email == "" {
		user.Email = fmt.Sprintf("temp_%d@example.com", time.Now().UnixNano())
	}

	if err := uc.UserService.CreateUser(c.Request.Context(), &user); err != nil {
		log.Printf("創建用戶失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "創建用戶失敗",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "新增用户成功！",
		"data": gin.H{
			"id":       user.ID,
			"identity": user.Identity,
		},
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的用戶ID"})
		return
	}
	if err := uc.UserService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "刪除用戶成功"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}

	if err := uc.UserService.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "修改用户成功！", "data": user})
}

func (uc *UserController) FindUserByNameAndPwd(c *gin.Context) {
	// 一次讀取並保留內容用於 debug
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("讀取請求體失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "請求體讀取失敗"})
		return
	}

	// 印出原始請求資料
	log.Printf("登錄請求原始數據: %s", string(body))

	// 定義登錄請求結構體
	var loginData struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// 重設 Body 後再解析 JSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.Unmarshal(body, &loginData); err != nil {
		log.Printf("JSON 解析錯誤: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "JSON 解析錯誤", "error": err.Error()})
		return
	}

	// 記錄解析後的數據
	log.Printf("解析後的登錄數據: %+v, name 長度: %d, password 長度: %d",
		loginData, len(loginData.Name), len(loginData.Password))

	// 檢查必要字段
	if loginData.Name == "" || loginData.Password == "" {
		log.Printf("用戶名或密碼為空: name=%s, password=%s", loginData.Name, loginData.Password)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "用戶名或密碼不能為空"})
		return
	}

	user, err := uc.UserService.FindUserByNameAndPwd(c.Request.Context(), loginData.Name, loginData.Password)
	if err != nil {
		log.Printf("登錄失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"code": -1, "message": err.Error()})
		return
	}

	log.Printf("登錄成功: %+v", user)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "登錄成功", "data": user})
}

func (uc *UserController) SearchFriend(c *gin.Context) {
	//userIdStr := c.Query("userId")
	//userId, err := strconv.Atoi(userIdStr)
	//log.Printf("用戶ID: %s", userIdStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "無效的用戶ID"})
	//	return
	//}
	var SearchFriendRequest struct {
		UserID int `json:"userId"`
	}
	if err := c.ShouldBindJSON(&SearchFriendRequest); err != nil {
		log.Println("解析 SearchFriendRequest-JSON 失敗:", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": err.Error()})
		return
	}
	log.Printf("解析後資料: %+v\n", SearchFriendRequest)

	friends, err := uc.UserService.SearchFriend(c.Request.Context(), uint(SearchFriendRequest.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "搜尋好友成功", "rows": friends})
}

// 加入好友
func (uc *UserController) AddFriend(c *gin.Context) {
	//targetNameStr := c.Query("targetName")
	//targetName, _ := strconv.Atoi(targetNameStr)
	//log.Printf("targetName: %s", targetName)
	//userIdStr := c.Query("userId")
	//userId, _ := strconv.Atoi(userIdStr)
	//log.Printf("userId: %s", userId)
	var AddFriendRequest struct {
		UserID   int `json:"userId"`
		TargetID int `json:"targetId"`
	}

	// 一次讀取並保留內容用於 debug
	body, err := io.ReadAll(c.Request.Body)
	log.Println("📦 JSON 原始內容: ", string(body))
	if err != nil {
		log.Printf("讀取請求體失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "請求體讀取失敗"})
		return
	}

	// 重設 Body 後再解析 JSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.Unmarshal(body, &AddFriendRequest); err != nil {
		log.Printf("JSON 解析錯誤: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "JSON 解析錯誤", "error": err.Error()})
		return
	}

	// 記錄解析後的數據
	log.Printf("解析後的登錄數據: %+v", AddFriendRequest)
	//userId, err := strconv.Atoi(c.PostForm("userId"))

	//if err != nil {
	//	utils.RespFail(c.Writer, "無效的用戶ID")
	//	return
	//}
	//targetName := c.PostForm("targetName")
	//targetName := req.TargetName
	//if targetName == "" {
	//	utils.RespFail(c.Writer, "目標用戶名不能為空")
	//	return
	//}
	//
	//err := uc.UserService.AddFriend(c.Request.Context(), uint(req.UserID), targetName)
	//if err != nil {
	//	utils.RespFail(c.Writer, err.Error())
	//	return
	//}
	//utils.RespOK(c.Writer, nil, "添加好友成功")
	err = uc.UserService.AddFriend(c.Request.Context(), uint(AddFriendRequest.UserID), uint(AddFriendRequest.TargetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	utils.RespOK(c.Writer, nil, "添加好友成功")
}

func (uc *UserController) FindUserByID(c *gin.Context) {
	var req struct {
		UserID int `json:"userId"`
	}

	// 檢查請求體
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("解析請求體失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "無效的請求格式",
			"error":   err.Error(),
		})
		return
	}

	// 檢查用戶ID是否有效
	if req.UserID <= 0 {
		log.Printf("無效的用戶ID: %d", req.UserID)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "無效的用戶ID",
		})
		return
	}

	// 查找用戶
	user, err := uc.UserService.FindUserByID(c.Request.Context(), uint(req.UserID))
	if err != nil {
		log.Printf("查找用戶失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "查找用戶失敗",
			"error":   err.Error(),
		})
		return
	}

	if user == nil {
		log.Printf("用戶不存在: %d", req.UserID)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用戶不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    user,
	})
}

// UploadFile 處理檔案上傳
// @Summary 上傳檔案
// @Description 上傳檔案並返回檔案URL
// @Tags 檔案
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "檔案"
// @Param type formData string true "檔案類型(avatar/file)"
// @Param userId formData int false "用戶ID(當type=avatar時必填)"
// @Success 200 {object} response.Response
// @Router /attach/upload [post]
func (c *UserController) UploadFile(ctx *gin.Context) {
	log.Printf("開始處理檔案上傳請求")
	log.Printf("請求方法: %s", ctx.Request.Method)
	log.Printf("Content-Type: %s", ctx.GetHeader("Content-Type"))

	// 檢查請求方法
	if ctx.Request.Method != http.MethodPost {
		log.Printf("錯誤：不支援的請求方法 %s", ctx.Request.Method)
		ctx.JSON(http.StatusMethodNotAllowed, response.Error("只允許 POST 方法"))
		return
	}

	// 檢查 Content-Type
	contentType := ctx.GetHeader("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Printf("錯誤：無效的 Content-Type %s", contentType)
		ctx.JSON(http.StatusBadRequest, response.Error("Content-Type 必須是 multipart/form-data"))
		return
	}

	// 獲取上傳的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("錯誤：獲取文件失敗 - %v", err)
		ctx.JSON(http.StatusBadRequest, response.Error("獲取文件失敗："+err.Error()))
		return
	}
	log.Printf("成功獲取文件：%s, 大小：%d bytes", file.Filename, file.Size)

	// 獲取檔案類型
	fileType := ctx.PostForm("type")
	log.Printf("檔案類型：%s", fileType)
	if fileType == "" {
		log.Printf("錯誤：未指定檔案類型")
		ctx.JSON(http.StatusBadRequest, response.Error("檔案類型不能為空，請指定 type=avatar 或 type=file"))
		return
	}

	// 根據檔案類型進行不同的處理
	switch fileType {
	case "avatar":
		// 驗證用戶ID
		userId := ctx.PostForm("userId")
		log.Printf("用戶ID：%s", userId)
		if userId == "" {
			log.Printf("錯誤：未提供用戶ID")
			ctx.JSON(http.StatusBadRequest, response.Error("上傳頭像時必須提供 userId"))
			return
		}

		// 驗證文件類型
		contentType := file.Header.Get("Content-Type")
		log.Printf("文件 Content-Type：%s", contentType)
		if !isValidImageType(contentType) {
			log.Printf("錯誤：不支援的文件類型 %s", contentType)
			ctx.JSON(http.StatusBadRequest, response.Error("頭像只支持 JPG、PNG 格式的圖片"))
			return
		}

		// 驗證文件大小（限制為 2MB）
		if file.Size > 2*1024*1024 {
			log.Printf("錯誤：文件大小超限 %d bytes", file.Size)
			ctx.JSON(http.StatusBadRequest, response.Error("頭像大小不能超過 2MB"))
			return
		}

		// 生成文件名
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("avatar_%s_%d%s", userId, time.Now().Unix(), ext)
		log.Printf("生成文件名：%s", fileName)

		// 確保上傳目錄存在
		uploadDir := "web/asset/avatars"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("錯誤：創建目錄失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("創建上傳目錄失敗"))
			return
		}

		// 保存文件
		filePath := filepath.Join(uploadDir, fileName)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("錯誤：保存文件失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("保存文件失敗"))
			return
		}
		log.Printf("文件已保存到：%s", filePath)

		// 更新用戶頭像路徑
		avatarPath := fmt.Sprintf("/web/asset/avatars/%s", fileName)
		userIdInt, err := strconv.ParseUint(userId, 10, 32)
		if err != nil {
			log.Printf("錯誤：無效的用戶ID格式 - %v", err)
			ctx.JSON(http.StatusBadRequest, response.Error("無效的用戶ID格式"))
			return
		}
		if err := c.UserService.UpdateAvatar(ctx, uint(userIdInt), avatarPath); err != nil {
			log.Printf("錯誤：更新用戶頭像失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("更新用戶頭像失敗"))
			return
		}

		log.Printf("頭像上傳成功：%s", avatarPath)
		ctx.JSON(http.StatusOK, response.Success(gin.H{
			"file_url":  avatarPath,
			"file_name": file.Filename,
			"file_size": file.Size,
		}))

	case "group":
		// 驗證用戶ID（群創建者ID）
		userId := ctx.PostForm("userId")
		log.Printf("群創建者ID：%s", userId)
		if userId == "" {
			log.Printf("錯誤：未提供用戶ID")
			ctx.JSON(http.StatusBadRequest, response.Error("上傳群頭像時必須提供 userId"))
			return
		}

		// 驗證文件類型
		contentType := file.Header.Get("Content-Type")
		log.Printf("文件 Content-Type：%s", contentType)
		if !isValidImageType(contentType) {
			log.Printf("錯誤：不支援的文件類型 %s", contentType)
			ctx.JSON(http.StatusBadRequest, response.Error("群頭像只支持 JPG、PNG 格式的圖片"))
			return
		}

		// 驗證文件大小（限制為 2MB）
		if file.Size > 2*1024*1024 {
			log.Printf("錯誤：文件大小超限 %d bytes", file.Size)
			ctx.JSON(http.StatusBadRequest, response.Error("群頭像大小不能超過 2MB"))
			return
		}

		// 生成文件名
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("group_%s_%d%s", userId, time.Now().Unix(), ext)
		log.Printf("生成文件名：%s", fileName)

		// 確保上傳目錄存在
		uploadDir := "web/asset/groups"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("錯誤：創建目錄失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("創建上傳目錄失敗"))
			return
		}

		// 保存文件
		filePath := filepath.Join(uploadDir, fileName)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("錯誤：保存文件失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("保存文件失敗"))
			return
		}
		log.Printf("文件已保存到：%s", filePath)

		// 返回群頭像 URL
		groupAvatarPath := fmt.Sprintf("/web/asset/groups/%s", fileName)
		log.Printf("群頭像上傳成功：%s", groupAvatarPath)
		ctx.JSON(http.StatusOK, response.Success(gin.H{
			"file_url":  groupAvatarPath,
			"file_name": file.Filename,
			"file_size": file.Size,
		}))

	case "file":
		// 驗證文件大小（限制為 10MB）
		if file.Size > 10*1024*1024 {
			log.Printf("錯誤：文件大小超限 %d bytes", file.Size)
			ctx.JSON(http.StatusBadRequest, response.Error("文件大小不能超過 10MB"))
			return
		}

		// 生成文件名
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("file_%d%s", time.Now().UnixNano(), ext)
		log.Printf("生成文件名：%s", fileName)

		// 確保上傳目錄存在
		uploadDir := "web/asset/files"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("錯誤：創建目錄失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("創建上傳目錄失敗"))
			return
		}

		// 保存文件
		filePath := filepath.Join(uploadDir, fileName)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("錯誤：保存文件失敗 - %v", err)
			ctx.JSON(http.StatusInternalServerError, response.ServerError("保存文件失敗"))
			return
		}
		log.Printf("文件已保存到：%s", filePath)

		// 返回文件URL
		fileURL := fmt.Sprintf("/web/asset/files/%s", fileName)
		log.Printf("文件上傳成功：%s", fileURL)
		ctx.JSON(http.StatusOK, response.Success(gin.H{
			"file_url":  fileURL,
			"file_name": file.Filename,
			"file_size": file.Size,
		}))

	default:
		log.Printf("錯誤：不支援的檔案類型 %s", fileType)
		ctx.JSON(http.StatusBadRequest, response.Error("不支持的檔案類型，請使用 avatar、group 或 file"))
	}
}

// isValidImageType 驗證文件類型是否為有效的圖片
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	return validTypes[contentType]
}
