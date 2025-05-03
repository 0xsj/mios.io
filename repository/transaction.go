// repository/transaction.go
package repository

import (
	"context"

	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/jackc/pgx/v4"
)

type contextKey string

const txContextKey contextKey = "transaction"

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type PgxTxManager struct {
	conn *pgx.Conn
}

func NewTxManager(conn *pgx.Conn) TxManager {
	return &PgxTxManager{conn: conn}
}

func (m *PgxTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txContextKey, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return errors.Wrap(rbErr, "rollback failed: "+err.Error())
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func GetTxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	return tx, ok
}
