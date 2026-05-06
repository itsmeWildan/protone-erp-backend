package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// contextKey untuk menyimpan transaction dalam context
type contextKey string

const txKey contextKey = "tx"

// TxManager adalah abstraksi untuk mengelola database transaction.
// Diinject ke usecase agar business logic tidak tahu tentang implementasi DB.
type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithTx menjalankan fn dalam satu transaction.
// Kalau fn return error, transaction di-rollback otomatis.
// Kalau fn berhasil, transaction di-commit.
func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Simpan tx ke context agar repository bisa menggunakannya
	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// ExtractTx mengambil transaction dari context.
// Jika tidak ada tx, gunakan pool langsung (auto-commit).
func ExtractTx(ctx context.Context, pool *pgxpool.Pool) pgxDBConn {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return pool
}

// pgxDBConn adalah interface untuk pgx.Tx dan pgxpool.Pool (keduanya implement ini)
type pgxDBConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
