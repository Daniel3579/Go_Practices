# Коляда Даниил
## Практическая работа №6

### Краткое описание
ORM (Object–Relational Mapping) — это технология сопоставления объектно-ориентированных сущностей программы с таблицами реляционной базы данных

---

### Зачем нужен ORM
1. Ускорение разработки — меньше шаблонного SQL, код компактнее
2. Безопасность — ORM автоматически использует параметризацию (? или $1) и защищает от SQL-инъекций
3. Поддержка миграций — создание/обновление схемы таблиц прямо из кода
4. Кросс-СУБД — один и тот же код работает с PostgreSQL, MySQL, SQLite

---

### Чем помог GORM
1. CRUD через GORM пишется в несколько строк кода
2. Миграции (AutoMigrate) позволяют автоматически создавать таблицы

---

> [!TIP]
> ## Сервер запущен
> ### Реализованные эндпоинты
> | Тип | Адрес | Заголовок | Тело |
> |-|-|-|-|
> | POST | my.domain:8080/users | Content-Type: application/json | {"name":"...", "email":"..."} |
> | POST | my.domain:8080/notes | Content-Type: application/json | {"title":"...", "content":"...", "userId":..., "tags":["...", "..."]} |
> | GET | my.domain:8080/notes/{id} |
> | GET | [my.domain/notes](http://79.164.43.106:2000/notes) |
> | GET | [my.domain/tags](http://79.164.43.106:2000/tags) |
> | GET | [my.domain/users](http://79.164.43.106:2000/users) |

---

### Cоздание пользователя
```
curl -X POST localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name":"Daniel", "email":"daniel@example.com"}'
```
![Screenshot](./screenshots/Screenshot_1.png)

---

### Создание заметки
```
curl -X POST localhost:8080/notes \
-H "Content-Type: application/json" \
-d '{"title":"Моя заметка", "content":"Заметка заметка какая замечательная заметка", "userId":7, "tags":["go", "gorm"]}'
```
![Screenshot](./screenshots/Screenshot_2.png)

---

### Получение заметки
```
curl localhost:8080/notes/3
```
![Screenshot](./screenshots/Screenshot_3.png)

---

### Получение всех заметок
```
curl localhost:8080/notes
```
![Screenshot](./screenshots/Screenshot_4.png)

---

### Получение всех тэгов
```
curl localhost:8080/tags
```
![Screenshot](./screenshots/Screenshot_5.png)

---

### Получение всех пользователей
```
curl localhost:8080/users
```
![Screenshot](./screenshots/Screenshot_6.png)

---

### Фрагмент схемы БД
![Screenshot](./screenshots/Screenshot_7.png)

---

### Дерево проекта
```
prc_gorm
├── README.md
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   └── postgres.go
│   ├── httpapi
│   │   ├── handlers.go
│   │   └── router.go
│   └── models
│       └── models.go
└── screenshots
    ├── ...
```
