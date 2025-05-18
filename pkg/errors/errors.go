// Package errors 提供增強的錯誤處理功能
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error 自定義錯誤接口
type Error interface {
	error
	Code() int                             // 返回錯誤碼
	Key() string                           // 返回錯誤鍵名
	Message() string                       // 返回用戶友好訊息
	StatusCode() int                       // 返回HTTP狀態碼
	WithError(err error) Error             // 添加原始錯誤
	WithDetails(details interface{}) Error // 添加詳細信息
	LogError()                             // 記錄錯誤
}

// AppError 應用錯誤結構
type AppError struct {
	code       int         // 錯誤碼
	key        string      // 錯誤鍵名
	message    string      // 用戶友好訊息
	devMessage string      // 開發者訊息
	origError  error       // 原始錯誤
	statusCode int         // HTTP狀態碼
	details    interface{} // 詳細信息
	stack      string      // 堆疊追蹤
	logged     bool        // 是否已被記錄
}

// Error 實現 error 接口
func (e *AppError) Error() string {
	if e.origError != nil {
		return fmt.Sprintf("[%d|%s] %s: %v", e.code, e.key, e.message, e.origError)
	}
	return fmt.Sprintf("[%d|%s] %s", e.code, e.key, e.message)
}

// 實現 Error 接口方法
func (e *AppError) Code() int       { return e.code }
func (e *AppError) Key() string     { return e.key }
func (e *AppError) Message() string { return e.message }
func (e *AppError) StatusCode() int { return e.statusCode }

// WithError 添加原始錯誤
func (e *AppError) WithError(err error) Error {
	e.origError = err
	return e
}

// WithDetails 添加詳細信息
func (e *AppError) WithDetails(details interface{}) Error {
	e.details = details
	return e
}

// WithDevMessage 添加開發者訊息
func (e *AppError) WithDevMessage(msg string) Error {
	e.devMessage = msg
	return e
}

// WithStatusCode 設置HTTP狀態碼
func (e *AppError) WithStatusCode(code int) Error {
	e.statusCode = code
	return e
}

// DevMessage 獲取開發者訊息
func (e *AppError) DevMessage() string {
	return e.devMessage
}

// Details 獲取詳細信息
func (e *AppError) Details() interface{} {
	return e.details
}

// OriginalError 獲取原始錯誤
func (e *AppError) OriginalError() error {
	return e.origError
}

// Stack 獲取堆疊追蹤
func (e *AppError) Stack() string {
	return e.stack
}

// Unwrap 實現 Unwrap 接口，支持 errors.Is 和 errors.As
func (e *AppError) Unwrap() error {
	return e.origError
}

// LogError 方法在 log.go 中實現

// New 創建新的應用錯誤
func New(code int, key, message string) Error {
	return &AppError{
		code:       code,
		key:        key,
		message:    message,
		statusCode: mapCodeToStatus(code),
		stack:      captureStack(2),
	}
}

// Wrap 封裝現有錯誤
func Wrap(err error, code int, key, message string) Error {
	if err == nil {
		return New(code, key, message)
	}

	// 檢查是否已是 AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		// 保留原始AppError的訊息，但使用新的碼和鍵名
		return &AppError{
			code:       code,
			key:        key,
			message:    message,
			devMessage: appErr.devMessage,
			origError:  appErr.origError, // 使用最初的錯誤
			statusCode: mapCodeToStatus(code),
			details:    appErr.details,
			stack:      captureStack(2),
		}
	}

	return &AppError{
		code:       code,
		key:        key,
		message:    message,
		origError:  err,
		statusCode: mapCodeToStatus(code),
		stack:      captureStack(2),
	}
}

// IsAppError 檢查錯誤是否為應用錯誤
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError 從錯誤中提取應用錯誤
func GetAppError(err error) (Error, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetErrorCode 獲取錯誤碼
func GetErrorCode(err error) int {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Code()
	}
	return 0 // 無錯誤碼
}

// 捕獲堆疊信息
func captureStack(skip int) string {
	var buf [2 << 10]byte
	n := runtime.Stack(buf[:], false)
	stack := string(buf[:n])

	lines := strings.Split(stack, "\n")
	if len(lines) > skip*2+2 {
		// 跳過前幾行，從實際調用點開始
		return strings.Join(lines[skip*2:], "\n")
	}
	return stack
}
