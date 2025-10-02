.PHONY: build run test

# Установите имя вашего приложения
APP_NAME=prc_http

# Сборка приложения
build:
	go build -o ./bin/${APP_NAME}.app ./cmd/server

# Запуск приложения
run:
	go run ./cmd/server

# Тестирование приложения
test:
	go test -v ./internal/test