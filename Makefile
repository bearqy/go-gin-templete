APP_NAME ?= go-gin-templete
CONFIG ?= config/default.yaml

.PHONY: run test vet lint build swag docker-build

run:
	go run ./cmd/server -config $(CONFIG)

test:
	go test ./...

vet:
	go vet ./...

lint: vet

build:
	go build -o bin/$(APP_NAME) ./cmd/server

swag:
	swag init -g cmd/server/main.go -o docs

docker-build:
	docker build -t $(APP_NAME):latest .
