package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
* Http request response common util
 */
const (
	SUCCESS        int = 200
	ERROR          int = 500
	INVALID_PARAMS int = 400
	UNAUTHORIZED   int = 401
	FORBIDDEN      int = 403
	NOT_FOUND      int = 404
)

// response message
var ResMsg = map[int]string{
	SUCCESS:        "成功",
	ERROR:          "服务器内部错误",
	INVALID_PARAMS: "请求参数错误",
	UNAUTHORIZED:   "未授权访问",
	FORBIDDEN:      "禁止访问",
	NOT_FOUND:      "资源不存在",
}

// Response struct
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// http response success
func Success(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK, Response{
			Code:    SUCCESS,
			Message: ResMsg[SUCCESS],
			Data:    data,
		})
}

// http response success
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(
		http.StatusOK, Response{
			Code:    SUCCESS,
			Message: message,
			Data:    data,
		})
}

// http response fail
func Fail(c *gin.Context, code int, data interface{}) {
	c.JSON(
		http.StatusOK, Response{
			Code:    code,
			Message: ResMsg[code],
			Data:    data,
		})
}

// http response fail
func FailWithMessage(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(
		http.StatusOK, Response{
			Code:    code,
			Message: message,
			Data:    data,
		})
}

// get response message
func GetMsg(code int) string {
	msg, ok := ResMsg[code]
	if ok {
		return msg
	} else {
		return ResMsg[ERROR]
	}
}

// get http status code from the business code
func getHttpStatusByCode(code int) int {
	switch code {
	case INVALID_PARAMS:
		return http.StatusBadRequest
	case UNAUTHORIZED:
		return http.StatusUnauthorized
	case FORBIDDEN:
		return http.StatusForbidden
	case NOT_FOUND:
		return http.StatusNotFound
	case ERROR:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
