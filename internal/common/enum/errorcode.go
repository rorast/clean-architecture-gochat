// Package enum 提供系統中使用的枚舉常量
package enum

// ErrorCode 定義錯誤碼類型
type ErrorCode int

// 錯誤碼常量 - 分類定義
const (
	// 輸入驗證錯誤 (1xxx)
	ErrInvalidInput     ErrorCode = 1000
	ErrInvalidJSON      ErrorCode = 1001
	ErrValidationFailed ErrorCode = 1002
	ErrMissingField     ErrorCode = 1003
	ErrInvalidFormat    ErrorCode = 1004

	// 身份驗證錯誤 (2xxx)
	ErrUnauthorized       ErrorCode = 2000
	ErrInvalidCredentials ErrorCode = 2001
	ErrTokenExpired       ErrorCode = 2002
	ErrAccessDenied       ErrorCode = 2003

	// 用戶錯誤 (3xxx)
	ErrUserNotFound      ErrorCode = 3000
	ErrUserAlreadyExists ErrorCode = 3001
	ErrUserUpdateFailed  ErrorCode = 3002
	ErrUserDeleteFailed  ErrorCode = 3003

	// 聊天錯誤 (4xxx)
	ErrMessageSendFailed ErrorCode = 4000
	ErrMessageInvalid    ErrorCode = 4001
	ErrReceiverNotFound  ErrorCode = 4002
	ErrHistoryLoadFailed ErrorCode = 4003

	// 群組錯誤 (5xxx)
	ErrGroupNotFound     ErrorCode = 5000
	ErrGroupCreateFailed ErrorCode = 5001
	ErrGroupJoinFailed   ErrorCode = 5002
	ErrNotGroupMember    ErrorCode = 5003

	// 檔案錯誤 (6xxx)
	ErrFileUploadFailed  ErrorCode = 6000
	ErrFileInvalidFormat ErrorCode = 6001
	ErrFileTooLarge      ErrorCode = 6002
	ErrFileNotFound      ErrorCode = 6003

	// WebSocket錯誤 (7xxx)
	ErrWSConnectFailed ErrorCode = 7000
	ErrWSDisconnected  ErrorCode = 7001
	ErrWSMessageFailed ErrorCode = 7002

	// Redis錯誤 (8xxx)
	ErrRedisConnectionFailed ErrorCode = 8000
	ErrRedisOperationFailed  ErrorCode = 8001
	ErrRedisCacheMiss        ErrorCode = 8002

	// 系統錯誤 (9xxx)
	ErrInternalServer     ErrorCode = 9000
	ErrDatabaseError      ErrorCode = 9001
	ErrTimeout            ErrorCode = 9002
	ErrServiceUnavailable ErrorCode = 9003
)

// ErrorCodeDetails 錯誤碼詳細信息
var ErrorCodeDetails = map[ErrorCode]struct {
	Key     string
	Message string
}{
	ErrInvalidInput:     {"INVALID_INPUT", "無效的輸入"},
	ErrInvalidJSON:      {"INVALID_JSON", "無效的JSON格式"},
	ErrValidationFailed: {"VALIDATION_FAILED", "數據驗證失敗"},
	ErrMissingField:     {"MISSING_FIELD", "缺少必要欄位"},
	ErrInvalidFormat:    {"INVALID_FORMAT", "格式錯誤"},

	ErrUnauthorized:       {"UNAUTHORIZED", "未授權訪問"},
	ErrInvalidCredentials: {"INVALID_CREDENTIALS", "無效的憑證"},
	ErrTokenExpired:       {"TOKEN_EXPIRED", "令牌已過期"},
	ErrAccessDenied:       {"ACCESS_DENIED", "拒絕訪問"},

	ErrUserNotFound:      {"USER_NOT_FOUND", "用戶不存在"},
	ErrUserAlreadyExists: {"USER_EXISTS", "用戶已存在"},
	ErrUserUpdateFailed:  {"USER_UPDATE_FAILED", "用戶更新失敗"},
	ErrUserDeleteFailed:  {"USER_DELETE_FAILED", "用戶刪除失敗"},

	ErrMessageSendFailed: {"MESSAGE_SEND_FAILED", "消息發送失敗"},
	ErrMessageInvalid:    {"MESSAGE_INVALID", "無效的消息"},
	ErrReceiverNotFound:  {"RECEIVER_NOT_FOUND", "接收者不存在"},
	ErrHistoryLoadFailed: {"HISTORY_LOAD_FAILED", "歷史記錄加載失敗"},

	ErrGroupNotFound:     {"GROUP_NOT_FOUND", "群組不存在"},
	ErrGroupCreateFailed: {"GROUP_CREATE_FAILED", "創建群組失敗"},
	ErrGroupJoinFailed:   {"GROUP_JOIN_FAILED", "加入群組失敗"},
	ErrNotGroupMember:    {"NOT_GROUP_MEMBER", "非群組成員"},

	ErrFileUploadFailed:  {"FILE_UPLOAD_FAILED", "文件上傳失敗"},
	ErrFileInvalidFormat: {"FILE_INVALID_FORMAT", "無效的文件格式"},
	ErrFileTooLarge:      {"FILE_TOO_LARGE", "文件過大"},
	ErrFileNotFound:      {"FILE_NOT_FOUND", "文件不存在"},

	ErrWSConnectFailed: {"WS_CONNECT_FAILED", "WebSocket連接失敗"},
	ErrWSDisconnected:  {"WS_DISCONNECTED", "WebSocket連接中斷"},
	ErrWSMessageFailed: {"WS_MESSAGE_FAILED", "WebSocket消息發送失敗"},

	ErrRedisConnectionFailed: {"REDIS_CONNECTION_FAILED", "Redis連接失敗"},
	ErrRedisOperationFailed:  {"REDIS_OPERATION_FAILED", "Redis操作失敗"},
	ErrRedisCacheMiss:        {"REDIS_CACHE_MISS", "Redis緩存未命中"},

	ErrInternalServer:     {"INTERNAL_SERVER_ERROR", "服務器內部錯誤"},
	ErrDatabaseError:      {"DATABASE_ERROR", "數據庫錯誤"},
	ErrTimeout:            {"TIMEOUT", "請求超時"},
	ErrServiceUnavailable: {"SERVICE_UNAVAILABLE", "服務暫時不可用"},
}

// GetErrorDetails 獲取錯誤碼詳情
func GetErrorDetails(code ErrorCode) (string, string) {
	if details, exists := ErrorCodeDetails[code]; exists {
		return details.Key, details.Message
	}
	return "UNKNOWN_ERROR", "未知錯誤"
}
