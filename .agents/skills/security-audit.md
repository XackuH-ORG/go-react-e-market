---
description: Аудит безопасности (SecOps)
---
# Security Audit Guidelines

- **Auth**: Пароли хэшировать `bcrypt`. Сессии — JWT (Access < 15m, Refresh).
- **Роли**: Строгая проверка RBAC (Guest, User, Admin) в HTTP middleware.
- **SQLi**: Только prepared statements (через `sqlc`). Конкатенация запрещена.
- **XSS/CORS**: Разрешать только конкретные Origins (без `*`). Отдавать CSP-заголовки.
- **Limiting**: Rate limit на `login` / `register` от brute-force.
- **Валидация**: Жесткая проверка всех входящих JSON/Path данных (через `validator/v10`) до сервисного слоя.
