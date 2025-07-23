# HTML2Markdown API 服务

一个高性能的HTML转Markdown转换服务，提供HTTP REST API和GRPC接口，基于[html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown)库构建。

## 🚀 特性

- **多协议支持**: 同时提供HTTP REST API和GRPC接口
- **高性能转换**: 基于html-to-markdown v2库，支持CommonMark规范
- **插件系统**: 支持表格、删除线等扩展插件
- **批量处理**: 支持批量HTML转换，提高处理效率
- **统一响应格式**: HTTP接口采用`{code, msg, data}`统一响应格式
- **自动文档**: 集成Swagger UI，自动生成API文档
- **健康检查**: 内置服务健康检查和监控接口
- **配置灵活**: 支持环境变量配置
- **优雅关闭**: 支持服务的优雅启动和关闭

## 🛠 技术栈

- **后端框架**: Gin (HTTP) + gRPC
- **HTML转换**: [html-to-markdown v2](https://github.com/JohannesKaufmann/html-to-markdown)
- **文档生成**: Swagger/OpenAPI 3.0
- **配置管理**: 环境变量 + 配置文件
- **协议缓冲**: Protocol Buffers

## 📦 安装与使用

### 方式1: 直接运行

```bash
# 克隆项目
git clone https://github.com/relaxcloud-cn/html2md.git
cd html2md

# 安装依赖
make deps

# 构建项目
make build

# 运行服务
./bin/html2md
```

### 方式2: 开发模式

```bash
# 安装开发工具
make install-tools

# 生成协议和文档
make proto
make swagger

# 开发模式运行
make dev
```

### 方式3: Docker

```bash
# 方式3.1: 单容器运行
make docker          # 构建镜像
make docker-run      # 运行容器

# 方式3.2: Docker Compose（简化版）
make docker-compose-simple    # 启动服务
make docker-compose-simple-down  # 停止服务

# 方式3.3: Docker Compose（完整版，包含nginx、redis、监控）
make docker-compose-up        # 启动所有服务
make docker-compose-down      # 停止所有服务
make docker-compose-logs      # 查看日志
```

#### Docker访问地址

- **单容器模式**:
  - HTTP API: http://localhost:8080
  - GRPC API: localhost:9090

- **Docker Compose模式**:
  - 通过Nginx访问: http://localhost (端口80)
  - 直接访问: http://localhost:8080
  - GRPC API: localhost:9090
  - Prometheus监控: http://localhost:9091 (完整版)

## 🌐 API 接口

服务启动后，可以通过以下地址访问：

- **HTTP服务**: http://localhost:8080
- **GRPC服务**: localhost:9090  
- **Swagger文档**: http://localhost:8080/docs/index.html
- **演示页面**: http://localhost:8080/api/v1/demo

### HTTP API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/convert` | 转换HTML为Markdown |
| `GET` | `/api/v1/convert/simple` | 简单转换（GET方式） |
| `POST` | `/api/v1/convert/batch` | 批量转换 |
| `POST` | `/api/v1/convert/url` | 从URL转换 |
| `GET` | `/api/v1/health` | 健康检查 |
| `GET` | `/api/v1/info` | 转换器信息 |
| `GET` | `/api/v1/demo` | 演示页面 |

### 示例请求

#### 基本转换

```bash
curl -X POST http://localhost:8080/api/v1/convert \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<h1>Hello World</h1><p>This is a <strong>bold</strong> text.</p>",
    "plugins": ["commonmark"]
  }'
```

#### 响应格式

```json
{
  "code": 200,
  "msg": "success", 
  "data": {
    "markdown": "# Hello World\n\nThis is a **bold** text.",
    "stats": {
      "input_size": 65,
      "output_size": 42,
      "processing_time": "2.5ms",
      "elements_count": 3,
      "converted_count": 3,
      "skipped_count": 0,
      "plugins_used": ["commonmark"]
    }
  }
}
```

#### 简单转换 (GET)

```bash
curl "http://localhost:8080/api/v1/convert/simple?html=<h1>Test</h1>&plugins=commonmark"
```

## 🔧 配置

### 环境变量配置

项目根目录包含 `env.example` 文件，展示了所有可用的环境变量配置。使用前请复制为 `.env` 文件：

```bash
cp env.example .env
# 然后编辑 .env 文件
```

### 主要环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `HTTP_PORT` | `8080` | HTTP服务端口 |
| `GRPC_PORT` | `9090` | GRPC服务端口 |
| `ENVIRONMENT` | `development` | 运行环境 |
| `LOG_LEVEL` | `info` | 日志级别 |
| `CONVERTER_MAX_INPUT_SIZE` | `10485760` | 最大输入大小(10MB) |
| `CONVERTER_MAX_BATCH_SIZE` | `100` | 最大批量数量 |

### 配置示例

```bash
# 设置端口
export HTTP_PORT=8888
export GRPC_PORT=9999

# 设置日志级别
export LOG_LEVEL=debug

# 启动服务
./bin/html2md
```

## 🐳 Docker 部署

> 📖 **完整部署指南**: 查看 [Docker部署指南](docs/DOCKER_DEPLOYMENT.md) 获取详细说明和故障排查

### Docker 镜像特性

- **多阶段构建**: 优化镜像大小，最终镜像约20MB
- **非root用户**: 使用专用用户运行，提高安全性
- **健康检查**: 内置健康检查机制
- **时区设置**: 默认设置为Asia/Shanghai
- **资源限制**: 配置了合理的CPU和内存限制

### 部署选项

#### 1. 单容器部署（推荐用于开发）

```bash
# 构建镜像
make docker

# 运行容器
make docker-run

# 手动运行（自定义配置）
docker run -d \
  --name html2md \
  -p 8080:8080 \
  -p 9090:9090 \
  -e LOG_LEVEL=debug \
  -e CONVERTER_MAX_INPUT_SIZE=20971520 \
  html2md:latest
```

#### 2. Docker Compose 简化部署（推荐用于生产）

```bash
# 启动服务
make docker-compose-simple

# 停止服务
make docker-compose-simple-down

# 查看日志
docker-compose -f docker-compose.simple.yml logs -f
```

#### 3. Docker Compose 完整部署（包含监控和缓存）

```bash
# 启动所有服务
make docker-compose-up

# 查看服务状态
docker-compose ps

# 停止所有服务
make docker-compose-down
```

完整部署包含的服务：
- **html2md**: 主服务
- **nginx**: 反向代理和负载均衡
- **redis**: 缓存服务（可选）
- **prometheus**: 监控服务（可选）

### Docker 环境配置

Docker Compose支持通过`.env`文件进行配置：

```bash
# 复制配置文件
cp env.example .env

# 编辑配置
vim .env

# 使用配置启动
docker-compose up -d
```

### 常用Docker命令

```bash
# 查看运行状态
docker-compose ps

# 查看实时日志
make docker-compose-logs

# 重启服务
docker-compose restart html2md

# 更新服务（重新构建）
docker-compose up -d --build

# 清理资源
make docker-clean
```

## 🔌 GRPC 客户端

### Go 客户端示例

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    pb "github.com/relaxcloud-cn/html2md/api/grpc/proto"
)

