package message

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository interface {
	Create(ctx context.Context, message *Message) error
	GetByEmail(ctx context.Context, email string) ([]*Message, error)
}

type repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, message *Message) error {
	_, err := r.db.NewInsert().Model(message).Exec(ctx)
	if err != nil {
		return err
	}
	// Reload to get DB-generated timestamps
	return r.db.NewSelect().Model(message).WherePK().Scan(ctx)
}

func (r *repository) GetByEmail(ctx context.Context, email string) ([]*Message, error) {
	var messages []*Message
	err := r.db.NewSelect().
		Model(&messages).
		Where("email = ?", email).
		Order("created_at DESC").
		Scan(ctx)
	return messages, err
}
