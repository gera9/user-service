package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

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

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := utils.ParseAndValidateToken(config.Config{}, token)
		if err != nil {
			render.Render(w, r, models.ErrUnauthorized(err))
			return
		}

		w.Header().Set("Authorization", token)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		ctx := context.WithValue(r.Context(), ClaimsCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
