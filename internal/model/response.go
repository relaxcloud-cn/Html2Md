package model

import "time"

// ConvertResponse HTML转Markdown响应数据
type ConvertResponse struct {
	Markdown string           `json:"markdown" example:"# Hello World"` // 转换后的Markdown内容
	Stats    *ConversionStats `json:"stats,omitempty"`                  // 转换统计信息
}

// ConversionStats 转换统计信息
type ConversionStats struct {
	InputSize      int           `json:"input_size" example:"1024"`                            // 输入HTML大小（字节）
	OutputSize     int           `json:"output_size" example:"512"`                            // 输出Markdown大小（字节）
	ProcessingTime time.Duration `json:"processing_time" swaggertype:"string" example:"100ms"` // 处理时间
	PluginsUsed    []string      `json:"plugins_used" example:"base,commonmark"`               // 使用的插件列表
}

// HealthResponse 健康检查响应数据
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`                      // 服务状态
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T12:00:00Z"` // 检查时间
	Version   string    `json:"version" example:"1.0.0"`                  // 服务版本
	Uptime    string    `json:"uptime" example:"1h30m"`                   // 运行时间
	Memory    *MemInfo  `json:"memory,omitempty"`                         // 内存信息
}

// MemInfo 内存信息
type MemInfo struct {
	Alloc      uint64 `json:"alloc" example:"1048576"`       // 已分配内存（字节）
	TotalAlloc uint64 `json:"total_alloc" example:"2097152"` // 总分配内存（字节）
	Sys        uint64 `json:"sys" example:"4194304"`         // 系统内存（字节）
	NumGC      uint32 `json:"num_gc" example:"5"`            // GC次数
}

// BatchConvertResponse 批量转换响应数据
type BatchConvertResponse struct {
	Results []BatchConvertItem `json:"results"` // 批量转换结果
	Summary *BatchSummary      `json:"summary"` // 批量转换摘要
}

// BatchConvertItem 批量转换项目结果
type BatchConvertItem struct {
	Index   int              `json:"index" example:"0"`              // 项目索引
	Success bool             `json:"success" example:"true"`         // 是否成功
	Result  *ConvertResponse `json:"result,omitempty"`               // 转换结果（成功时）
	Error   string           `json:"error,omitempty" example:"转换失败"` // 错误信息（失败时）
}

// BatchSummary 批量转换摘要
type BatchSummary struct {
	Total       int           `json:"total" example:"10"`                                // 总数
	Success     int           `json:"success" example:"8"`                               // 成功数
	Failed      int           `json:"failed" example:"2"`                                // 失败数
	TotalTime   time.Duration `json:"total_time" swaggertype:"string" example:"5s"`      // 总处理时间
	AverageTime time.Duration `json:"average_time" swaggertype:"string" example:"500ms"` // 平均处理时间
}
