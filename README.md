# go-gin-templete

一个基于 Gin 的 Go Web 服务脚手架，目标是提供轻量、清晰、可直接二次开发的工程基础。

内置能力：
- Gin HTTP 服务、分组路由和基础中间件
- Swagger UI: `/wiki` -> `/swagger/index.html`
- 健康检查 `/health`、示例接口 `/api/v1/home`、Prometheus 指标 `/metrics`
- YAML 配置、环境变量覆盖和基础配置校验
- slog 结构化运行日志与 Gin 访问日志
- 可选 MySQL 初始化，默认使用 xorm
- robfig/cron 秒级定时任务框架
- Makefile、Dockerfile、GitHub Actions CI、基础单元测试

## 目录结构

```text
cmd/server/        # 标准服务入口
config/            # YAML 配置
docs/              # Swagger 生成文件
internal/
  app/             # 应用组装：配置、日志、DB、任务、HTTP 生命周期
  api/             # 路由、控制器、响应工具
  cli/             # 命令行参数解析
  config/          # 配置加载、默认值、环境变量覆盖
  db/              # xorm 初始化和关闭
  entity/          # 示例数据库实体
  home/            # 示例业务代码
  job/             # 定时任务注册和管理
  logger/          # slog 和日志文件生命周期
  tool/            # 通用工具函数
  webserver/       # Gin engine 和 HTTP server
```

## 快速开始

```bash
make run
```

等价命令：

```bash
go run ./cmd/server -config config/default.yaml
```

启动后访问：
- Swagger: `http://127.0.0.1:8080/wiki`
- Liveness: `http://127.0.0.1:8080/health`
- Metrics: `http://127.0.0.1:8080/metrics`
- 示例接口: `http://127.0.0.1:8080/api/v1/home`

根目录 `main.go` 仍保留兼容入口，推荐新项目和容器构建使用 `./cmd/server`。

## 配置

默认配置在 `config/default.yaml`：

```yaml
web:
  address: ":8080"
  read_timeout_second: 10
  write_timeout_second: 30
  idle_timeout_second: 60

log:
  level: "info"
  accessLogfile: "log/access.log"
  runtimeLogfile: "log/runtime.log"

db:
  conn_str: ""
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime_second: 300

job:
  # demo:
  #   cron: "0 */5 * * * *"
```

常用环境变量覆盖：
- `GIN_TEMPLATE_WEB_ADDRESS`
- `GIN_TEMPLATE_LOG_LEVEL`
- `GIN_TEMPLATE_ACCESS_LOG_FILE`
- `GIN_TEMPLATE_RUNTIME_LOG_FILE`
- `GIN_TEMPLATE_DB_CONN_STR`
- `GIN_TEMPLATE_DB_MAX_OPEN_CONNS`
- `GIN_TEMPLATE_DB_MAX_IDLE_CONNS`
- `GIN_TEMPLATE_DB_CONN_MAX_LIFETIME_SECOND`

数据库连接串为空时会跳过数据库初始化，应用仍可正常启动。

## 数据库

当前默认保留 xorm。启动时如果配置了 `db.conn_str`，应用会：
- 创建 xorm engine
- 设置连接池参数
- Ping 数据库
- 使用 `Sync2` 同步示例实体 `entity.User`
- 关闭应用时关闭数据库连接

生产项目建议按团队习惯接入迁移工具，避免长期依赖自动建表。

## 定时任务

定时任务使用秒级 cron 表达式：

```go
job.Register("demo", func() {
    // do work
})
```

```yaml
job:
  demo:
    cron: "0 */5 * * * *"
```

配置里的任务名必须已注册，否则启动会直接返回错误，避免运行时才发现 nil 任务。

## 常用命令

```bash
make run
make test
make vet
make lint
make build
make swag
make docker-build
```

## Docker

```bash
docker build -t go-gin-templete:latest .
docker run --rm -p 8080:8080 go-gin-templete:latest
```

## Swagger

仓库包含已生成的 `docs/`。重新生成：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swag
```

## CI

`.github/workflows/ci.yml` 会执行：
- `go test ./...`
- `go vet ./...`
- `go build -o /tmp/go-gin-templete-server ./cmd/server`

## License

本项目示例代码仅供参考与二次开发使用。
