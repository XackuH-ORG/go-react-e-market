# WEB Marketplace — Project Context

## Business Requirements (Функционал)
Проект должен реализовывать следующий функционал:
1. Аутентификация и авторизация по ролям (Guest, User, Admin). Шифрование паролей.
2. **Админ-панель**: добавление, удаление, редактирование, просмотр товаров, работа с фото. Управление заказами (просмотр, смена статуса, поиск по 4 последним символам заказа). Просмотр пользователей и смена их ролей. Выход из ЛК.
3. **Пользовательский ЛК**: поиск, сортировка, фильтрация товаров. Просмотр карточек (фото, цена, наименование) и детальной информации. Управление корзиной (добавление, просмотр, удаление, оформление заказа, очистка после оформления). Просмотр истории заказов. Выход из ЛК.
4. **Гостевая зона (Главная)**: просмотр карточек товаров, детальная информация, поиск, фильтрация, сортировка (без возможности заказа).

## Tech Stack
- **Backend**: Go 1.26.3, `chi` router, PostgreSQL 18 (Docker), `sqlc` (type-safe SQL). Никаких ORM!
- **Frontend**: React 19, TypeScript, Vite, Tailwind CSS v4, Zustand, TanStack Query v5.

## Backend Architecture (Vertical Slices)
- `internal/adapters/postgresql/sqlc` — только автосгенерированный код БД. Никакой бизнес-логики.
- `internal/<domain>/` (например, auth, products, orders) — изолированные бизнес-модули.
  - `handlers.go`: HTTP-обработчики (chi), парсинг JSON, вызов сервиса.
  - `service.go`: Бизнес-правила, транзакции, вызов `sqlc` методов через интерфейс `repo.Querier`.
- Строгое DI: сервисы получают доступ к БД через интерфейс, хендлеры к сервисам — через структуры.
- Транзакции БД выполняются на уровне `service.go` с передачей `pgx.Tx` в `sqlc` (паттерн WithTx).

## Frontend Architecture (FSD Lite)
- `frontend/src/app/` — Инициализация (Router, Providers, глобальные стили).
- `frontend/src/shared/` — Общие UI-компоненты, Axios, утилиты.
- `frontend/src/features/<domain>/` — Изолированные бизнес-модули (auth, catalog, cart, orders).
- `frontend/src/pages/` — Компоновка фич в страницы.

## Agents Configuration

Определены следующие специализированные субагенты и закрепленные за ними модели (LLM).

### 1. `Manager`
**Роль**: Senior Architect / Team Lead
- **Model**: `gemini-3-1-pro`
- **Описание**: Планирование задач, распределение нагрузки, утверждение архитектурных решений, ревью сложных систем.

### 2. `BackendAgent`
**Роль**: Go Backend Developer
- **Model**: `gemini-3-5-flash`
- **Описание**: Разработка бизнес-логики (SQLC, chi, DI), написание хендлеров и сервисов. Опирается на `golang-pro` и `api-design`.

### 3. `FrontendAgent`
**Роль**: React/TypeScript Developer
- **Model**: `gemini-3-5-flash`
- **Описание**: Верстка Tailwind, управление состоянием (Zustand, TanStack Query), интеграция API.

### 4. `SecurityAgent`
**Роль**: Security Auditor / SecOps
- **Model**: `claude-opus-4-6`
- **Описание**: Точечный аудит API, проверка SQLi/XSS (Фронт+Бэк). Формирует баг-репорты для Developer/Tester в цикле до полного исправления. Использует MCP `context7` (`query-docs`) для документации по уязвимостям.

### 5. `DBAgent`
**Роль**: Database Engineer
- **Model**: `gemini-3-5-flash`
- **Описание**: Проектирование схемы БД, написание raw SQL, создание миграций. Опирается на `db-migrations`.

### 6. `DevOpsAgent`
**Роль**: DevOps Engineer
- **Model**: `gemini-3-5-flash`
- **Описание**: Подготовка локального и production развертывания в Docker. Фокус только на инфраструктурных файлах (`Dockerfile`, `docker-compose.yml`, `nginx.conf`). Использует MCP `context7` (`query-docs`) для доков по CI/CD и Docker.
