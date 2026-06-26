FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# 编译环境和运行环境的系统版本最好一致
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/server ./server
COPY --from=builder /app/config/default.yaml ./config/default.yaml
EXPOSE 8080
CMD ["./server", "-config", "./config/default.yaml"]
