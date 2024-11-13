package db

import (
	"context"

	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// StorageConnection struct of pgsql connection
type StorageConnection struct {
	DB        *sqlx.DB
	Config    *ConnectionConfig
	pgxConfig *pgx.ConnConfig
}

// NewConnection return new connection with pgsql
func NewConnection(config *ConnectionConfig) *StorageConnection {
	return &StorageConnection{
		Config: config,
	}
}

// Open connection to pgsql
func (conn *StorageConnection) Open() error {
	var err error = nil
	if conn.pgxConfig, err = conn.getPgxConfig(); err != nil {
		return err
	}
	return conn.openSqlxViaPooler()
}

func (conn *StorageConnection) openSqlxViaPooler() error {
	db := stdlib.OpenDB(*conn.pgxConfig)
	conn.DB = sqlx.NewDb(db, "pgx")
	conn.DB.SetMaxOpenConns(conn.Config.ConnectMaxOpens)
	return nil
}

func (conn *StorageConnection) getPgxConfig() (*pgx.ConnConfig, error) {
	pgxConfig, err := pgx.ParseConfig(conn.Config.ConnectionDSN)
	if err != nil {
		logger.Log.Errorf("Unable to parse DSN: %s", err.Error())
		return nil, err
	}
	return pgxConfig, nil
}

// Ping test connection with pgsql
func (conn *StorageConnection) Ping() error {
	return conn.DB.PingContext(context.Background())
}

// Close connection with pgsql
func (conn *StorageConnection) Close() {
	if conn.Ping() == nil {
		logger.Log.Info("close storage connection")
		if err := conn.DB.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}
}
