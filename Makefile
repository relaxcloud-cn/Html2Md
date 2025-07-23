.PHONY: help build clean test proto swagger docker run

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  build     - 构建项目"
	@echo "  clean     - 清理构建文件"
	@echo "  test      - 运行测试"
	@echo "  proto     - 生成protobuf代码"
	@echo "  swagger   - 生成swagger文档"
	@echo "  run       - 运行服务"
	@echo "  deps      - 安装依赖"
	@echo ""
	@echo "Docker相关命令:"
	@echo "  docker                    - 构建Docker镜像"
	@echo "  docker-run               - 运行单个Docker容器"
	@echo "  docker-compose-simple    - 启动简化版Docker Compose"
	@echo "  docker-compose-up        - 启动完整版Docker Compose"
	@echo "  docker-compose-down      - 停止Docker Compose"
	@echo "  docker-compose-logs      - 查看Docker Compose日志"
	@echo "  docker-clean             - 清理Docker资源"
	@echo "  docker-clean-all         - 完全清理Docker资源"

# 构建项目
build:
	go build -o bin/html2md cmd/server/main.go

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf docs/swagger.json
	rm -rf docs/swagger.yaml

# 运行测试
test:
	go test -v ./...

# 生成protobuf代码
proto:
	@echo "生成protobuf代码..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/grpc/proto/convert.proto

# 生成swagger文档
swagger:
	@echo "生成swagger文档..."
	swag init -g api/http/router.go -o docs/

# 安装依赖
deps:
	go mod download
	go mod tidy

# 运行服务
run: build
	./bin/html2md

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 安装工具
install-tools:
	@echo "安装开发工具..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 初始化项目
init: install-tools deps proto

# 开发模式运行
dev:
	go run cmd/server/main.go

# 构建Docker镜像
docker:
	docker build -t html2md:latest .

# Docker运行（单容器）
docker-run:
	docker run -p 8080:8080 -p 9090:9090 html2md:latest

# Docker Compose构建
docker-compose-build:
	docker-compose build

# Docker Compose运行（完整版）
docker-compose-up:
	docker-compose up -d

# Docker Compose运行（简化版）
docker-compose-simple:
	docker-compose -f docker-compose.simple.yml up -d

# Docker Compose停止
docker-compose-down:
	docker-compose down

# Docker Compose停止（简化版）
docker-compose-simple-down:
	docker-compose -f docker-compose.simple.yml down

# Docker Compose查看日志
docker-compose-logs:
	docker-compose logs -f html2md

# Docker清理
docker-clean:
	docker system prune -f
	docker volume prune -f

# Docker完全清理（包括镜像）
docker-clean-all:
	docker system prune -a -f
	docker volume prune -f 