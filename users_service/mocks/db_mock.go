package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type PgxRowMock struct {
	mock.Mock
}

func (r *PgxRowMock) Scan(dest ...any) error {
	return r.Called(dest...).Error(0)
}

type DBMock struct {
	mock.Mock
}

func (m *DBMock) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.Called(ctx, sql, args).Get(0).(pgx.Row)
}

func (m *DBMock) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)

	var tx pgx.Tx
	if args.Get(0) != nil {
		tx = args.Get(0).(pgx.Tx)
	}

	return tx, args.Error(1)
}

func (m *DBMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m *DBMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}
