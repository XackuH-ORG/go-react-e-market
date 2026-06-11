---
description: Стандарты написания и ревью Go-кода (уровень Senior)
---
# Go 1.26 Backend Pro Rules

- **Роутинг**: `github.com/go-chi/chi/v5`.
- **Ошибки**: Оборачивай `fmt.Errorf("pkg.Func: %w", err)`. Пиши `log.Error("ошибка", sl.Err(err))`.
- **БД**: Только `sqlc` (чистый SQL в `queries.sql`). ORM строго запрещены.
- **Типы**: UUID -> `github.com/google/uuid`, JSONB -> `json.RawMessage`.
- **Транзакции**: Строго `pgx.Tx` (WithTx) при операциях с >1 таблицей.
- **Логирование**: Только `log/slog` с оберткой из `[@backend/internal/lib/logger/sl/sl.go]`.
- **Контекст**: `ctx context.Context` — всегда первый аргумент.
- **Архитектура**: DI через интерфейсы репозиториев в сервисах.
