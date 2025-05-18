package errors

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel 定義日誌級別
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// 預設日誌記錄器
var defaultLogger = log.New(os.Stderr, "", log.LstdFlags)

// 自定義日誌記錄函數類型
type LoggerFunc func(level LogLevel, msg string, fields map[string]interface{})

// 當前使用的日誌記錄函數
var loggerFunc LoggerFunc = defaultLoggerFunc

// 設置自定義日誌記錄器
func SetLogger(fn LoggerFunc) {
	if fn != nil {
		loggerFunc = fn
	}
}

// LogError 記錄應用錯誤
func (e *AppError) LogError() {
	// 避免重複記錄
	if e.logged {
		return
	}

	fields := map[string]interface{}{
		"error_code":    e.code,
		"error_key":     e.key,
		"error_message": e.message,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	if e.devMessage != "" {
		fields["dev_message"] = e.devMessage
	}

	if e.origError != nil {
		fields["original_error"] = e.origError.Error()
	}

	if e.details != nil {
		fields["details"] = e.details
	}

	if e.stack != "" {
		fields["stack"] = e.stack
	}

	// 根據錯誤碼選擇日誌級別
	level := selectLogLevel(e.code)

	// 記錄錯誤
	loggerFunc(level, e.Error(), fields)

	// 標記為已記錄
	e.logged = true
}

// 預設日誌記錄實現
func defaultLoggerFunc(level LogLevel, msg string, fields map[string]interface{}) {
	levelStr := "INFO"
	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
	case LevelInfo:
		levelStr = "INFO"
	case LevelWarn:
		levelStr = "WARN"
	case LevelError:
		levelStr = "ERROR"
	case LevelFatal:
		levelStr = "FATAL"
	}

	// 構建日誌訊息
	fieldsStr := ""
	for k, v := range fields {
		fieldsStr += fmt.Sprintf(" %s=%v", k, v)
	}

	// 打印日誌
	defaultLogger.Printf("[%s] %s%s", levelStr, msg, fieldsStr)
}

// 根據錯誤碼選擇日誌級別
func selectLogLevel(code int) LogLevel {
	// 基於錯誤碼的第一位數字確定日誌級別
	switch code / 1000 {
	case 1, 2: // 輸入、身份驗證錯誤 - 一般是正常操作
		return LevelInfo
	case 3, 4, 5, 6, 7: // 業務邏輯錯誤
		return LevelWarn
	case 8: // WebSocket錯誤
		return LevelError
	case 9: // 系統錯誤
		return LevelError
	default:
		return LevelInfo
	}
}
