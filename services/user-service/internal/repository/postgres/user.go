package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kimvlry/sales-sync/shared/pkg/db/tx"
	"time"
	"user-service/internal/models"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, telegramID, name string) (*models.User, error) {
	user := &models.User{
		ID:         uuid.NewString(),
		TelegramID: telegramID,
		Name:       name,
		CreatedAt:  time.Now(),
		Accounts:   []models.MarketplaceAccount{},
	}

	err := tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx,
			`INSERT INTO users (id, telegram_id, name, created_at) VALUES ($1, $2, $3, $4)`,
			user.ID, user.TelegramID, user.Name, user.CreatedAt,
		)
		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx,
			`UPDATE users SET telegram_id=$1, name=$2 WHERE id=$3`,
			user.TelegramID, user.Name, user.ID,
		)
		return err
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	return tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
		return err
	})
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{ID: id}

	err := r.pool.QueryRow(ctx,
		`SELECT telegram_id, name, created_at FROM users WHERE id=$1`,
		id,
	).Scan(&user.TelegramID, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	accounts, err := r.GetMarketplaceAccounts(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Accounts = accounts

	return user, nil
}
