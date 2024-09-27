package user

import (
	"context"

	"project-template/database/sqlc"
)

type Service struct {
	DB      sqlc.DBTX
	Queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Service {
	return &Service{
		DB:      db,
		Queries: sqlc.New(db),
	}
}

func (s *Service) Create(ctx context.Context, name, email string) (*sqlc.Users, error) {
	return s.Queries.CreateUser(ctx, &sqlc.CreateUserParams{
		Name:  name,
		Email: email,
	})
}
