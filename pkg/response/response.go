package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 統一的響應結構
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功響應
func Success(data interface{}) gin.H {
	return gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	}
}

// Error 返回錯誤響應
func Error(message string) gin.H {
	return gin.H{
		"code":    400,
		"message": message,
	}
}

// ServerError 返回服務器錯誤響應
func ServerError(message string) gin.H {
	return gin.H{
		"code":    500,
		"message": message,
	}
}

// Unauthorized 返回未授權響應
func Unauthorized(message string) gin.H {
	return gin.H{
		"code":    401,
		"message": message,
	}
}

// NotFound 返回未找到響應
func NotFound(message string) gin.H {
	return gin.H{
		"code":    404,
		"message": message,
	}
}

// JSON 返回 JSON 響應
func JSON(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}
