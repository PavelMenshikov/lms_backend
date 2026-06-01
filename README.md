# Cap Education — LMS Backend

Бэкенд образовательной платформы для онлайн-школы. Управление курсами, уроками, посещаемостью, домашними заданиями, чатами и отчётами. Шесть ролей: студент, родитель, преподаватель, ментор, модератор, администратор.

## Стек

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)
![S3](https://img.shields.io/badge/S3-AWS-569A31?logo=amazons3)
![Chi](https://img.shields.io/badge/Chi_v5-router-5C2D91)
![WebSocket](https://img.shields.io/badge/WebSocket-gorilla-333)
![Swagger](https://img.shields.io/badge/Swagger-docs-85EA2D?logo=swagger)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker)
![Goose](https://img.shields.io/badge/Goose-migrations-00ADD8)

## Архитектура

Clean Architecture (delivery → usecase → repository). Все модули изолированы:

```
cmd/
  app/          — точка входа, роутер
  notifier/     — фоновый воркер уведомлений
internal/
  domain/       — модели данных
  auth/         — регистрация, логин, JWT-like токены
  dashboard/    — дашборды (ученик, преподаватель, админ)
  learning/     — курсы, уроки, ДЗ, посещаемость
  content_admin/ — админка: CRUD курсов, пользователей
  teacher_dashboard/ — отчёты преподавателя
  groups/       — управление группами и потоками
  schedule/     — расписание
  chat/         — вебсокет-чат
  review/       — ревью преподавателей
  statistics/   — статистика учеников
  reports/      — Excel-отчёты
  profile/      — профили
  notification/ + notifier/ — уведомления (Redis Pub/Sub)
  freeze/       — заморозки
  access/       — запросы на доступ
  banner/       — баннеры
  comment/      — комментарии
  audit/        — аудит
migrations/     — Goose-миграции (34 шт.)
pkg/             — database, storage (S3), broker (Redis)
```

## Быстрый старт

```bash
# .env
cp .env.example .env

# PostgreSQL + Redis + приложение
docker compose up -d

# Миграции накатятся автоматически при старте
```

## Переменные окружения

| Переменная      | Описание                     |
|-----------------|------------------------------|
| `API_PORT`      | Порт сервера (по умолч. 8080) |
| `DB_HOST`       | Хост PostgreSQL              |
| `DB_PORT`       | Порт PostgreSQL              |
| `DB_USER`       | Пользователь БД              |
| `DB_PASSWORD`   | Пароль БД                    |
| `DB_NAME`       | Имя БД                       |
| `REDIS_ADDR`    | Адрес Redis                  |
| `S3_ENDPOINT_URL` | S3 endpoint               |
| `S3_REGION`     | S3 регион                    |
| `S3_BUCKET_NAME` | S3 бакет                   |
| `SYSTEM_SECRET`  | Секрет для системных вызовов |

## Роли

- **student** — ученик
- **parent** — родитель
- **teacher** — преподаватель
- **curator** — МЗК (менеджер по заботе о клиентах)
- **moderator** — модератор
- **admin** — полный доступ

## Аутентификация

Токен хранится в куке `auth_token` (формат: `userID:role`). Middleware проверяет роль на каждый запрос, `RoleRequiredMiddleware` ограничивает доступ по endpoint'ам.

## Документация API

Swagger UI: `/swagger/` (генерация через `swag init -g cmd/app/main.go --parseInternal`)

## Оптимизация

Проект писался на Go с прицелом на производительность — и это оправдалось.

- **Connection pool**: `MaxOpenConns=25`, `MaxIdleConns=10`, таймауты под нагрузкой
- **N+1 запросы**: исправлено 6 кейсов — рейтинги студентов/преподавателей переписаны с подзапросов на `LEFT JOIN`, статистика ментора — через pre-agg вместо LATERAL
- **Redis-кэш**: студенческий дашборд (активные курсы, %, upcoming lessons) + админские метрики, TTL 5 мин
- **GetCourseContent**: с 5 последовательных запросов до 2 параллельных через errgroup
- **Пагинация**: `LIMIT` на все list-запросы (50–500) — никаких `SELECT ... FROM users` без ограничения
- **Админские метрики**: 3 отдельных запроса (`PerformanceStats`, `HwPerformanceStats`, `AttendancePerformanceStats`) объединены в один `GetAllPerformanceStats` через общий CTE
- **Индексы**: миграция 000034 — 13 композитных индексов под частые фильтры (role+created_at, user_id+status, teacher_id и т.д.)

## Разработка

```bash
# Локальный запуск
go run ./cmd/app

# Тесты
go test ./...

# Миграции вручную
goose -dir migrations postgres "$CONN_STRING" up
```

Требуется Go 1.25+. При первом запуске Go сам скачает toolchain — достаточно чистой установки Go, отдельно ставить ничего не нужно.
