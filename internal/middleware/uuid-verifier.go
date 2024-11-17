package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gera9/user-service/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const (
	IdCtxKey CtxKey = "id"
)

func (m *MiddlewareManager) VerifyUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")

		id, err := uuid.Parse(idString)
		if err != nil {
			render.Render(w, r, models.ErrBadRequest(errors.New("invalid id")))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), IdCtxKey, id)))
	})
}
