# Коляда Даниил

Развернули рабочее окружение Go на MacOS, создали минимальный HTTP-сервис на net/http, подключили и использовали внешнюю зависимость, собрали и проверили приложение.
- Установили Go и Git, проверили версии.
- Инициализировали модуль Go в новом проекте.
- Реализовали HTTP-сервер с маршрутами /hello (текст) и /user (JSON).
- Подключили внешнюю библиотеку (генерация UUID) и использовали её в /user.
- Запустили и проверили ответы curl/браузером.
- Собрали бинарник .app и подготовить README и отчёт.

Команда запуска:
```bash
go run ./cmd/server
```

Команда сборки под Mac:
```bash
go build -o helloapi.app ./cmd/server
```

Команда сборки под Windows:
```bash
go build -o helloapi.exe ./cmd/server
```

Команда для настройки переменной окружения в Windows:
```bash
$env:APP_PORT="8081"
```

Команда для настройки переменной окружения в MacOS:
```bash
export APP_PORT="8081"
```

Команда для запуска скомпилированной программы с кастомным портом:
```bash
APP_PORT="8081" ./helloapi.app
```

Дерево проекта:
```
helloapi
├── README.md
├── Screenshot 2025-09-09 at 8.12.51 PM.png
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
└── helloapi.app
```

Скриншот работы программы
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/helloapi/Screenshot%202025-09-09%20at%208.12.51%E2%80%AFPM.png)
