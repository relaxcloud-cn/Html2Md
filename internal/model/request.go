package model

// ConvertRequest HTML转Markdown请求参数
type ConvertRequest struct {
	HTML    string          `json:"html" binding:"required" example:"<h1>Hello World</h1>"` // HTML内容
	Options *ConvertOptions `json:"options,omitempty"`                                      // 转换选项
	Plugins []string        `json:"plugins,omitempty" example:"commonmark,table"`           // 启用的插件
	Domain  string          `json:"domain,omitempty" example:"https://example.com"`         // 基础域名，用于转换相对链接
}

// ConvertOptions 转换选项
type ConvertOptions struct {
	// 是否移除空白字符
	TrimSpaces bool `json:"trim_spaces,omitempty" example:"true"`

	// 是否保留HTML标签（对于不支持的标签）
	KeepUnknownTags bool `json:"keep_unknown_tags,omitempty" example:"false"`

	// 链接处理选项
	LinkTarget string `json:"link_target,omitempty" example:"_blank"` // 链接目标

	// 图片处理选项
	ImageAbsolutePath bool `json:"image_absolute_path,omitempty" example:"true"` // 是否转换为绝对路径

	// 代码块处理
	CodeBlockStyle string `json:"code_block_style,omitempty" example:"fenced"` // 代码块样式: fenced 或 indented

	// 表格处理
	TableCompact bool `json:"table_compact,omitempty" example:"false"` // 是否压缩表格格式

	// 强调标记样式
	EmphasisStyle string `json:"emphasis_style,omitempty" example:"*"` // 强调样式: * 或 _
	BoldStyle     string `json:"bold_style,omitempty" example:"**"`    // 粗体样式: ** 或 __

	// 标题样式
	HeadingStyle string `json:"heading_style,omitempty" example:"atx"` // 标题样式: atx (#) 或 setext

	// 列表样式
	BulletListMarker  string `json:"bullet_list_marker,omitempty" example:"-"`  // 无序列表标记: -, *, +
	OrderedListMarker string `json:"ordered_list_marker,omitempty" example:"."` // 有序列表标记: . 或 )
}

// HealthRequest 健康检查请求
type HealthRequest struct {
	// 可以为空，用于扩展
}

// BatchConvertRequest 批量转换请求
type BatchConvertRequest struct {
	Items []ConvertRequest `json:"items" binding:"required,min=1,max=100"` // 批量转换项目，最多100个
}

// ConvertFromURLRequest 从URL转换请求
type ConvertFromURLRequest struct {
	URL     string          `json:"url" binding:"required,url" example:"https://example.com/page.html"` // 网页URL
	Options *ConvertOptions `json:"options,omitempty"`                                                  // 转换选项
	Plugins []string        `json:"plugins,omitempty"`                                                  // 启用的插件

	// 选择器选项
	IncludeSelector string `json:"include_selector,omitempty" example:"article"` // 只包含匹配的元素
	ExcludeSelector string `json:"exclude_selector,omitempty" example:".ad"`     // 排除匹配的元素

	// 请求选项
	Timeout int               `json:"timeout,omitempty" example:"30"` // 请求超时时间（秒）
	Headers map[string]string `json:"headers,omitempty"`              // 自定义请求头
}
