# Коляда Даниил
## Практическая 4

Конфигурация порта была вынесена в переменную окружения [PRC_HTTP_PORT](./cmd/server/main.go#L42)

---

Запустить программу можно так:
```
PRC_HTTP_PORT=3000 make run
```

Или так:
```
export PRC_HTTP_PORT=3000
make run
```

---

- Был создан
[Makefile](./Makefile)
с целями run, buid, test

- Были созданы
[юнит-тесты](./internal/test/handlers_test.go)
для обработчиков

- Набор тестовых запросов лежит в [request.md](./requests.md)

- Были добавлены заголовки CORS в
[отдельный](./internal/api/cors.go)
middleware

- [Добавлена](./internal/api/handlers.go#L76-L83)
валидация длины `title`

- Добавлен метод
[PATCH](./internal/api/handlers.go#L106-L119)
для поля `Done`

- Добавлен метод [DELETE](./internal/api/handlers.go#L122-L135)

- Был
[сделан](./cmd/server/main.go#L55-L81)
Graceful shutdown через http.Server и контекст

---

### Дерево проекта
```
prc_http
├── Readme.md
├── Makefile
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── internal
│   ├── api
│   │   ├── cors.go
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── responses.go
│   ├── storage
│   │   └── memory.go
│   └── test
│       └── handlers_test.go
├── requests.md
└── screenshots
    ├── ...
```
