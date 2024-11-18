package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gera9/user-service/pkg/utils"
)

func TestMiddlewareManager_Auth(t *testing.T) {
	tests := []struct {
		name  string
		m     *MiddlewareManager
		want  int
		token string
	}{
		{
			name:  "Test expired token",
			m:     &MiddlewareManager{},
			want:  http.StatusUnauthorized,
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJuYW1lIiwiZW1haWwiOiJlbWFpbEBlbWFpbC5jb20iLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJzdWIiOiI4YmI0ODkyMC1iYzYzLTRjODQtOWJhYy02ZTYwY2ZkMDZmMjciLCJleHAiOjE3MzE4NzgyNzB9.q1pChFiY7Oqj_RvRnVdkzebtzFcZQLoI6zL3TtTgPSU",
		},
		{
			name:  "Test invalid token",
			m:     &MiddlewareManager{},
			want:  http.StatusUnauthorized,
			token: "invalid_token",
		},
		{
			name:  "Test manipulated token",
			m:     &MiddlewareManager{},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJuYW1lIiwiZW1haWwiOiJlbWFpbEBlbWFpbC5jb20iLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJzdWIiOiI4YmI0ODkyMC1iYzYzLTRjODQtOWJhYy02ZTYwY2ZkMDZmMjciLCJleHAiOjE3MzE4NzgyNzB9.q1pChFiY7Oqj_RvRnVdkzebtzFcZQLoI6zL3TtTgPS",
			want:  http.StatusUnauthorized,
		},
		{
			name:  "Test empty token",
			m:     &MiddlewareManager{},
			token: "",
			want:  http.StatusUnauthorized,
		},
		{
			name:  "Test token with Bearer prefix",
			m:     &MiddlewareManager{},
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJuYW1lIiwiZW1haWwiOiJlbWFpbEBlbWFpbC5jb20iLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJzdWIiOiI4YmI0ODkyMC1iYzYzLTRjODQtOWJhYy02ZTYwY2ZkMDZmMjciLCJleHAiOjE3MzE4NzgyNzB9.q1pChFiY7Oqj_RvRnVdkzebtzFcZQLoI6zL3TtTgPSU",
			want:  http.StatusUnauthorized,
		},
		{
			name:  "Test valid token",
			m:     &MiddlewareManager{},
			token: utils.CreateTestingToken("username", "email"),
			want:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MiddlewareManager{}

			rr := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			r.Header.Set("Authorization", tt.token)

			m.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rr, r)

			if got := rr.Result().StatusCode; got != tt.want {
				fmt.Printf("rr.Body.String(): %v\n", rr.Body.String())
				t.Errorf("MiddlewareManager.Auth() = %v, want %v", got, tt.want)
			}
		})
	}
}
