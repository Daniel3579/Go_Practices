# Коляда Даниил

Изучили структуру Go-проекта, создали минимальный скелет приложения и запустили его

Запустить можно командой:
```bash
go run ./cmd/myapp
```
Или же запустив само приложение `myapp.app`

---

- Конфигурационные файлы такие как `config.yaml` или `config.json`, следует разместить в папке `configs/`.

- Вся документация проекта, включая инструкции по установке и использованию, должна находиться в папке `docs/`.

- Для интеграционных тестов и тестирования функциональности нужно создать папку `test/`.

- Скрипты сборки и автоматизации такие как Makefile, следует разместить в папке `scripts/`.

- Все контракты API должны находиться в папке `api/`. Это позволит четко отделить спецификации от основной логики приложения и упростит процесс их обновления. Хранение API-контрактов в отдельной папке также способствует лучшему пониманию взаимодействия между компонентами системы.

---

Дерево проекта:
```
myapp
├── README.md
├── bin
│   └── myapp.app
├── cmd
│   └── myapp
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   └── app
│       ├── app.go
│       └── handlers
│           └── ping.go
├── screenshots
│   ├── Screenshot Logs.png
│   └── Screenshot Requests.png
└── utils
    ├── httpjson.go
    └── logger.go
```

Скриншот работы программы
![Screenshot](https://github.com/Daniel3579/Go_Practices/blob/main/myapp/screenshots/Screenshot%20Requests.png)
Скриншот отображение логов
![Logs](https://github.com/Daniel3579/Go_Practices/blob/main/myapp/screenshots/Screenshot%20Logs.png)
