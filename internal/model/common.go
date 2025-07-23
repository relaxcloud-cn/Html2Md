package model

// APIResponse 统一API响应格式
type APIResponse struct {
	Code int         `json:"code" example:"200"`    // 业务响应码
	Msg  string      `json:"msg" example:"success"` // 接口消息
	Data interface{} `json:"data"`                  // 返回数据
}

// APIError 错误响应格式
type APIError struct {
	Code int         `json:"code" example:"400"`        // 错误码
	Msg  string      `json:"msg" example:"Bad Request"` // 错误消息
	Data interface{} `json:"data"`                      // 错误详情
}

// 定义业务响应码常量
const (
	// 成功码
	CodeSuccess = 200

	// 客户端错误码
	CodeBadRequest     = 400
	CodeUnauthorized   = 401
	CodeForbidden      = 403
	CodeNotFound       = 404
	CodeTooManyRequest = 429

	// 服务端错误码
	CodeInternalError = 500
	CodeServiceError  = 502
	CodeTimeout       = 504
)

// 定义响应消息常量
const (
	MsgSuccess        = "success"
	MsgBadRequest     = "请求参数错误"
	MsgUnauthorized   = "未授权访问"
	MsgForbidden      = "禁止访问"
	MsgNotFound       = "资源不存在"
	MsgTooManyRequest = "请求过于频繁"
	MsgInternalError  = "内部服务器错误"
	MsgServiceError   = "服务暂时不可用"
	MsgTimeout        = "请求超时"
)

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, msg string, data interface{}) APIResponse {
	return APIResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
