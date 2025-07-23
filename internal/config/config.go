package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server ServerConfig `json:"server"`

	// 日志配置
	Log LogConfig `json:"log"`

	// 转换器配置
	Converter ConverterConfig `json:"converter"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// HTTP服务器配置
	HTTP HTTPConfig `json:"http"`

	// GRPC服务器配置
	GRPC GRPCConfig `json:"grpc"`

	// 通用配置
	Environment string        `json:"environment"` // 环境: development, production
	Name        string        `json:"name"`        // 服务名称
	Version     string        `json:"version"`     // 服务版本
	Timeout     time.Duration `json:"timeout"`     // 请求超时时间
}

// HTTPConfig HTTP服务器配置
type HTTPConfig struct {
	Port         int           `json:"port"`          // 端口
	Host         string        `json:"host"`          // 主机地址
	ReadTimeout  time.Duration `json:"read_timeout"`  // 读取超时
	WriteTimeout time.Duration `json:"write_timeout"` // 写入超时
	IdleTimeout  time.Duration `json:"idle_timeout"`  // 空闲超时
}

// GRPCConfig GRPC服务器配置
type GRPCConfig struct {
	Port           int           `json:"port"`              // 端口
	Host           string        `json:"host"`              // 主机地址
	MaxRecvMsgSize int           `json:"max_recv_msg_size"` // 最大接收消息大小
	MaxSendMsgSize int           `json:"max_send_msg_size"` // 最大发送消息大小
	Timeout        time.Duration `json:"timeout"`           // 超时时间
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level"`  // 日志级别: debug, info, warn, error
	Format string `json:"format"` // 日志格式: json, text
	Output string `json:"output"` // 输出: stdout, stderr, file
	File   string `json:"file"`   // 日志文件路径
}

// ConverterConfig 转换器配置
type ConverterConfig struct {
	MaxInputSize   int           `json:"max_input_size"`  // 最大输入大小（字节）
	MaxBatchSize   int           `json:"max_batch_size"`  // 最大批量转换数量
	DefaultPlugins []string      `json:"default_plugins"` // 默认启用的插件
	Timeout        time.Duration `json:"timeout"`         // 转换超时时间
	EnableCache    bool          `json:"enable_cache"`    // 是否启用缓存
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			HTTP: HTTPConfig{
				Port:         getEnvAsInt("HTTP_PORT", 8080),
				Host:         getEnvAsString("HTTP_HOST", "0.0.0.0"),
				ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", "30s"),
				WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", "30s"),
				IdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT", "60s"),
			},
			GRPC: GRPCConfig{
				Port:           getEnvAsInt("GRPC_PORT", 9090),
				Host:           getEnvAsString("GRPC_HOST", "0.0.0.0"),
				MaxRecvMsgSize: getEnvAsInt("GRPC_MAX_RECV_MSG_SIZE", 4*1024*1024), // 4MB
				MaxSendMsgSize: getEnvAsInt("GRPC_MAX_SEND_MSG_SIZE", 4*1024*1024), // 4MB
				Timeout:        getEnvAsDuration("GRPC_TIMEOUT", "30s"),
			},
			Environment: getEnvAsString("ENVIRONMENT", "development"),
			Name:        getEnvAsString("SERVICE_NAME", "html2md"),
			Version:     getEnvAsString("SERVICE_VERSION", "1.0.0"),
			Timeout:     getEnvAsDuration("SERVICE_TIMEOUT", "30s"),
		},
		Log: LogConfig{
			Level:  getEnvAsString("LOG_LEVEL", "info"),
			Format: getEnvAsString("LOG_FORMAT", "json"),
			Output: getEnvAsString("LOG_OUTPUT", "stdout"),
			File:   getEnvAsString("LOG_FILE", ""),
		},
		Converter: ConverterConfig{
			MaxInputSize:   getEnvAsInt("CONVERTER_MAX_INPUT_SIZE", 10*1024*1024), // 10MB
			MaxBatchSize:   getEnvAsInt("CONVERTER_MAX_BATCH_SIZE", 100),
			DefaultPlugins: getEnvAsStringSlice("CONVERTER_DEFAULT_PLUGINS", []string{"base", "commonmark"}),
			Timeout:        getEnvAsDuration("CONVERTER_TIMEOUT", "30s"),
			EnableCache:    getEnvAsBool("CONVERTER_ENABLE_CACHE", false),
		},
	}

	return config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证端口
	if c.Server.HTTP.Port <= 0 || c.Server.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTP.Port)
	}

	if c.Server.GRPC.Port <= 0 || c.Server.GRPC.Port > 65535 {
		return fmt.Errorf("invalid GRPC port: %d", c.Server.GRPC.Port)
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	// 验证环境
	validEnvs := map[string]bool{
		"development": true, "production": true, "testing": true,
	}
	if !validEnvs[c.Server.Environment] {
		return fmt.Errorf("invalid environment: %s", c.Server.Environment)
	}

	// 验证转换器配置
	if c.Converter.MaxInputSize <= 0 {
		return fmt.Errorf("max input size must be positive")
	}

	if c.Converter.MaxBatchSize <= 0 {
		return fmt.Errorf("max batch size must be positive")
	}

	return nil
}

// GetHTTPAddress 获取HTTP服务器地址
func (c *Config) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.HTTP.Host, c.Server.HTTP.Port)
}

// GetGRPCAddress 获取GRPC服务器地址
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.GRPC.Host, c.Server.GRPC.Port)
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// 辅助函数

// getEnvAsString 获取环境变量（字符串）
func getEnvAsString(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量（整数）
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量（布尔）
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration 获取环境变量（时间）
func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Second * 30 // 回退默认值
}

// getEnvAsStringSlice 获取环境变量（字符串切片）
func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// 简单的逗号分割，实际项目中可能需要更复杂的解析
		return []string{value}
	}
	return defaultValue
}
