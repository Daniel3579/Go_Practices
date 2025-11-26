# Коляда Даниил
## Практическая работа №8

### Краткое описание

**MongoDB** – документно-ориентированная NoSQL база данных:
- Данные хранятся в коллекциях документов
- Каждый документ – Binary JSON объект
- Документы могут иметь разные поля и вложенные структуры

---

### Основные понятия

1. **Документ:**
- Аналог строки в таблице, но гораздо гибче
- Представлен в формате BSON (Binary JSON)

2. **Коллекция:**
- Набор документов (аналог таблицы в SQL)
- Нет строгой схемы — в одной коллекции могут лежать документы с разными полями

3. **База данных:**
- Логическое объединение коллекций
- Пример: `prc_mongo` – база для нашего практического задания, внутри коллекция `notes`

4. **Индексы:**
- Обычный (по одному полю)
- Составной (по нескольким полям)
- Уникальный (не допускает дубликатов)
- TTL (документы удаляются автоматически после заданного времени)

---

**Преимущества:**
- Простая работа с JSON-подобными данными
- Нет «жёсткой схемы» — можно быстро изменять структуру
- Масштабируется горизонтально
- Поддерживает репликацию, индексы, агрегирующие запросы

---

**Недостатки:**
- Сложнее обеспечить строгую согласованность, чем в SQL
- Транзакции появились только в новых версиях (и медленнее, чем в PostgreSQL)
- Возможен рост объёма данных из-за дублирования полей
- Ограниченная поддержка сложных связей (JOIN делается вручную)

---

## Начало работы
**Требования (Go ≥ 1.21, Docker, curl/Postman)**

### Настройка конфигурации docker и переменных окружения:

1. Перейдем в папку проекта
2. Переименуем `docker-compose.yml.example` в `docker-compose.yml`
3. Настроим файл конфигурации

<details>
    <summary>Настройка файла конфигурации</summary>

1. Откроем файл
2. Изменим имя контейнера

```
container_name: your_container_name
```

3. Установим порты

```
ports:
      - "external_port:internal_port"
```

4. Зададим имя пользователя и пароль

```
MONGO_INITDB_ROOT_USERNAME: your_root_username
MONGO_INITDB_ROOT_PASSWORD: your_root_password
```

5. Сохраним изменения

</details>

4. Переименуем `.env.example` в `.env`
5. Настроим файл окружения

<details>
    <summary>Настройка файла окружения</summary>

1. Откроем файл
2. Зададим:
- `username` — то же самое, что и `MONGO_INITDB_ROOT_USERNAME`
- `password` — то же самое, что и `MONGO_INITDB_ROOT_USERNAME`
- `address` — адрес сервера, на котором работает контейнер с MongoDB 
- `port` — то же самое, что и `external_port`

```
MONGO_URI=mongodb://username:password@address:port
```

3. Изменим `your_db_name` на имя вашей DB

```
MONGO_DB=your_db_name
```

4. Изменим `your_server_port` на порт, на котором будет работать GO сервер

```
HTTP_ADDR=your_server_port
```

5. Сохраним изменения

</details>

---

### Запуск

1. Запустим docker контейнер

```
docker compose up -d
```

2. Запустим GO сервер

```
go run ./cmd/api
```

---

> [!WARNING]
> ## Сервер не запущен
> ### Реализованные эндпоинты
> | Тип | Адрес | Заголовок | Тело |
> |-|-|-|-|
> | **POST** | my.domain:8080/api/v1/notes | Content-Type: application/json | {"title": "...", "content": "..."} |
> | **POST** | my.domain:8080/api/v1/notes | Content-Type: application/json | {"title": "...", "content": "...", "expireTime": ...} |
> | **PATCH** | my.domain:8080/api/v1/notes/{**`id`**} | Content-Type: application/json | {"content": "..."} |
> | **DELETE** | my.domain:8080/api/v1/notes/{**`id`**} |
> | **GET** | my.domain:8080/api/v1/notes/{**`id`**} |
> | **GET** | my.domain:8080/api/v1/notes\?after={**`id`**} |
> | **GET** | my.domain:8080/api/v1/notes\?q={**`search`**} |
> | **GET** | my.domain:8080/api/v1/notes\?limit={**`limit`**} |
> | **GET** | my.domain:8080/api/v1/notes\?skip={**`skip`**} |
> | **GET** | [my.domain:8080/api/v1/notes/stats](https://google.com) |
> | **GET** | [my.domain:8080/api/v1/notes](https://google.com) |

---

### Тестирование

<details>
    <summary>Результаты тестирования</summary>

```
curl -i -X POST localhost:8080/api/v1/notes \
    -H 'Content-Type: application/json' \
    -d '{
            "title": "Fiveth note",
            "content": "Hello Mongo!"
        }'
```
![Screenshot](./screenshots/Screenshot_1.png)

```
curl -i localhost:8080/api/v1/notes\?limit=2\&skip=2
```
![Screenshot](./screenshots/Screenshot_2.png)

```
curl -i localhost:8080/api/v1/notes/6926b207ed6d88dd9154e3f9
```
![Screenshot](./screenshots/Screenshot_3.png)

```
curl -i -X PATCH localhost:8080/api/v1/notes/6926b207ed6d88dd9154e3f9 \
    -H 'Content-Type: application/json' \
    -d '{"content": "Updated content"}'
```
![Screenshot](./screenshots/Screenshot_4.png)

```
curl -X DELETE localhost:8080/api/v1/notes/6926b207ed6d88dd9154e3f9
```
![Screenshot](./screenshots/Screenshot_5.png)

```
go test -v ./internal/tests
```
![Screenshot](./screenshots/Screenshot_6.png)

```
curl -i localhost:8080/api/v1/notes\?q=update
```
![Screenshot](./screenshots/Screenshot_7.png)

```
curl -i -X POST localhost:8080/api/v1/notes \
    -H 'Content-Type: application/json' \
    -d '{
            "title": "Third note",
            "content": "Hello Mongo!",
            "expireTime": 15
        }'
```
![Screenshot](./screenshots/Screenshot_8.png)

```
curl -i localhost:8080/api/v1/notes\?after=6926c63a7beea8301be9bce8
```
![Screenshot](./screenshots/Screenshot_9.png)

```
curl -i localhost:8080/api/v1/notes/stats
```
![Screenshot](./screenshots/Screenshot_10.png)

</details>

---

### Отладка и типовые ошибки
- **`connection refused`** или **`i/o timeout`** — проверьте, что контейнер с Mongo запущен и порт проброшен

- Аутентификация – используйте **`?authSource=admin`** и правильные `root/secret`
- **`duplicate key error`** — сработал уникальный индекс **`title`**. `HTTP 409` возвращаем в `create`
- **`invalid ObjectID`** — при неверном **`id`** возвращайте `404` (мы мапим ошибку на **`ErrNotFound`**)
- **`CORS`** (если будете дергать из браузера) — добавьте **`middleware CORS`** либо тестируйте через `curl/Postman`

---

### Дерево проекта
```
prc_mongo
├── .env.example
├── .vscode
│   └── launch.json
├── README.md
├── cmd
│   └── api
│       └── main.go
├── docker-compose.yml.example
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   └── mongo.go
│   ├── notes
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── repo.go
│   └── tests
│       └── repo_test.go
└── screenshots
    ├── ...

9 directories, 22 files
```