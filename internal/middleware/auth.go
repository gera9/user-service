package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gera9/user-service/config"
	"github.com/gera9/user-service/pkg/models"
	"github.com/gera9/user-service/pkg/utils"
	"github.com/go-chi/render"
)

const (
	ClaimsCtxKey CtxKey = "claims"
)

func (m *MiddlewareManager) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.Render(w, r, models.ErrUnauthorized(errors.New("authorization token is required")))
			return
		}

		claims, err := utils.ParseAndValidateToken(config.Config{}, token)
		if err != nil {
			render.Render(w, r, models.ErrUnauthorized(err))
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
