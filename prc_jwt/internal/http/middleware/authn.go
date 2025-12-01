package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"example.com/prc_jwt/internal/platform/jwt"
)

type ctxKey int

const CtxClaimsKey ctxKey = iota

func AuthN(v jwt.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Parse(raw)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claimsMap := make(map[string]any)
			for key, value := range claims {
				claimsMap[key] = value
			}

			log.Printf("AuthN: token valid, user ID: %v, role: %v", claimsMap["sub"], claimsMap["role"])

			ctx := context.WithValue(r.Context(), CtxClaimsKey, claimsMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
