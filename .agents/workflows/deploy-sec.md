---
description: Автономный цикл безопасного развертывания (Local & Remote Docker + Security Audit).
---
# Security & Deployment Pipeline

## 1. DevOps Phase (Docker)
- Агент `DevOpsAgent` создает минималистичные `Dockerfile` для `backend` (multistage, alpine) и `frontend` (nginx, multistage).
- Обновляет `docker-compose.yml` для локального запуска всего стека (postgres, backend, frontend).
- Создает `docker-compose.prod.yml` для удаленного продакшен-сервера.
- Экономия токенов: Агент читает ТОЛЬКО `go.mod`, `package.json` и `vite.config.ts`. Никаких исходников.

## 2. Security Audit Phase
- Агент `SecurityAgent` проводит статический аудит:
  1. Бэкенд (`handlers.go`, `service.go`): SQLi (убедиться, что везде sqlc с биндингами), XSS, JWT-конфигурации (HttpOnly cookies или заголовки).
  2. Фронтенд: хранение токенов (не в localStorage, если можно), XSS (экранирование).
  3. API-тестирование: проверка эндпоинтов (роли).
- Использует MCP сервер `context7/query-docs` для проверки best practices.
- Экономия токенов: Читает только контроллеры/хендлеры и слой сервисов. Не лезет в утилиты и верстку.

## 3. Healing Loop
- Если найдены уязвимости, Оркестратор поднимает `DeveloperAgent` для точечного фикса и `TesterAgent` для проверки.
- Цикл повторяется, пока `SecurityAgent` не выдаст "All Secure".
