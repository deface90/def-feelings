package rest

import (
	"context"
	"github.com/deface90/def-feelings/storage"
	"net/http"
)

type ContextKey string

type Auth struct {
	Engine storage.Engine
}

func (a *Auth) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.URL.Query()["token"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := a.Engine.GetUserByToken(r.Context(), token[0])
		if err != nil || user.Status != storage.StatusActive {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("current_user"), user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
