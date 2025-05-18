package errors

import (
	"net/http"
	"os"
)

// 將錯誤碼映射到HTTP狀態碼
func mapCodeToStatus(code int) int {
	// 基於錯誤碼的第一位數字來確定HTTP狀態
	switch code / 1000 {
	case 1: // 輸入錯誤
		return http.StatusBadRequest
	case 2: // 身份驗證錯誤
		if code == 2003 { // 假設2003是權限不足
			return http.StatusForbidden
		}
		return http.StatusUnauthorized
	case 3: // 資源錯誤
		if code == 3000 { // 假設3000是資源不存在
			return http.StatusNotFound
		}
		if code == 3001 { // 假設3001是資源已存在
			return http.StatusConflict
		}
		return http.StatusBadRequest
	case 4, 5, 6, 7: // 業務邏輯錯誤
		return http.StatusBadRequest
	case 8: // WebSocket錯誤
		return http.StatusServiceUnavailable
	case 9: // 系統錯誤
		if code == 9002 { // 假設9002是超時
			return http.StatusRequestTimeout
		}
		if code == 9003 { // 假設9003是服務不可用
			return http.StatusServiceUnavailable
		}
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ToResponse 將錯誤轉換為HTTP響應格式
func ToResponse(err error) (int, interface{}) {
	if err == nil {
		return http.StatusOK, map[string]interface{}{
			"success": true,
		}
	}

	// 提取應用錯誤
	appErr, ok := GetAppError(err)
	if !ok {
		// 非應用錯誤，包裝為系統錯誤
		appErr = Wrap(err, 9000, "SYSTEM_ERROR", "系統錯誤").(*AppError)
	}

	// 記錄錯誤
	appErr.LogError()

	// 構建響應
	response := map[string]interface{}{
		"success": false,
		"code":    appErr.Code(),
		"key":     appErr.Key(),
		"message": appErr.Message(),
	}

	// 添加詳細信息 (如果有)
	details := appErr.(*AppError).Details()
	if details != nil {
		response["details"] = details
	}

	// 在非生產環境添加開發者訊息
	// 此處可以根據環境變量來判斷
	if appErr.(*AppError).DevMessage() != "" && !isProduction() {
		response["dev_message"] = appErr.(*AppError).DevMessage()
	}

	return appErr.StatusCode(), response
}

// 檢查是否為生產環境
func isProduction() bool {
	// 這裡可以根據環境變量判斷
	return os.Getenv("GO_ENV") == "production"
}
