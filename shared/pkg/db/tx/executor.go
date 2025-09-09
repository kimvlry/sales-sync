package tx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactional func(ctx context.Context, tx pgx.Tx) error

func WithTx(ctx context.Context, db *pgxpool.Pool, fn Transactional) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				fmt.Printf("rollback error: %v\n", rbErr)
			}
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	committed = true
	return nil
}
