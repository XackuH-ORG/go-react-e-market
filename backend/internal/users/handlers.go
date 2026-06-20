package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// GetUsers обрабатывает GET /api/v1/admin/users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Парсим query-параметры для пагинации
	limit := int32(20)
	offset := int32(0)

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = int32(parsed)
		}
	}

	users, err := h.service.ListUsers(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}

	// Если пользователей нет, sqlc может вернуть nil, отдаем пустой массив
	if users == nil {
		users = []database.GetUsersRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

// UpdateRole обрабатывает PATCH /api/v1/admin/users/{id}/role
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Достаем ID из URL
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "ID пользователя обязателен", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Передаем распарсенный UUID в сервис
	user, err := h.service.UpdateRole(r.Context(), userUUID, req.Role)
	if err != nil {
		if err.Error() == "недопустимая роль" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Ошибка при обновлении роли", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
