package errorcode

// 通用錯誤 (1000-1999)
const (
	// API綁定錯誤
	ErrInvalidJSONSyntax  = 1000 // JSON 格式錯誤
	ErrInvalidJSONType    = 1001 // JSON 欄位型別錯誤
	ErrInvalidQueryParams = 1002 // 查詢參數錯誤
	ErrInvalidFormData    = 1003 // 表單數據錯誤
	ErrInvalidURIParams   = 1004 // URI 參數錯誤

	// 輸入驗證錯誤
	ErrValidationFailed = 1100 // 欄位驗證失敗
	ErrMissingField     = 1101 // 缺少必要欄位
)

// 使用者錯誤 (2000-2999)
const (
	ErrUserNotFound       = 2000 // 使用者不存在
	ErrUserAlreadyExists  = 2001 // 使用者已存在
	ErrInvalidCredentials = 2002 // 無效的憑證
	ErrUserUpdateFailed   = 2003 // 使用者更新失敗
	ErrUserDeleteFailed   = 2004 // 使用者刪除失敗
	ErrUnauthorized       = 2005 // 未授權訪問
)

// 聊天功能錯誤 (3000-3999)
const (
	ErrMessageSendFailed    = 3000 // 消息發送失敗
	ErrReceiverNotFound     = 3001 // 接收者不存在
	ErrInvalidMessageFormat = 3002 // 無效的消息格式
	ErrMessageNotFound      = 3003 // 消息不存在
	ErrHistoryLoadFailed    = 3004 // 歷史記錄載入失敗
)

// 群組聊天錯誤 (4000-4999)
const (
	ErrGroupCreateFailed           = 4000 // 群組創建失敗
	ErrGroupNotFound               = 4001 // 群組不存在
	ErrGroupJoinFailed             = 4002 // 加入群組失敗
	ErrNotGroupMember              = 4003 // 非群組成員
	ErrInsufficientGroupPermission = 4004 // 群組權限不足
)

// 檔案/媒體錯誤 (5000-5999)
const (
	ErrFileUploadFailed      = 5000 // 檔案上傳失敗
	ErrInvalidFileFormat     = 5001 // 無效的檔案格式
	ErrFileTooLarge          = 5002 // 檔案過大
	ErrFileNotFound          = 5003 // 檔案不存在
	ErrImageProcessingFailed = 5004 // 圖片處理失敗
)

// WebSocket錯誤 (6000-6999)
const (
	ErrWSConnectionFailed  = 6000 // WebSocket連接失敗
	ErrWSConnectionClosed  = 6001 // WebSocket連接關閉
	ErrWSMessageSendFailed = 6002 // WebSocket消息發送失敗
	ErrWSInvalidMessage    = 6003 // 無效的WebSocket消息
)

// 系統錯誤 (9000-9999)
const (
	ErrInternalServer      = 9000 // 內部伺服器錯誤
	ErrDatabaseConnection  = 9001 // 資料庫連接錯誤
	ErrRequestTimeout      = 9002 // 請求超時
	ErrResourceUnavailable = 9003 // 資源不可用
	ErrUnexpectedSystem    = 9999 // 未預期的系統錯誤
)
