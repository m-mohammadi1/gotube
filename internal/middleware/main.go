package middleware

import (
	"context"
	"gotube/internal/config"
	"gotube/internal/utils/authutil"
	"gotube/pkg/repository"
	"net/http"
)

type Middleware struct {
	config   config.Data
	userRepo repository.UserRepository
}

func New(config config.Data, userRepo repository.UserRepository) Middleware {
	return Middleware{
		config:   config,
		userRepo: userRepo,
	}
}

func (m *Middleware) EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := authutil.VerifyAuthTokenInRequestHeader(r, m.config)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.Find(r.Context(), claims.UserID)
		if err != nil || user == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
