package message

import (
	"context"
	"errors"
)

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrInvalidInput    = errors.New("invalid input")
)

type Service interface {
	GetMessagesByEmail(ctx context.Context, email string) ([]*Message, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetMessagesByEmail(ctx context.Context, email string) ([]*Message, error) {
	if email == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByEmail(ctx, email)
}
