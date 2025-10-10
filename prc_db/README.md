# Коляда Даниил

Команда создания датабазы `prc_todo`
```postgres
create database prc_todo;
```
![Screenshot](./screenshots/Screenshot_1.png)

Команда создания таблицы `tasks`
```postgres
create table if not exists tasks (
    id serial primary key,
    title text not null,
    done boolean not null default false,
    created_at timestamptz not null default now()
);
```
![Screenshot](./screenshots/Screenshot_2.png)

Команда вывода всех записей из таблицы `tasks`
```postgres
select * from tasks order by id;
```
![Screenshot](./screenshots/Screenshot_3.png)

Вывод после исполнения
```
go run .
```
![Screenshot](./screenshots/Screenshot_4.png)

### Ответы на вопросы
Пул соединений `*sql.DB` – это коллекция соединений с базой данных, которая управляет эффективным использованием ресурсов

`db.SetMaxOpenConns(10)` — устанавливает максимальное кол-во открытых соединений\
`db.SetMaxIdleConns(5)` — устанавливает максимальное кол-во неактивных соединений\
`db.SetConnMaxLifetime(30 * time.Minute)` — устанавливает время жизни соединения\

---

Плейсхолдеры (`$1`, `$2` и т.д.) используются в SQL-запросах для безопасной передачи параметров

---

| Функция |Возвращаемое значение | Описание |
|-|-|-|
| **Query** | `*Rows`, `error` | Выполняет SQL-запрос и возвращает множество строк в виде `*Rows` |
| **QueryRow** | `*Row`, `error` | Выполняет SQL-запрос и возвращает только одну строку |
| **Exec** | `Result`, `error` | Выполняет SQL-запрос, который не возвращает строк (например, `INSERT`, `UPDATE`, или `DELETE`) |

---

Транзакции в базах данных обеспечивают атомарность, целостность, изолированность и постоянство (ACID) операций