func main() {
    // 连接服务器
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // 创建客户端
    client := pb.NewConvertServiceClient(conn)
    
    // 调用转换接口
    resp, err := client.Convert(context.Background(), &pb.ConvertRequest{
        Html: "<h1>Hello GRPC</h1>",
        Plugins: []string{"commonmark"},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("转换结果: %s", resp.Markdown)
}
```

## 🎯 支持的插件

- `base` - 基础功能插件 (默认启用)
- `commonmark` - CommonMark规范插件 (默认启用)  
- `table` - 表格转换插件 (TODO)
- `strikethrough` - 删除线插件 (TODO)

## 📊 性能监控

### 健康检查

```bash
curl http://localhost:8080/api/v1/health
```

### 服务信息

```bash  
curl http://localhost:8080/api/v1/info
```

## 🏗 开发

### 项目结构

```
Html2Md/
├── cmd/server/           # 主服务入口
├── api/
│   ├── http/            # HTTP API
│   │   ├── handler/     # 处理器
│   │   ├── middleware/  # 中间件
│   │   └── router.go    # 路由
│   └── grpc/            # GRPC API
│       ├── proto/       # 协议文件
│       └── server/      # 服务器
├── internal/
│   ├── config/          # 配置管理
│   ├── service/         # 业务逻辑
│   └── model/           # 数据模型
├── pkg/converter/       # 核心转换器
├── docs/                # API文档
└── Makefile            # 构建脚本
```

### 可用命令

```bash
make help           # 查看所有命令
make build          # 构建项目
make dev            # 开发模式运行
make test           # 运行测试
make proto          # 生成protobuf代码
make swagger        # 生成swagger文档
make clean          # 清理构建文件
make fmt            # 格式化代码

# Docker相关命令
make docker                    # 构建Docker镜像
make docker-run               # 运行单个Docker容器
make docker-compose-simple    # 启动简化版Docker Compose
make docker-compose-up        # 启动完整版Docker Compose
make docker-compose-down      # 停止Docker Compose
make docker-compose-logs      # 查看Docker Compose日志
make docker-clean             # 清理Docker资源
```

### 添加新功能

1. 在 `internal/model/` 中定义数据模型
2. 在 `api/grpc/proto/` 中更新协议文件 
3. 在 `api/http/handler/` 中添加HTTP处理器
4. 在 `api/grpc/server/` 中添加GRPC处理器
5. 运行 `make proto swagger` 重新生成代码和文档

## 🤝 贡献

欢迎提交Issue和Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交变更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - 核心转换库
- [Gin](https://github.com/gin-gonic/gin) - HTTP Web框架
- [gRPC](https://grpc.io/) - 高性能RPC框架
- [Swagger](https://swagger.io/) - API文档工具

---

**⭐ 如果这个项目对你有帮助，请给它一个Star！**
