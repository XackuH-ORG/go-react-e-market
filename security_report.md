# Security Audit Report

## 1. Critical RBAC Bypass Flaw
**Location:** `backend/internal/admin/handlers.go` (`AdminMiddleware`)
**Severity:** Critical 🚨

**Issue:** 
The `AdminMiddleware` falls back to reading the `X-User-Role` header from the HTTP request if the context role is missing or invalid. 
```go
role := r.Header.Get("X-User-Role")
ctxRole := r.Context().Value("user_role")
if ctxRole != nil {
	if str, ok := ctxRole.(string); ok {
		role = str
	}
}
```
Because HTTP headers can be easily spoofed by the client, any standard user (or guest) can send a request with `X-User-Role: Admin` and completely bypass the admin authorization check.

**Recommendation:** 
Remove the fallback to the `X-User-Role` header. The role must be strictly extracted from a securely validated JWT token and placed into the request context by an upstream Auth middleware.

---

## 2. JWT Storage and Architecture Flaws
**Location:** `backend/internal/auth/handlers.go` and `frontend/src/features/auth/store.ts`
**Severity:** High ⚠️

**Issue:** 
The backend `Login` handler returns the JWT token in the JSON response body (`{"token": token}`). While the React `store.ts` does not explicitly show how it stores this token, passing the token in JSON typically forces the frontend to store it in `localStorage` or `sessionStorage`. This makes the application highly vulnerable to XSS (Cross-Site Scripting) attacks, as any malicious JavaScript can steal the token.

**Recommendation:** 
Follow modern security best practices:
1. **Backend:** Change the login handler to set the JWT via an `HttpOnly`, `Secure`, and `SameSite=Strict` cookie. This prevents client-side JavaScript from accessing the token entirely.
2. **Frontend:** Rely on the browser to automatically send the cookie with API requests (using `withCredentials` in Axios/fetch) instead of manually managing token storage in Zustand or `localStorage`.

---

## 3. SQL Injection (SQLi)
**Status:** Secure 🟢
The backend stack utilizes `sqlc`. By design, `sqlc` generates fully type-safe and parameterized queries (e.g., `WHERE id = $1`). As long as the data access layer relies strictly on the generated `sqlc` methods rather than manual string concatenation, the application is protected against SQL injection attacks.

---

## 4. Cross-Site Scripting (XSS)
**Status:** Mostly Secure 🟢
React inherently protects against XSS by auto-escaping values rendered in JSX. A review of the `auth/store.ts` and standard frontend practices did not reveal the use of dangerous APIs like `dangerouslySetInnerHTML`. However, ensuring that the JWT is not exposed to JS (as mentioned in point 2) is critical to maintaining this security posture.
