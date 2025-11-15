package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, user CreateUser) (uuid.UUID, error) {
	newUser := User{
		CreateUser: user,
		ID:         uuid.New(),
	}

	id, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, userID uuid.UUID) (*ResponseUser, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	respUser := ResponseUser{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GitUrl:     user.GitUrl,
		LikedinUrl: user.LikedinUrl,
	}
	return &respUser, err
}
