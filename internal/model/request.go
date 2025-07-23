package model

// ConvertRequest HTML转Markdown请求参数
type ConvertRequest struct {
	HTML string `json:"html" binding:"required" example:"<h1>Hello World</h1>"` // HTML内容
}

// HealthRequest 健康检查请求
type HealthRequest struct {
	// 可以为空，用于扩展
}

// BatchConvertRequest 批量转换请求
type BatchConvertRequest struct {
	Items []ConvertRequest `json:"items" binding:"required,min=1,max=100"` // 批量转换项目，最多100个
}
