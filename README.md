# Golang Todo App

REST API приложение на Go — учебный проект. Реализует управление пользователями, задачами и статистикой с JWT-аутентификацией и Redis-кешированием.

```bash
# Управление проектом
make help
```

## Технологический стек

| Компонент         | Технология                                    |
|-------------------|-----------------------------------------------|
| Язык              | Go 1.25+                                      |
| HTTP-фреймворк    | Стандартный `net/http` (без внешних фреймворков) |
| База данных       | PostgreSQL                                    |
| Драйвер БД        | `jackc/pgx/v5`                                |
| Кеш               | Redis (`go-redis/redis/v9`)                   |
| Аутентификация    | JWT (`golang-jwt/jwt/v5`) + bcrypt            |
| Логгер            | `go.uber.org/zap`                             |
| Конфигурация      | `kelseyhightower/envconfig`                   |
| Валидация         | `go-playground/validator/v10`                 |
| Документация API  | Swagger (`swaggo/swag`)                       |
| Миграции БД       | golang-migrate                                |
| Деплой            | Docker                                        |

---

## Архитектура

Проект следует принципам **чистой архитектуры** (Clean Architecture).
Каждая фича (`auth`, `users`, `tasks`, `statistics`, `web`) разделена на три слоя:

```
Transport (HTTP Handler)
      │   Декодирует запрос, вызывает сервис, формирует ответ
      ↓
Service (Business Logic)
      │   Валидация, оркестрация вызовов, доменная логика
      ↓
Repository (Data Access)
      └─  SQL-запросы к PostgreSQL / Redis, маппинг моделей

Domain (Core)
          Сущности, инварианты, бизнес-правила — без зависимостей
```

Ключевое отличие от обычной слоистой архитектуры — **инверсия зависимостей (DIP)**:
интерфейсы определяются не в реализующем слое, а в потребляющем.

- `TasksRepository` интерфейс живёт в пакете `tasks_service` — сервис владеет контрактом
- `TasksService` интерфейс живёт в пакете `tasks_transport_http` — транспорт владеет контрактом
- `domain` не импортирует ни один другой внутренний пакет — он полностью независим

Благодаря этому зависимости всегда направлены **внутрь**, к домену, а не наружу к инфраструктуре.

**Dependency Injection** (ручная инъекция зависимостей) реализован в `cmd/todoapp/main.go`:
```
Repository → Service → HTTP Handler
```

### Аутентификация

Аккаунты (email + пароль) хранятся отдельно от доменных пользователей:

- `todoapp.accounts` — учётные данные для входа (email, bcrypt-хеш пароля)
- `todoapp.users` — доменные сущности (ФИО, телефон, привязка к задачам)

Регистрация и логин возвращают JWT-токен (HS256). Все API-роуты, кроме `/auth/*`, защищены JWT-middleware — токен передаётся в заголовке `Authorization: Bearer <token>`.

### Redis-кеширование

Статистика задач кешируется в Redis по паттерну **cache-aside** с TTL.

Реализовано через паттерн **декоратор**: `CachedStatisticsRepository` оборачивает `PostgresStatisticsRepository`, не меняя его контракт. Если данные есть в кеше — Postgres не опрашивается.

---

## Структура проекта

```
.
├── cmd/
│   └── todoapp/
│       ├── main.go              # Точка входа: инициализация и запуск
│       └── Dockerfile
├── internal/
│   ├── core/                    # Общие компоненты, не зависящие от фич
│   │   ├── auth/                # JWT: генерация и валидация токенов, конфиг
│   │   ├── config/              # Общая конфигурация приложения (часовой пояс)
│   │   ├── domain/              # Доменные сущности: Task, User, Statistics, Nullable
│   │   ├── errors/              # Sentinel-ошибки: ErrNotFound, ErrConflict, ErrUnauthorized…
│   │   ├── logger/              # Структурированный логгер (zap) + «logger in context»
│   │   └── repository/
│   │       ├── postgres/pool/   # Интерфейс пула + реализация на pgx
│   │       └── redis/pool/      # Интерфейс пула + реализация на go-redis
│   │   └── transport/http/
│   │       ├── context/         # Хранение account_id в context.Context
│   │       ├── middleware/      # CORS, RequestID, Logger, Trace, Panic, JWT
│   │       ├── request/         # Декодирование тела, path/query параметры
│   │       ├── response/        # HTTPResponseHandler, ResponseWriter, ErrorResponse
│   │       ├── server/          # HTTPServer, APIVersionRouter, Route
│   │       └── types/           # Nullable[T] с UnmarshalJSON для PATCH-запросов
│   └── features/                # Бизнес-фичи приложения
│       ├── auth/                # Регистрация и логин
│       ├── tasks/               # CRUD задач
│       ├── users/               # CRUD пользователей
│       ├── statistics/          # Статистика по задачам (с Redis-кешем)
│       └── web/                 # Отдача HTML-страниц (login, app)
├── migrations/                  # SQL-миграции (golang-migrate)
├── public/                      # Статические файлы (login.html, index.html)
├── docs/                        # Автогенерированная Swagger-документация
├── docker-compose.yaml          # Локальная инфраструктура (PostgreSQL, Redis, etc.)
└── Makefile                     # Удобные команды для разработки
```

