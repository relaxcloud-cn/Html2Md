package model

import "time"

// ConvertResponse HTML转Markdown响应数据
type ConvertResponse struct {
	Markdown string           `json:"markdown" example:"# Hello World"` // 转换后的Markdown内容
	Stats    *ConversionStats `json:"stats,omitempty"`                  // 转换统计信息
	Warnings []string         `json:"warnings,omitempty"`               // 转换警告信息
	Metadata *ConversionMeta  `json:"metadata,omitempty"`               // 转换元数据
}

// ConversionStats 转换统计信息
type ConversionStats struct {
	InputSize      int           `json:"input_size" example:"1024"`                            // 输入HTML大小（字节）
	OutputSize     int           `json:"output_size" example:"512"`                            // 输出Markdown大小（字节）
	ProcessingTime time.Duration `json:"processing_time" swaggertype:"string" example:"100ms"` // 处理时间
	ElementsCount  int           `json:"elements_count" example:"25"`                          // HTML元素数量
	ConvertedCount int           `json:"converted_count" example:"23"`                         // 成功转换的元素数量
	SkippedCount   int           `json:"skipped_count" example:"2"`                            // 跳过的元素数量
	PluginsUsed    []string      `json:"plugins_used" example:"commonmark,table"`              // 使用的插件列表
}

// ConversionMeta 转换元数据
type ConversionMeta struct {
	Title       string            `json:"title,omitempty" example:"Page Title"`             // 页面标题
	Description string            `json:"description,omitempty" example:"Page description"` // 页面描述
	Keywords    []string          `json:"keywords,omitempty"`                               // 关键词
	Author      string            `json:"author,omitempty" example:"Author Name"`           // 作者
	Language    string            `json:"language,omitempty" example:"zh-CN"`               // 语言
	Images      []ImageInfo       `json:"images,omitempty"`                                 // 图片信息
	Links       []LinkInfo        `json:"links,omitempty"`                                  // 链接信息
	CustomMeta  map[string]string `json:"custom_meta,omitempty"`                            // 自定义元数据
}

// ImageInfo 图片信息
type ImageInfo struct {
	Src    string `json:"src" example:"https://example.com/image.jpg"` // 图片源地址
	Alt    string `json:"alt,omitempty" example:"Image description"`   // 替代文本
	Title  string `json:"title,omitempty" example:"Image title"`       // 图片标题
	Width  int    `json:"width,omitempty" example:"800"`               // 宽度
	Height int    `json:"height,omitempty" example:"600"`              // 高度
}

// LinkInfo 链接信息
type LinkInfo struct {
	Href   string `json:"href" example:"https://example.com"`   // 链接地址
	Text   string `json:"text" example:"Example Link"`          // 链接文本
	Title  string `json:"title,omitempty" example:"Link title"` // 链接标题
	Target string `json:"target,omitempty" example:"_blank"`    // 链接目标
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

// ConvertFromURLResponse 从URL转换响应数据
type ConvertFromURLResponse struct {
	ConvertResponse
	URL             string            `json:"url" example:"https://example.com/page.html"`  // 源URL
	FetchTime       time.Duration     `json:"fetch_time" swaggertype:"string" example:"2s"` // 获取时间
	ContentType     string            `json:"content_type" example:"text/html"`             // 内容类型
	StatusCode      int               `json:"status_code" example:"200"`                    // HTTP状态码
	FinalURL        string            `json:"final_url,omitempty"`                          // 最终URL（重定向后）
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`                   // 响应头
}
