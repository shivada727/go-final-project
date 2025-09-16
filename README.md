# ToDo Scheduler (Practicum Final)

## Описание
Веб-сервер написанный на Go для списка задач. Хранит задачи в SQLite, умеет добавлять/редактировать/удалять, отмечать выполненными и считать следующую дату по правилу повторения. Работает вместе с фронтендом из папки `web`.

---

## Задания со звёздочкой
- Поддержка переменной окружения `TODO_DBFILE` — можно задать путь к файлу базы данных (если пусто, берётся `./scheduler.db`).
- Расширенные правила повторений в `NextDate()` — помимо `d N` поддерживаются:
  - `w 1,3,5` — дни недели (1=Mon … 7=Sun),
  - `m 1,15,-1 [месяцы]` — дни месяца (поддержаны `-1` — последний, `-2` — предпоследний; можно ограничить месяцами).

---

## Как запустить локально
Требуется Go 1.21+.

Переменные окружения:
- `TODO_PORT` — порт HTTP-сервера (по умолчанию `7540`)
- `TODO_DBFILE` — путь к SQLite-файлу (по умолчанию `./scheduler.db`)

Пример `.env`:

TODO_PORT=7540
TODO_DBFILE=./scheduler.db

Запуск:
```bash
go run .
# или
go build -o todo-server .
./todo-server
```

Открой в браузере: http://localhost:7540/

## Как запустить тесты

1. Запусти сервер:
```bash
go run .
```

2. В другом терминале:

```bash
#По шагам
go test -run ^TestAddTask$   ./tests
go test -run ^TestTasks$     ./tests
go test -run ^TestTask$      ./tests
go test -run ^TestEditTask$  ./tests
go test -run ^TestDone$      ./tests
go test -run ^TestDelTask$   ./tests

# все тесты
go test ./tests
```

# Настройки tests/settings.go
var BaseURL = "http://localhost:7540"
var Port    = "7540"