---

## Локальный запуск

#### Предварительные требования

- Docker и Docker Compose
- Go 1.25+
- make

#### Шаги

```bash
# 1. Создать .env по примеру
cp .env.example .env

# 2. Выставить недостающие переменные окружения
code .env

# 3. Поднять окружение (PostgreSQL + Redis)
make env-up

# 4. Применить миграции БД
make migrate-up

# 5. Открыть порты сервисов окружения
make env-port-forward

# 6. Запустить приложение локально
make todoapp-run
```

После запуска:
- Страница входа: `http://127.0.0.1:5050/`
- Приложение: `http://127.0.0.1:5050/app`
- Swagger UI: `http://127.0.0.1:5050/swagger/`
- API: `http://127.0.0.1:5050/api/v1/`

## Деплой

```bash
make env-up
make migrate-up
make todoapp-deploy
```

---

## Переменные окружения

| Переменная               | Описание                                      | Пример                                              |
|--------------------------|-----------------------------------------------|-----------------------------------------------------|
| `TIME_ZONE`              | Часовой пояс (IANA)                           | `Europe/Moscow`                                     |
| `LOGGER_LEVEL`           | Уровень логирования                           | `DEBUG`                                             |
| `LOGGER_FOLDER`          | Директория для лог-файлов                     | `out/logs`                                          |
| `POSTGRES_HOST`          | Хост PostgreSQL                               | `localhost`                                         |
| `POSTGRES_PORT`          | Порт PostgreSQL                               | `5432`                                              |
| `POSTGRES_USER`          | Пользователь БД                               | `test_user`                                         |
| `POSTGRES_PASSWORD`      | Пароль БД                                     | `test_pass`                                         |
| `POSTGRES_DB`            | Имя базы данных                               | `todoapp`                                           |
| `POSTGRES_TIMEOUT`       | Таймаут запроса к БД                          | `10s`                                               |
| `REDIS_HOST`             | Хост Redis                                    | `localhost`                                         |
| `REDIS_PORT`             | Порт Redis                                    | `6379`                                              |
| `REDIS_DB`               | Номер базы Redis                              | `0`                                                 |
| `REDIS_TTL`              | TTL кеша статистики                           | `60s`                                               |
| `JWT_SECRET`             | Секрет для подписи JWT (держи в тайне)        | `change-me-in-production`                           |
| `JWT_EXPIRES`            | Время жизни JWT-токена                        | `24h`                                               |
| `HTTP_ADDR`              | Адрес и порт HTTP-сервера                     | `:5050`                                             |
| `HTTP_ALLOWED_ORIGINS`   | Разрешённые CORS origins (через запятую)      | `http://localhost:5050,null`                        |
| `SHUTDOWN_TIMEOUT`       | Таймаут graceful shutdown                     | `30s`                                               |
| `PROJECT_ROOT`           | Корень проекта для полных путей               | `/home/user/projects/todoApp-golang`                |

---

## API

Все эндпоинты, кроме `/auth/*`, требуют заголовок:
```
Authorization: Bearer <token>
```

### Аутентификация `/api/v1/auth`

| Метод  | Путь               | Описание                              | Авторизация |
|--------|--------------------|---------------------------------------|-------------|
| `POST` | `/auth/register`   | Зарегистрировать аккаунт, получить JWT | Нет        |
| `POST` | `/auth/login`      | Войти, получить JWT                   | Нет         |

### Пользователи `/api/v1/users`

| Метод    | Путь               | Описание                         |
|----------|--------------------|----------------------------------|
| `POST`   | `/users`           | Создать пользователя             |
| `GET`    | `/users`           | Список пользователей (пагинация) |
| `GET`    | `/users/{id}`      | Получить пользователя по ID      |
| `PATCH`  | `/users/{id}`      | Частично обновить пользователя   |
| `DELETE` | `/users/{id}`      | Удалить пользователя             |

### Задачи `/api/v1/tasks`

| Метод    | Путь               | Описание                                         |
|----------|--------------------|--------------------------------------------------|
| `POST`   | `/tasks`           | Создать задачу                                   |
| `GET`    | `/tasks`           | Список задач (пагинация + фильтр по `user_id`)   |
| `GET`    | `/tasks/{id}`      | Получить задачу по ID                            |
| `PATCH`  | `/tasks/{id}`      | Частично обновить задачу                         |
| `DELETE` | `/tasks/{id}`      | Удалить задачу                                   |

### Статистика `/api/v1/statistics`

| Метод  | Путь            | Описание                                                                  |
|--------|-----------------|---------------------------------------------------------------------------|
| `GET`  | `/statistics`   | Статистика задач (фильтры: `user_id`, `from`, `to` в формате YYYY-MM-DD)  |

Полная интерактивная документация доступна в **Swagger UI** по адресу `/swagger/`
