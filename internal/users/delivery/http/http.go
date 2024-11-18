package http

import (
	"errors"
	"net/http"

	"github.com/gera9/user-service/internal/middleware"
	"github.com/gera9/user-service/internal/users"
	"github.com/gera9/user-service/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer = otel.Tracer("user-service")
)

type usersHandler struct {
	usersService users.UsersService
	mm           *middleware.MiddlewareManager
}

func NewUsersHandler(usersService users.UsersService, mm *middleware.MiddlewareManager) *usersHandler {
	return &usersHandler{
		usersService: usersService,
		mm:           mm,
	}
}

func (h *usersHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/signup", h.Signup)
	r.Post("/login", h.Login)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(h.mm.Auth, h.mm.VerifyUUID)

		r.Get("/", h.GetById)
		r.Patch("/", h.UpdateById)
		r.Delete("/", h.DeleteById)
	})

	return r
}

func (h *usersHandler) Signup(w http.ResponseWriter, r *http.Request) {
	userPayload := models.UserPayload{}
	err := render.Bind(r, &userPayload)
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(errors.New("invalid body")))
		return
	}

	id, err := h.usersService.Register(r.Context(), userPayload)
	if err != nil {
		render.Render(w, r, models.ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"id": id.String()})
}

func (h *usersHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "Login", trace.WithAttributes(
		attribute.String("layer", "http"),
	))
	defer span.End()

	userPayload := models.UserPayload{}
	err := render.Bind(r, &userPayload)
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(errors.New("invalid body")))
		return
	}

	token, err := h.usersService.LoginByUsername(ctx, userPayload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			render.Render(w, r, models.ErrNotFound(errors.New("invalid email")))
			return
		}
		render.Render(w, r, models.ErrUnauthorized(err))
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (h *usersHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.IdCtxKey).(uuid.UUID)
	if !ok {
		render.Render(w, r, models.ErrInternalServer(errors.New("unable to get id from context")))
		return
	}

	user, err := h.usersService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			render.Render(w, r, models.ErrNotFound(errors.New("user not found")))
			return
		}
		render.Render(w, r, models.ErrInternalServer(err))
		return
	}

	resp := models.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp)
}

func (h *usersHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.IdCtxKey).(uuid.UUID)
	if !ok {
		render.Render(w, r, models.ErrInternalServer(errors.New("unable to get id from context")))
		return
	}

	userPayload := models.UserPayload{}
	err := render.Bind(r, &userPayload)
	if err != nil {
		render.Render(w, r, models.ErrBadRequest(errors.New("invalid body")))
		return
	}

	err = h.usersService.UpdateById(r.Context(), id, userPayload)
	if err != nil {
		render.Render(w, r, models.ErrInternalServer(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *usersHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.IdCtxKey).(uuid.UUID)
	if !ok {
		render.Render(w, r, models.ErrInternalServer(errors.New("unable to get id from context")))
		return
	}

	err := h.usersService.DeleteById(r.Context(), id)
	if err != nil {
		render.Render(w, r, models.ErrInternalServer(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
