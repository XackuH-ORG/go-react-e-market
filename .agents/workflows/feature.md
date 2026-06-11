---
description: Команда /feature запускает полный цикл разработки новой фичи (Fullstack)
---
Реализуй фичу согласно описанию ниже, строго следуя AGENTS.md и всем скиллам.

Порядок выполнения:
1. Субагент DB: Напиши схему таблиц и raw SQL запросы для `sqlc` в `internal/adapters/postgresql/sqlc/queries.sql`. Сгенерируй код.
2. Субагент BACKEND: Создай домен (например `internal/products`), реализуй `service.go` с бизнес-логикой и `handlers.go` для chi-роутера.
3. Субагент FRONTEND: Создай фичу в `frontend/src/features/`, напиши хуки TanStack Query, UI компоненты и добавь их на нужную страницу в `frontend/src/pages/`.
4. Субагент REVIEW: Проверь безопасность (валидация данных, проверка ролей Admin/User) и соответствие архитектуре.

Фича: {{args}}
