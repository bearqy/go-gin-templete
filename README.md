# go-gin-templete
gin脚手架

一个基于 Gin 的轻量级服务模板，内置：
- HTTP 服务与路由
- Swagger 接口文档（/wiki -> /swagger）
- 健康检查（/health）与 Prometheus 指标（/metrics）
- MySQL 数据库初始化（xorm）
- 定时任务框架（robfig/cron）
- 统一响应结构与分页/排序工具
- 结构化日志（slog）与访问日志

### 目录结构

```
config/            # 配置文件（YAML）
docs/              # Swagger 文档（生成文件与导出 YAML/JSON）
internal/
  api/
    controller/    # 控制器（业务接口）
    util/          # 接口响应模板、分页/排序等
    router.go      # 路由注册（/wiki, /swagger, /home, /api/home 等）
  cli/             # 命令行参数（-config 指定配置文件）
  config/          # 配置加载与 SwaggerInfo 初始化
  db/              # 数据库引擎初始化与建表（xorm）
  entity/          # 数据库实体定义
  home/            # 示例业务与定时任务
  job/             # 定时任务注册与启动
  logger/          # slog 日志初始化、访问日志文件
  webserver/       # Gin 封装（健康、指标、优雅关停）
main.go            # 程序入口
```

### 快速开始

1) 准备环境
- Go 1.20+（项目 go.mod 使用 go 1.23 语法，请使用较新版本）

2) 启动服务

```bash
go run main.go -config config/default.yaml
```

启动后访问：
- Swagger 文档: `http://127.0.0.1:8080/wiki`（重定向到 `/swagger/index.html`）
- 健康检查: `http://127.0.0.1:8080/health`
- 指标采集: `http://127.0.0.1:8080/metrics`
- 示例接口: `GET /home` 或 `GET /api/home`

### 配置说明（config/default.yaml）

```yaml
web:
  address: ":8080"   # 监听地址

log:
  level: "info"       # 日志级别：debug/info/warn/error
  accessLogfile: "log/access.log"   # 访问日志输出
  runtimeLogfile: "log/runtime.log" # 运行日志输出

db:
  conn_str: ""        # MySQL 连接串，可留空（留空时将跳过数据库初始化）

# 精确到秒的 Cron 表达式
job:
  # job_name:
  #   cron: "* * * * * *"
```

说明：
- `db.conn_str` 留空时，应用仍可正常启动，仅跳过数据库连接与建表。
- 配置路径可通过命令行 `-config` 覆盖，指向任意 YAML 文件。

### 命令行参数

```bash
-config string   配置文件路径（默认：config/default.yaml）
```

示例：
```bash
go run main.go -config /path/to/your.yaml
```

### 接口与路由

- `GET /home`：健康检测示例，返回统一响应结构，`result: "Hello, World!"`
- `GET /api/home`：同上（分组路由示例）
- `GET /health`：健康检查，返回纯文本 `OK`
- `GET /metrics`：Prometheus 指标（text/plain）
- `GET /wiki`：跳转到 Swagger UI（`/swagger/index.html`）

统一响应结构（部分）：

```json
{
  "code": 1,
  "msg": "",
  "result": {},
  "error": null,
  "help": "暂不提供帮助信息"
}
```

分页/排序参数（示例）：

```json
{
  "page_num": 1,
  "page_size": 10,
  "sort": { "field": "created_at", "direction": "desc" }
}
```

### 数据库

- 使用 `xorm` 连接 MySQL，启动时（当连接串不为空）会自动 `Ping` 并建表：
  - 当前包含实体：`entity.User`
- 连接串示例：

```text
user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&loc=Asia%2FShanghai
```

### 定时任务

- 基于 `robfig/cron/v3`
- 使用方法：
  1. 在代码中注册任务：`job.Register("taskName", func() { /* do work */ })`
  2. 在配置 `job:` 下添加同名任务与 Cron 表达式
  3. 应用启动后自动加载并运行

示例（代码位置可参考 `internal/job/aaa.go`）：

```go
func init() {
    job.Register("demo", func() { /* ... */ })
}
```

```yaml
job:
  demo:
    cron: "0 */5 * * * *"   # 每 5 分钟
```

### 日志

- 使用 `slog` 输出结构化运行日志，日志级别由 `log.level` 控制
- 访问日志通过 Gin 中间件写入 `log/access.log`
- 日志文件路径会在启动时自动创建目录

### Swagger 文档

- 访问：`/wiki`（重定向至 `/swagger/index.html`）
- 生成方式：本仓库已包含生成好的 `docs/` 目录与 `swagger.yaml/json`
- 如需重新生成：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go -o docs
```

注意：`config.InitSwagger()` 可设置标题、版本等信息；如需生效，请在程序启动时调用它（当前示例未调用，仅作示例保留）。

### Docker（参考）

当前仓库的 `Dockerfile` 为示例，若你的项目入口为 `main.go`，可参考如下更简化的版本：

```Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /app/server ./server
COPY --from=builder /app/config/default.yaml ./config/default.yaml
EXPOSE 8080
CMD ["./server", "-config", "./config/default.yaml"]
```

### 许可

本项目示例代码仅供参考与二次开发使用。
