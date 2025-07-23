package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "github.com/relaxcloud-cn/html2md/api/grpc/proto"
	"github.com/relaxcloud-cn/html2md/api/grpc/server"
	httpApi "github.com/relaxcloud-cn/html2md/api/http"
	"github.com/relaxcloud-cn/html2md/internal/config"
)

// @title HTML2Markdown API
// @version 1.0
// @description HTML转Markdown转换服务的REST API，支持多种转换选项和插件
// @termsOfService https://github.com/relaxcloud-cn/html2md
// @contact.name API Support
// @contact.url https://github.com/relaxcloud-cn/html2md/issues
// @contact.email support@example.com
// @license.name MIT
// @license.url https://github.com/relaxcloud-cn/html2md/blob/main/LICENSE
// @host localhost:8080
// @BasePath /

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	log.Printf("HTML2Markdown 服务启动中...")
	log.Printf("环境: %s", cfg.Server.Environment)
	log.Printf("版本: %s", cfg.Server.Version)
	log.Printf("HTTP 服务地址: %s", cfg.GetHTTPAddress())
	log.Printf("GRPC 服务地址: %s", cfg.GetGRPCAddress())

	// 创建上下文和等待组
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// 启动HTTP服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startHTTPServer(ctx, cfg); err != nil {
			log.Printf("HTTP服务器错误: %v", err)
		}
	}()

	// 启动GRPC服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startGRPCServer(ctx, cfg); err != nil {
			log.Printf("GRPC服务器错误: %v", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("HTML2Markdown 服务已启动，按 Ctrl+C 停止服务")

	// 等待中断信号
	<-c
	log.Println("收到停止信号，正在关闭服务...")

	// 取消上下文
	cancel()

	// 等待所有服务器停止
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 等待服务器关闭或超时
	select {
	case <-done:
		log.Println("服务器已安全关闭")
	case <-time.After(30 * time.Second):
		log.Println("服务器关闭超时")
	}
}

// startHTTPServer 启动HTTP服务器
func startHTTPServer(ctx context.Context, cfg *config.Config) error {
	// 创建路由器
	router := httpApi.NewRouter()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         cfg.GetHTTPAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	}

	// 在goroutine中启动服务器
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP服务器启动在: %s", cfg.GetHTTPAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("HTTP服务器启动失败: %w", err)
		}
	}()

	// 等待上下文取消或服务器错误
	select {
	case <-ctx.Done():
		log.Println("正在关闭HTTP服务器...")

		// 创建关闭上下文
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		// 优雅关闭
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP服务器关闭错误: %v", err)
			return err
		}

		log.Println("HTTP服务器已关闭")
		return nil

	case err := <-serverErr:
		return err
	}
}

// startGRPCServer 启动GRPC服务器
func startGRPCServer(ctx context.Context, cfg *config.Config) error {
	// 创建监听器
	lis, err := net.Listen("tcp", cfg.GetGRPCAddress())
	if err != nil {
		return fmt.Errorf("GRPC监听失败: %w", err)
	}

	// 创建GRPC服务器
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.Server.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.Server.GRPC.MaxSendMsgSize),
	)

	// 注册服务
	convertServer := server.NewConvertServer()
	pb.RegisterConvertServiceServer(grpcServer, convertServer)

	// 在goroutine中启动服务器
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("GRPC服务器启动在: %s", cfg.GetGRPCAddress())
		if err := grpcServer.Serve(lis); err != nil {
			serverErr <- fmt.Errorf("GRPC服务器启动失败: %w", err)
		}
	}()

	// 等待上下文取消或服务器错误
	select {
	case <-ctx.Done():
		log.Println("正在关闭GRPC服务器...")

		// 优雅关闭GRPC服务器
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		// 等待优雅关闭或强制关闭
		select {
		case <-stopped:
			log.Println("GRPC服务器已优雅关闭")
		case <-time.After(10 * time.Second):
			log.Println("GRPC服务器优雅关闭超时，强制关闭")
			grpcServer.Stop()
		}

		return nil

	case err := <-serverErr:
		return err
	}
}
