---
description: Работа с БД и миграциями (PostgreSQL)
---
# DB & Migrations Rules

- **Инструмент**: `golang-migrate` или `goose`.
- **Формат**: Строго up/down файлы: `<version>_<name>.[up|down].sql`.
- **DDL**: Избегай удаления колонок (`DROP COLUMN`) без многоступенчатого деплоя (no breaking changes).
- **Индексы**: Создание в проде только `CONCURRENTLY` (в отдельной транзакции).
- **Консистентность**: constraints и foreign keys на уровне БД обязательны.
