package server

import (
	"context"
	"runtime"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/relaxcloud-cn/html2md/api/grpc/proto"
	"github.com/relaxcloud-cn/html2md/internal/model"
	"github.com/relaxcloud-cn/html2md/internal/service"
)

// ConvertServer GRPC转换服务器
type ConvertServer struct {
	pb.UnimplementedConvertServiceServer
	service   *service.ConvertService
	startTime time.Time
}

// NewConvertServer 创建GRPC转换服务器
func NewConvertServer() *ConvertServer {
	return &ConvertServer{
		service:   service.NewConvertService(),
		startTime: time.Now(),
	}
}

// Convert 转换HTML为Markdown
func (s *ConvertServer) Convert(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	// 将protobuf请求转换为内部模型
	modelReq := &model.ConvertRequest{
		HTML: req.Html,
	}

	// 执行转换
	result, err := s.service.Convert(modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "转换失败: %v", err)
	}

	// 转换结果为protobuf响应
	response := &pb.ConvertResponse{
		Markdown: result.Markdown,
	}

	// 转换统计信息
	if result.Stats != nil {
		response.Stats = &pb.ConversionStats{
			InputSize:      int32(result.Stats.InputSize),
			OutputSize:     int32(result.Stats.OutputSize),
			ProcessingTime: durationpb.New(result.Stats.ProcessingTime),
		}
	}

	return response, nil
}

// ConvertBatch 批量转换HTML为Markdown
func (s *ConvertServer) ConvertBatch(ctx context.Context, req *pb.BatchConvertRequest) (*pb.BatchConvertResponse, error) {
	// 转换请求
	modelReq := &model.BatchConvertRequest{
		Items: make([]model.ConvertRequest, len(req.Items)),
	}

	for i, item := range req.Items {
		modelReq.Items[i] = model.ConvertRequest{
			HTML: item.Html,
		}
	}

	// 执行批量转换
	result, err := s.service.ConvertBatch(modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量转换失败: %v", err)
	}

	// 转换结果
	response := &pb.BatchConvertResponse{
		Results: make([]*pb.BatchConvertItem, len(result.Results)),
	}

	for i, item := range result.Results {
		response.Results[i] = &pb.BatchConvertItem{
			Index:   int32(item.Index),
			Success: item.Success,
			Error:   item.Error,
		}

		if item.Result != nil {
			// 转换成功的结果
			pbResult := &pb.ConvertResponse{
				Markdown: item.Result.Markdown,
			}

			if item.Result.Stats != nil {
				pbResult.Stats = &pb.ConversionStats{
					InputSize:      int32(item.Result.Stats.InputSize),
					OutputSize:     int32(item.Result.Stats.OutputSize),
					ProcessingTime: durationpb.New(item.Result.Stats.ProcessingTime),
				}
			}

			response.Results[i].Result = pbResult
		}
	}

	// 转换摘要信息
	if result.Summary != nil {
		response.Summary = &pb.BatchSummary{
			Total:       int32(result.Summary.Total),
			Success:     int32(result.Summary.Success),
			Failed:      int32(result.Summary.Failed),
			TotalTime:   durationpb.New(result.Summary.TotalTime),
			AverageTime: durationpb.New(result.Summary.AverageTime),
		}
	}

	return response, nil
}

// HealthCheck 健康检查
func (s *ConvertServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// 获取内存统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(s.startTime)

	response := &pb.HealthCheckResponse{
		Status:    "ok",
		Timestamp: timestamppb.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Memory: &pb.MemInfo{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGc:      m.NumGC,
		},
	}

	return response, nil
}

// GetConverterInfo 获取转换器信息
func (s *ConvertServer) GetConverterInfo(ctx context.Context, req *pb.GetConverterInfoRequest) (*pb.GetConverterInfoResponse, error) {
	info := s.service.GetConverterInfo()

	response := &pb.GetConverterInfoResponse{
		Version: info["version"].(string),
		Config:  make(map[string]string),
	}

	// 转换支持的插件
	if plugins, ok := info["supported_plugins"].([]string); ok {
		response.SupportedPlugins = plugins
	}

	// 转换功能特性
	if features, ok := info["features"].([]string); ok {
		response.Features = features
	}

	return response, nil
}
