# Коляда Даниил

### Цели работы:
1. Освоить базовую маршрутизацию HTTP-запросов в Go на примере роутера chi
2. Научиться строить REST-маршруты и обрабатывать методы GET/POST/PUT/DELETE
3. Реализовать небольшой CRUD-сервис «ToDo» (без БД, хранение в памяти)
4. Добавить простое middleware (логирование, CORS)
5. Научиться тестировать API запросами через curl/Postman/HTTPie

---

### Итоги работы:
- Набор тестовых запросов и ответов лежит в
[request.md](./request.md)

- Добавлена
[валидация](./internal/task/handler.go#L66-L73)
длины `title` от 3 до 100 символов

- Добавлена
[пагинация](./internal/task/handler.go#L150-L180)
и
[фильтрация](./internal/task/handler.go#L136-L148)
по `done`  
Если параметры невалидные или отсутствуют, то будут использоваться по умолчанию `page=1`, `limit=5`

- Реализовано
[чтение](./internal/task/repo.go#L77-L101)
и
[сохранение](./internal/task/repo.go#L104-L114)
в `JSON`

- Добавлено
[версионирование](./cmd/server/main.go#L34-L36)
`API`

---

### Ошибки, их описание и коды ответа

| Код | Описание | Обработка |
|-|-|-|
| 200 | OK | Запрос выполнен успешно, возвращены необходимые данные |
| 201 | Created | Запрос успешно выполнен, создана новая запись |
| 204 | No Content | Запрос выполнен успешно, но нет содержимого для возврата |
| 400 | Bad Request | Ошибка в запросе. Сервер не может его обработать |
| 404 | Not Found | Запрашиваемый ресурс не найден |

---

### Запросы и их ошибки

| Маршрут | Запрос | Ответ |
|-|-|-|
| **GET /tasks** | С параметрами или без | 204 No Content, когда нет записей для вывода |
| **GET /tasks** | Без параметров | 200 OK + все записи |
| **GET /tasks** | `?done=true` | 200 OK + записи, у которых `done=true` |
| **GET /tasks** | `?page=2&limit=10` | 200 OK + до 10 записей, начиная со второй страницы |
| **GET /tasks/{id}** | `.../task/-1` | 400 Bad Request `{error: "invalid id"}` |
| **GET /tasks/{id}** | `.../task/adc` | 400 Bad Request `{error: "invalid id"}` |
| **GET /tasks/{id}** | `.../task/1` (отсутствует) | 404 Not Found `{error: "task not found"}` |
| **GET /tasks/{id}** | `.../task/1` (присутствует) | 200 OK + запись |
| **POST /tasks** | `{"title":""}` | 400 Bad Request `{error: "invalid json: require non-empty title"}` |
| **POST /tasks** | `{"title":"a"}` | 400 Bad Request `{error: "the title length must be at least 3"}` |
| **POST /tasks** | Длина заголовка больше 100 | 400 Bad Request `{error: "the title length must be less than 100"}` |
| **POST /tasks** | `{"title":"Новая задача"}` | 201 Created + запись |
| **PUT /tasks/{id}** | `.../task/-1` + `{"title":"New", "done":true}` | 400 Bad Request `{error: "invalid id"}` |
| **PUT /tasks/{id}** | `.../task/abc` + `{"title":"New", "done":true}` | 400 Bad Request `{error: "invalid id"}` |
| **PUT /tasks/{id}** | `.../task/1` (отсутствует) + `{"title":"New", "done":true}` | 404 Not Found `{error: "task not found"}` |
| **PUT /tasks/{id}** | `.../task/1` (присутствует) + `{"title":"", "done":true}` | 400 Bad Request `{error: "invalid json: require non-empty title"}` |
| **PUT /tasks/{id}** | `.../task/1` (присутствует) + `{"title":"New", "done":true}` | 200 OK + обновленная запись |
| **DELETE /tasks/{id}** | `.../tasks/-1` | 400 Bad Request `{error: "invalid id"}` |
| **DELETE /tasks/{id}** | `.../tasks/abc` | 400 Bad Request `{error: "invalid id"}` |
| **DELETE /tasks/{id}** | `.../tasks/1` (отсутствует) | 404 Not Found `{error: "task not found"}` |
| **DELETE /tasks/{id}** | `.../tasks/1` (присутствует) | 204 No Content, запись удалена |

---

### Дерево проекта
```
prc_todo
├── README.md
├── cmd
│   └── server
│       └── main.go
├── db.json
├── go.mod
├── go.sum
├── internal
│   └── task
│       ├── handler.go
│       ├── model.go
│       └── repo.go
├── pkg
│   └── middleware
│       ├── cors.go
│       └── logger.go
├── request.md
└── screenshots
    ├── ...
```
