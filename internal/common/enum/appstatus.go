package enum

// AppStatus 定義應用狀態碼類型
type AppStatus int

// 應用狀態碼常量
const (
	// 成功狀態
	StatusSuccess AppStatus = 0

	// 一般失敗狀態
	StatusFailed AppStatus = 1

	// 登錄相關狀態
	StatusNotLoggedIn     AppStatus = 401 // 未登錄
	StatusLoginFailed     AppStatus = 402 // 登錄失敗
	StatusTokenExpired    AppStatus = 403 // Token過期
	StatusInvalidToken    AppStatus = 407 // 無效Token
	StatusAccountDisabled AppStatus = 409 // 賬號已禁用
	StatusAccessDenied    AppStatus = 406 // 拒絕訪問

	// 請求參數相關狀態
	StatusInvalidParam     AppStatus = 400 // 無效參數
	StatusResourceNotFound AppStatus = 404 // 資源不存在
	StatusMethodNotAllowed AppStatus = 405 // 方法不允許
	StatusRequestTimeout   AppStatus = 408 // 請求超時
	StatusConflict         AppStatus = 410 // 資源衝突

	// 服務器相關狀態
	StatusServerError        AppStatus = 500 // 服務器錯誤
	StatusServiceUnavailable AppStatus = 503 // 服務不可用
	StatusDatabaseError      AppStatus = 510 // 數據庫錯誤
	StatusRedisError         AppStatus = 511 // Redis錯誤
)

// StatusMessage 狀態碼對應的訊息
var StatusMessage = map[AppStatus]string{
	StatusSuccess:            "操作成功",
	StatusFailed:             "操作失敗",
	StatusNotLoggedIn:        "未登錄",
	StatusLoginFailed:        "登錄失敗",
	StatusTokenExpired:       "令牌已過期",
	StatusInvalidToken:       "無效的令牌",
	StatusAccountDisabled:    "賬號已禁用",
	StatusAccessDenied:       "拒絕訪問",
	StatusInvalidParam:       "無效的參數",
	StatusResourceNotFound:   "資源不存在",
	StatusMethodNotAllowed:   "方法不允許",
	StatusRequestTimeout:     "請求超時",
	StatusConflict:           "資源衝突",
	StatusServerError:        "服務器錯誤",
	StatusServiceUnavailable: "服務暫時不可用",
	StatusDatabaseError:      "數據庫錯誤",
	StatusRedisError:         "Redis服務錯誤",
}

// GetStatusMessage 獲取狀態碼對應的訊息
func GetStatusMessage(status AppStatus) string {
	if msg, exists := StatusMessage[status]; exists {
		return msg
	}
	return "未知狀態"
}
