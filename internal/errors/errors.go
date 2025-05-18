// Package errors 提供應用特定的錯誤處理功能
package errors

import (
	"clean-architecture-gochat/internal/common/enum"
	baseErrors "clean-architecture-gochat/pkg/errors"

	pkgErrors "github.com/pkg/errors"
)

// 創建帶有錯誤碼的應用錯誤
func New(code enum.ErrorCode, details ...interface{}) baseErrors.Error {
	key, message := enum.GetErrorDetails(code)

	appErr := baseErrors.New(int(code), key, message)

	if len(details) > 0 && details[0] != nil {
		appErr = appErr.WithDetails(details[0])
	}

	return appErr
}

// 封裝現有錯誤
func Wrap(err error, code enum.ErrorCode, details ...interface{}) baseErrors.Error {
	if err == nil {
		return New(code, details...)
	}

	key, message := enum.GetErrorDetails(code)

	// 使用 pkg/errors.Cause 確保獲取根本錯誤
	cause := pkgErrors.Cause(err)
	appErr := baseErrors.Wrap(cause, int(code), key, message)

	if len(details) > 0 && details[0] != nil {
		appErr = appErr.WithDetails(details[0])
	}

	return appErr
}

// WithDevMessage 為錯誤添加開發者訊息
func WithDevMessage(err error, devMessage string) baseErrors.Error {
	if appErr, ok := baseErrors.GetAppError(err); ok {
		return appErr.(*baseErrors.AppError).WithDevMessage(devMessage)
	}

	// 如果不是應用錯誤，先封裝成系統錯誤
	return Wrap(err, enum.ErrInternalServer).(*baseErrors.AppError).WithDevMessage(devMessage)
}

// ToResponse 將錯誤轉換為HTTP響應
func ToResponse(err error) (int, interface{}) {
	return baseErrors.ToResponse(err)
}

// 以下是常見錯誤的便捷創建函數

// NewInvalidInput 創建無效輸入錯誤
func NewInvalidInput(details ...interface{}) baseErrors.Error {
	return New(enum.ErrInvalidInput, details...)
}

// NewUnauthorized 創建未授權錯誤
func NewUnauthorized(details ...interface{}) baseErrors.Error {
	return New(enum.ErrUnauthorized, details...)
}

// NewAccessDenied 創建訪問拒絕錯誤
func NewAccessDenied(details ...interface{}) baseErrors.Error {
	return New(enum.ErrAccessDenied, details...)
}

// NewNotFound 創建資源不存在錯誤
func NewNotFound(resource string, details ...interface{}) baseErrors.Error {
	var code enum.ErrorCode

	switch resource {
	case "user":
		code = enum.ErrUserNotFound
	case "group":
		code = enum.ErrGroupNotFound
	case "file":
		code = enum.ErrFileNotFound
	default:
		code = enum.ErrInternalServer
	}

	return New(code, details...)
}

// NewInternalError 創建內部服務器錯誤
func NewInternalError(err error, details ...interface{}) baseErrors.Error {
	return Wrap(err, enum.ErrInternalServer, details...)
}

// NewDBError 創建數據庫錯誤
func NewDBError(err error, details ...interface{}) baseErrors.Error {
	return Wrap(err, enum.ErrDatabaseError, details...)
}
