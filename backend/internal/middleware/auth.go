package middleware

import (
	"context"
	"net/http"
	"strings"

	database "github.com/XackuH-ORG/go-react-e-market/backend/internal/adapters/postgresql/sqlc"
)

// Определяем кастомный тип для ключей контекста, чтобы избежать коллизий
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

type AuthMiddleware struct {
	db *database.Queries
}

func NewAuthMiddleware(db *database.Queries) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

// RequireAuth проверяет Bearer токен и кладет данные юзера в контекст
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует заголовок Authorization", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Ищем сессию в БД
		session, err := m.db.GetSessionByToken(r.Context(), token)
		if err != nil {
			// Если сессия не найдена или просрочена
			http.Error(w, "Недействительный или просроченный токен", http.StatusUnauthorized)
			return
		}

		// Кладем user_id и role в контекст запроса
		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, session.Role)

		// Передаем управление следующему хэндлеру с новым контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin пускает только администраторов (должен использоваться ПОСЛЕ RequireAuth)
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(database.UserRole)
		if !ok || role != database.UserRoleADMIN {
			http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
