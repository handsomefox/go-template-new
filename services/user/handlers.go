package user

import (
	"net/http"

	"project-template/database/sqlc"
	"project-template/handler"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Bind() func(r *chi.Mux) {
	return func(r *chi.Mux) {
		r.Post("/user", handler.Create(s.HandleCreateUser))
	}
}

type CreateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type CreateUserResponse struct {
	*sqlc.Users
}

func (s *Service) HandleCreateUser(_ http.ResponseWriter, r *http.Request) (*CreateUserResponse, error) {
	req, err := handler.Decode[CreateUserRequest](r)
	if err != nil {
		return nil, err
	}
	user, err := s.Create(r.Context(), req.Name, req.Email)
	return &CreateUserResponse{Users: user}, err
}
