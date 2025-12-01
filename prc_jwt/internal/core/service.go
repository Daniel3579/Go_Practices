package core

import (
	"encoding/json"
	"net/http"

	"example.com/prc_jwt/internal/http/middleware"
)

type userRepo interface {
	CheckPassword(email, pass string) (User, error)
}

type jwtSigner interface {
	Sign(userID int64, email, role string) (string, error)
}

type Service struct {
	repo userRepo
	jwt  jwtSigner
}

func NewService(r userRepo, j jwtSigner) *Service { return &Service{repo: r, jwt: j} }

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		httpError(w, 400, "invalid_credentials")
		return
	}
	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}
	tok, err := s.jwt.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}
	jsonOK(w, map[string]any{"token": tok})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	claimsVal := r.Context().Value(middleware.CtxClaimsKey)
	if claimsVal == nil {
		http.Error(w, "No claims found", http.StatusUnauthorized)
		return
	}

	claims, ok := claimsVal.(map[string]any)
	if !ok {
		http.Error(w, "Invalid claims format", http.StatusInternalServerError)
		return
	}

	jsonOK(w, map[string]any{
		"id": claims["sub"], "email": claims["email"], "role": claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"users": 2, "version": "1.0"})
}

// утилиты и ключ для контекста — экспортируем из middleware
type ctxClaims struct{}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
