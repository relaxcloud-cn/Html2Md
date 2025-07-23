package service

import (
	"runtime"
	"time"

	"github.com/relaxcloud-cn/html2md/internal/model"
	"github.com/relaxcloud-cn/html2md/pkg/converter"
)

// ConvertService 转换服务
type ConvertService struct {
	converter *converter.Converter
	startTime time.Time
}

// NewConvertService 创建转换服务
func NewConvertService() *ConvertService {
	return &ConvertService{
		converter: converter.New(),
		startTime: time.Now(),
	}
}

// Convert 转换HTML为Markdown
func (s *ConvertService) Convert(req *model.ConvertRequest) (*model.ConvertResponse, error) {
	return s.converter.Convert(req)
}

// ConvertBatch 批量转换HTML为Markdown
func (s *ConvertService) ConvertBatch(req *model.BatchConvertRequest) (*model.BatchConvertResponse, error) {
	return s.converter.ConvertBatch(req)
}

// ConvertFromURL 从URL转换HTML为Markdown (TODO: 实现)
func (s *ConvertService) ConvertFromURL(req *model.ConvertFromURLRequest) (*model.ConvertFromURLResponse, error) {
	// TODO: 实现从URL获取HTML内容并转换
	// 这里需要HTTP客户端来获取URL内容
	return nil, nil
}

// Health 健康检查
func (s *ConvertService) Health() (*model.HealthResponse, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(s.startTime)

	response := &model.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Memory: &model.MemInfo{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
		},
	}

	return response, nil
}

// GetConverterInfo 获取转换器信息
func (s *ConvertService) GetConverterInfo() map[string]interface{} {
	return s.converter.GetConverterInfo()
}
