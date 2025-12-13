# Коляда Даниил
## Практическая работа №16

### Цели работы

1. Освоить интеграционное тестирование REST API: проверка «маршрут → хендлер → сервис → репозиторий → реальная БД»
2. Научиться поднимать изолированную тестовую среду БД в Docker
3. Научиться инициализировать схему БД (миграции/auto-migrate), сидировать тестовые данные, очищать окружение
4. Внедрить интеграционные проверки CRUD-эндпоинтов (статусы, заголовки, JSON-ответы, эффекты в БД)

---

### Эндпоинты

| Тип | Адрес |
|-|-|
| POST | my.domain/notes |
| GET | my.domain/notes  |
| GET | my.domain/notes/{id} |
| PUT | my.domain/notes/{id} |
| DELETE | my.domain/notes/{id} |

---

### Тесты

| Тест | Описание | Тип запроса | Код ответа |
|-|-|-|-|
| **Test_CreateAndGet_DockerCompose** | Создаёт заметку и затем получает её по ID | POST /notes | 201 (Создано) |
||| GET /notes/{id} | 200 (Успех) |
| **Test_Update_DockerCompose** | Обновляет существующую заметку и проверяет обновлённые данные | POST /notes | 201 (Создано) |
||| PUT /notes/{id} | 200 (Успех) |
||| GET /notes/{id} | 200 (Успех) |
| **Test_Delete_DockerCompose** | Удаляет заметку по ID и проверяет, что она удалена | POST /notes | 201 (Создано) |
||| DELETE /notes/{id} | 200 (Успех) |
||| GET /notes/{id} | 404 (Не найдено) |
| **Test_List_DockerCompose** | Создаёт несколько заметок и получает их список | POST /notes | 201 (Создано) |
||| GET /notes | 200 (Успех) |
| **Test_NotFound_DockerCompose** | Проверяет ответ на запрос несуществующей заметки по ID | GET /notes/{id} | 404 (Не найдено) |

---

### Запуск

Используем команду

```bash
make up
```

Или же пишем
```bash
docker-compose up -d
```

Результат поднятия docker контейнера
![Screenshot](./screenshots/Screenshot_1.png)

---

### Тестирование

Используем команду

```bash
make test-compose
```

Или же пишем
```bash
go test -v ./integration/... -db-dsn="postgres://test:test@localhost:54321/notes_test?sslmode=disable"
```

Результаты тестов
![Screenshot](./screenshots/Screenshot_2.png)

---

### Выводы

Освоили интеграционное тестирование REST API. Научились поднимать изолированную тестовую среду БД в Docker

Научились инициализировать схему БД и внедрять интеграционные проверки CRUD-эндпоинтов

---

### Дерево проекта
```
prc_integration_16
├── Makefile
├── README.md
├── cmd
│   └── api
│       └── main.go
├── go.mod
├── go.sum
├── integration
│   ├── notes_integration_test.go
│   └── notes_tc_integration_test.go
├── internal
│   ├── db
│   │   └── migrate.go
│   ├── httpapi
│   │   └── handlers.go
│   ├── models
│   │   └── note.go
│   ├── repo
│   │   └── postgres.go
│   └── service
│       └── service.go
└── screenshots
    ├── ...

11 directories, 14 files
```