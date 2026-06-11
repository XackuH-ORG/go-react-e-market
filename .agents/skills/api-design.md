---
description: Проектирование REST API
---
# REST API Design Standards

- **Методы**: GET (чтение), POST (создание), PUT/PATCH (мутация), DELETE (удаление).
- **Endpoint**: Существительные во множественном числе (напр. `/api/v1/orders/{id}`).
- **Ответы**: 
  - Успех: `{"data": {...}}` или `{"data": [...], "total": 100}` (для пагинации).
  - Ошибка: `{"error": "string_code", "message": "human text"}`.
- **HTTP Статусы**: 200 (OK), 201 (Created), 400 (Bad Req), 401 (Unauth), 403 (Forbidden), 404 (Not Found), 500 (Internal).
