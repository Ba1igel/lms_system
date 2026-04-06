package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/baigel/lms/user-service/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// loginRequest — тело запроса для POST /auth/login
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}


type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login — LT11: POST /auth/login Принимает username + password, возвращает access_token + refresh_token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Refresh — LT12: POST /auth/refresh принимает и возвращяет токены 
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		slog.Error("refresh token failed", "err", err)
		writeError(w, "failed to refresh token", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

// Logout — LT12: POST /auth/logout
// Отзывает refresh_token в Keycloak (сессия завершается)
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		slog.Error("logout failed", "err", err)
		writeError(w, "logout failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// Register — LT13: POST /auth/register
// Создаёт нового пользователя в Keycloak.
// Доступен ТОЛЬКО для пользователей с ролью "admin" (защищён middleware.RequireAdmin).
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Register(r.Context(), req.Username, req.Email, req.Password); err != nil {
		slog.Error("registration failed", "err", err)
		writeError(w, "registration failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
