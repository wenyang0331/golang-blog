package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Message)
		return
	}

	// 未知错误，记录日志但不暴露给客户端
	Error(c, 500, "Internal server error")
}

// HandleValidationError 处理 gin 参数校验错误，返回具体字段错误信息
func HandleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errMap := make(map[string]string)
		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				errMap[field] = field + " is required"
			case "min":
				errMap[field] = field + " must be at least " + fe.Param() + " characters"
			case "max":
				errMap[field] = field + " must be at most " + fe.Param() + " characters"
			case "email":
				errMap[field] = field + " must be a valid email address"
			default:
				errMap[field] = field + " is invalid"
			}
		}
		ValidationError(c, errMap)
		return
	}
	// 非 validator 错误（如 JSON 格式错误），走通用处理
	Error(c, http.StatusBadRequest, "Invalid request format")
}
