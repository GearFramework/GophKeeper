package gk

import (
	"context"
	"database/sql"
)

// Storable interface of data storages
type Storable interface{}

// DBStorable interface of DB storages
type DBStorable interface {
	Storable
	Begin(context.Context) (*sql.Tx, error)
	Commit(context.Context, *sql.Tx) error
	Rollback(context.Context, *sql.Tx) error
	Get(context.Context, any, string, ...any) error
	Insert(context.Context, string, ...any) (*sql.Row, error)
	Update(context.Context, string, ...any) (*sql.Row, error)
	Find(context.Context, string, ...any) (*sql.Rows, error)
	Delete(context.Context, string, ...any) error
}
