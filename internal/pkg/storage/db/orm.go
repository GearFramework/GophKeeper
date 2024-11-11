package db

import (
	"context"
	"database/sql"
)

// Begin transaction
func (s *Storage) Begin(ctx context.Context) (*sql.Tx, error) {
	return s.conn.DB.BeginTx(ctx, nil)
}

// Commit transaction
func (s *Storage) Commit(ctx context.Context, tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback transaction
func (s *Storage) Rollback(ctx context.Context, tx *sql.Tx) error {
	return tx.Rollback()
}

// Get data from DB-storage
func (s *Storage) Get(ctx context.Context, dest any, query string, args ...any) error {
	err := s.conn.DB.GetContext(ctx, dest, query, args...)
	return err
}

// Insert data into DB-storage
func (s *Storage) Insert(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	row := s.conn.DB.QueryRowContext(ctx, query, args...)
	return row, row.Err()
}

// Update data in DB-storage
func (s *Storage) Update(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	row := s.conn.DB.QueryRowContext(ctx, query, args...)
	return row, row.Err()
}

// Find data in DB-storage
func (s *Storage) Find(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.conn.DB.QueryContext(ctx, query, args...)
}

// Delete data from DB-storage
func (s *Storage) Delete(ctx context.Context, query string, args ...any) error {
	_, err := s.conn.DB.ExecContext(ctx, query, args...)
	return err
}
