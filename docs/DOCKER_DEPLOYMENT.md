# HTML2Markdown Docker部署指南

本文档提供HTML2Markdown服务的完整Docker部署指南。

## 📋 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [部署方式](#部署方式)
- [配置说明](#配置说明)
- [监控运维](#监控运维)
- [故障排查](#故障排查)

## 🔧 前置要求

### 系统要求

- Docker 20.10+
- Docker Compose 2.0+
- 可用内存: 至少512MB
- 可用磁盘: 至少1GB

### 端口要求

| 端口 | 服务 | 说明 |
|------|------|------|
| 8080 | HTTP API | REST API接口 |
| 9090 | GRPC API | GRPC服务接口 |
| 80 | Nginx | 反向代理（完整版） |
| 6379 | Redis | 缓存服务（可选） |
| 9091 | Prometheus | 监控服务（可选） |

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd Html2Md
```

### 2. 选择部署方式

#### 方式A: 单容器部署（开发推荐）

```bash
# 构建镜像
make docker

# 运行容器
make docker-run

# 验证服务
curl http://localhost:8080/api/v1/health
```

#### 方式B: Docker Compose简化部署（生产推荐）

```bash
# 启动服务
make docker-compose-simple

# 验证服务
curl http://localhost:8080/api/v1/health

# 停止服务
make docker-compose-simple-down
```

#### 方式C: Docker Compose完整部署（企业级）

```bash
# 启动所有服务
make docker-compose-up

# 查看服务状态
docker-compose ps

# 停止所有服务
make docker-compose-down
```

## 🏗 部署方式

### 单容器部署

最简单的部署方式，适合开发和小规模部署。

**优点:**
- 部署简单
- 资源占用少
- 易于调试

**缺点:**
- 没有反向代理
- 没有负载均衡
- 缺少监控

**使用场景:**
- 开发环境
- 个人项目
- 小规模应用

### Docker Compose简化部署

包含核心服务的生产部署方式。

**包含服务:**
- HTML2MD主服务

**优点:**
- 配置简单
- 生产可用
- 自动重启
- 健康检查

**使用场景:**
- 生产环境
- 中小型项目
- 快速部署

### Docker Compose完整部署

包含完整基础设施的企业级部署。

**包含服务:**
- HTML2MD主服务
- Nginx反向代理
- Redis缓存
- Prometheus监控

**优点:**
- 企业级特性
- 完整监控
- 高可用性
- 性能优化

**使用场景:**
- 企业环境
- 大型项目
- 高并发场景

## ⚙️ 配置说明

### 环境变量配置

复制并编辑配置文件：

```bash
cp env.example .env
vim .env
```

### 主要配置项

#### 服务配置

```bash
# 服务端口
HTTP_PORT=8080
GRPC_PORT=9090

# 运行环境
ENVIRONMENT=production

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
```

#### 转换器配置

```bash
# 输入限制
CONVERTER_MAX_INPUT_SIZE=10485760  # 10MB
CONVERTER_MAX_BATCH_SIZE=100

# 插件配置
CONVERTER_DEFAULT_PLUGINS=base,commonmark

# 超时设置
CONVERTER_TIMEOUT=30s
```

#### 资源限制

在docker-compose.yml中配置：

```yaml
deploy:
  resources:
    limits:
      memory: 512M
      cpus: '0.5'
    reservations:
      memory: 256M
      cpus: '0.25'
```

### Nginx配置

完整部署包含预配置的Nginx反向代理：

- **负载均衡**: 自动分发请求
- **静态资源缓存**: 优化性能
- **GZIP压缩**: 减少传输大小
- **安全头**: 增强安全性

配置文件位置: `nginx/conf.d/html2md.conf`

### Redis配置

可选的缓存服务，提高性能：

```bash
# Redis配置
REDIS_PASSWORD=html2md123
REDIS_HOST=redis
REDIS_PORT=6379
```

## 📊 监控运维

### 健康检查

所有部署方式都包含内置健康检查：

```bash
# 检查服务健康状态
curl http://localhost:8080/api/v1/health

# Docker容器健康状态
docker ps --format "table {{.Names}}\t{{.Status}}"
```

### 日志查看

```bash
# 单容器日志
docker logs html2md-test -f

# Docker Compose日志
make docker-compose-logs

# 特定服务日志
docker-compose logs html2md -f
```

### 性能监控

完整部署提供Prometheus监控：

- **访问地址**: http://localhost:9091
- **指标收集**: 自动采集服务指标
- **告警规则**: 可配置告警

### 服务管理

```bash
# 重启服务
docker-compose restart html2md

# 更新服务
docker-compose up -d --build

# 扩展服务
docker-compose up -d --scale html2md=3
```

## 🔍 故障排查

### 常见问题

#### 1. 容器启动失败

**症状**: 容器无法启动或立即退出

**排查步骤**:
```bash
# 查看容器日志
docker logs <container-name>

# 检查端口占用
netstat -tlnp | grep :8080

# 验证镜像
docker images | grep html2md
```

**解决方案**:
- 检查端口是否被占用
- 验证环境变量配置
- 确认Docker版本兼容性

#### 2. API请求失败

**症状**: 接口返回错误或超时

**排查步骤**:
```bash
# 检查服务状态
curl http://localhost:8080/api/v1/health

# 查看容器状态
docker ps

# 检查网络连接
docker network ls
```

**解决方案**:
- 验证端口映射
- 检查防火墙设置
- 确认容器网络配置

#### 3. 性能问题

**症状**: 响应缓慢或资源占用过高

**排查步骤**:
```bash
# 查看资源使用
docker stats

# 检查日志
docker-compose logs html2md --tail=100

# 监控指标
curl http://localhost:9091/metrics
```

**解决方案**:
- 调整资源限制
- 优化转换参数
- 启用Redis缓存

#### 4. 内存不足

**症状**: 容器被OOM killer终止

**排查步骤**:
```bash
# 查看内存使用
docker stats --no-stream

# 检查系统日志
dmesg | grep -i "killed process"
```

**解决方案**:
- 增加内存限制
- 减少并发处理
- 优化输入大小限制

### 调试技巧

#### 进入容器调试

```bash
# 进入运行中的容器
docker exec -it html2md-service sh

# 查看进程状态
ps aux

# 检查配置
env | grep -E "(HTTP|CONVERTER)"
```

#### 网络调试

```bash
# 测试容器间连接
docker exec html2md-service curl http://nginx/api/v1/health

# 查看网络配置
docker network inspect html2md_html2md-network
```

#### 性能调试

```bash
# 查看详细资源使用
docker exec html2md-service top

# 分析请求性能
time curl -X POST http://localhost:8080/api/v1/convert \
  -H "Content-Type: application/json" \
  -d '{"html": "<h1>Test</h1>", "plugins": ["base"]}'
```

## 📚 更多资源

- [API文档](../README.md#-api-接口)
- [配置参考](../env.example)
- [开发指南](../README.md#-开发)
- [项目主页](../README.md)

## 🆘 获取帮助

如果遇到问题，请：

1. 查看本文档的故障排查部分
2. 检查项目的Issue页面
3. 提交新的Issue并提供：
   - 错误信息
   - 日志输出
   - 环境信息
   - 复现步骤 