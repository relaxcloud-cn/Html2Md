package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"

	"github.com/relaxcloud-cn/html2md/internal/model"
)

// Converter HTML转Markdown转换器
type Converter struct {
	// 使用html-to-markdown v2的转换器
}

// New 创建新的转换器实例
func New() *Converter {
	return &Converter{}
}

// Convert 转换HTML为Markdown
func (c *Converter) Convert(req *model.ConvertRequest) (*model.ConvertResponse, error) {
	startTime := time.Now()

	// 验证HTML内容
	if err := c.ValidateHTML(req.HTML); err != nil {
		return nil, err
	}

	// 配置插件
	plugins := []converter.Plugin{
		base.NewBasePlugin(),
		commonmark.NewCommonmarkPlugin(),
	}

	// 添加额外的插件
	if len(req.Plugins) > 0 {
		additionalPlugins, err := c.configurePlugins(req.Plugins)
		if err != nil {
			return nil, fmt.Errorf("配置插件失败: %w", err)
		}
		plugins = append(plugins, additionalPlugins...)
	}

	// 创建转换器实例
	conv := converter.NewConverter(
		converter.WithPlugins(plugins...),
	)

	// 创建转换选项
	var options []converter.ConvertOptionFunc

	// 添加域名配置
	if req.Domain != "" {
		options = append(options, converter.WithDomain(req.Domain))
	}

	// 执行转换 - 使用转换器实例的ConvertString方法
	markdown, err := conv.ConvertString(req.HTML, options...)
	if err != nil {
		return nil, fmt.Errorf("HTML转换失败: %w", err)
	}

	// 统计信息
	processingTime := time.Since(startTime)
	stats := &model.ConversionStats{
		InputSize:      len(req.HTML),
		OutputSize:     len(markdown),
		ProcessingTime: processingTime,
		PluginsUsed:    req.Plugins,
	}

	response := &model.ConvertResponse{
		Markdown: markdown,
		Stats:    stats,
	}

	return response, nil
}

// ConvertBatch 批量转换HTML为Markdown
func (c *Converter) ConvertBatch(req *model.BatchConvertRequest) (*model.BatchConvertResponse, error) {
	startTime := time.Now()

	results := make([]model.BatchConvertItem, len(req.Items))
	successCount := 0
	failedCount := 0

	for i, item := range req.Items {
		result, err := c.Convert(&item)
		if err != nil {
			results[i] = model.BatchConvertItem{
				Index:   i,
				Success: false,
				Error:   err.Error(),
			}
			failedCount++
		} else {
			results[i] = model.BatchConvertItem{
				Index:   i,
				Success: true,
				Result:  result,
			}
			successCount++
		}
	}

	totalTime := time.Since(startTime)
	averageTime := totalTime / time.Duration(len(req.Items))

	summary := &model.BatchSummary{
		Total:       len(req.Items),
		Success:     successCount,
		Failed:      failedCount,
		TotalTime:   totalTime,
		AverageTime: averageTime,
	}

	return &model.BatchConvertResponse{
		Results: results,
		Summary: summary,
	}, nil
}

// configurePlugins 配置插件
func (c *Converter) configurePlugins(plugins []string) ([]converter.Plugin, error) {
	var result []converter.Plugin

	for _, plugin := range plugins {
		switch strings.ToLower(plugin) {
		case "commonmark", "base":
			// 已经默认加载，跳过
			continue
		default:
			return nil, fmt.Errorf("不支持的插件: %s", plugin)
		}
	}

	return result, nil
}

// GetSupportedPlugins 获取支持的插件列表
func (c *Converter) GetSupportedPlugins() []string {
	return []string{
		"base",
		"commonmark",
	}
}

// ValidateHTML 验证HTML内容
func (c *Converter) ValidateHTML(html string) error {
	if strings.TrimSpace(html) == "" {
		return fmt.Errorf("HTML内容不能为空")
	}

	// 基本的HTML标签检查
	if !strings.Contains(html, "<") || !strings.Contains(html, ">") {
		return fmt.Errorf("输入内容不是有效的HTML格式")
	}

	return nil
}

// GetConverterInfo 获取转换器信息
func (c *Converter) GetConverterInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":           "2.3.3", // html-to-markdown版本
		"supported_plugins": c.GetSupportedPlugins(),
		"features": []string{
			"CommonMark支持",
			"代码块转换",
			"链接和图片处理",
			"自定义规则支持",
		},
	}
}
