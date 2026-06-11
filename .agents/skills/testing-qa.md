---
description: Применяется ВСЕГДА при создании новых сервисов, хендлеров или UI-компонентов
---
# TDD & Testing Rules

1. **Backend (Go)**:
   - Пиши табличные тесты (Table-Driven Tests) для каждого метода в `service.go`.
   - Используй `go.uber.org/mock` (бывший gomock) для мока интерфейса `repo.Querier`.
   - Покрытие бизнес-логики (domain/service) должно быть не ниже 85%.

2. **Frontend (React)**:
   - Тестируй сложную логику (Zustand сторы, утилиты) через `vitest`.
   - Компоненты тестируй через `@testing-library/react`. Проверяй рендер и базовые user events (клик, ввод текста).
