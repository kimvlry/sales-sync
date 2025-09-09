package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kimvlry/sales-sync/shared/pkg/db/tx"
	"time"
	"user-service/internal/models"
)

func (r *UserRepository) CreateMarketplaceAccount(ctx context.Context, account models.MarketplaceAccount, userID string) error {
	account.ID = uuid.NewString()
	now := time.Now()

	return tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx,
			`INSERT INTO marketplace_accounts (id, user_id, marketplace_type, account_id, credentials, created_at)
             VALUES ($1, $2, $3, $4, $5, $6)`,
			account.ID, userID, account.MarketplaceType, account.AccountID, account.Credentials, now,
		)
		return err
	})
}

func (r *UserRepository) UpdateMarketplaceAccount(ctx context.Context, account models.MarketplaceAccount) error {
	return tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx,
			`UPDATE marketplace_accounts 
             SET marketplace_type=$1, account_id=$2, credentials=$3
             WHERE id=$4`,
			account.MarketplaceType, account.AccountID, account.Credentials, account.ID,
		)
		return err
	})
}

func (r *UserRepository) DeleteMarketplaceAccount(ctx context.Context, id string) error {
	return tx.WithTx(ctx, r.pool, func(ctx context.Context, t pgx.Tx) error {
		_, err := t.Exec(ctx, `DELETE FROM marketplace_accounts WHERE id=$1`, id)
		return err
	})
}

func (r *UserRepository) GetMarketplaceAccount(ctx context.Context, id string) (*models.MarketplaceAccount, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, marketplace_type, account_id, credentials 
         FROM marketplace_accounts WHERE id=$1`, id)

	var a models.MarketplaceAccount
	var mt string
	if err := row.Scan(&a.ID, &a.UserID, &mt, &a.AccountID, &a.Credentials); err != nil {
		return nil, err
	}
	a.MarketplaceType = models.MarketplaceType(mt)
	return &a, nil
}

func (r *UserRepository) GetMarketplaceAccounts(ctx context.Context, userID string) ([]models.MarketplaceAccount, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, marketplace_type, account_id, credentials 
         FROM marketplace_accounts WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.MarketplaceAccount
	for rows.Next() {
		var a models.MarketplaceAccount
		var mt string
		if err := rows.Scan(&a.ID, &mt, &a.AccountID, &a.Credentials); err != nil {
			return nil, err
		}
		a.UserID = userID
		a.MarketplaceType = models.MarketplaceType(mt)
		accounts = append(accounts, a)
	}
	return accounts, nil
}
