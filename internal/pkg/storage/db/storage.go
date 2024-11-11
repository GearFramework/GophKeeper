package db

import (
	"context"
	"sync"
)

// Storage struct of data store
type Storage struct {
	sync.RWMutex
	conn *StorageConnection
}

// NewStorage return new data store
func NewStorage(connectionDSN string) *Storage {
	return &Storage{
		conn: NewConnection(&ConnectionConfig{
			ConnectionDSN:   connectionDSN,
			ConnectMaxOpens: 10,
		}),
	}
}

// Init initialize storage
func (s *Storage) Init() error {
	if err := s.conn.Open(); err != nil {
		return err
	}
	_, err := s.conn.DB.ExecContext(context.Background(), `
		CREATE SCHEMA IF NOT EXISTS gks
	`)
	if err != nil {
		return err
	}
	_, err = s.conn.DB.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS gks.users (
			uuid 	    VARCHAR(255) NOT NULL PRIMARY KEY,
			username 	VARCHAR(128) UNIQUE,
			password 	VARCHAR(128) NOT NULL,
			CONSTRAINT uuid_idx UNIQUE (uuid)
		)	
	`)
	if err != nil {
		return err
	}
	_, err = s.conn.DB.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS gks.entities (
			guid 		VARCHAR(128) PRIMARY KEY,
			uuid 		VARCHAR(128) REFERENCES gks.users(uuid)
						ON DELETE RESTRICT 
						ON UPDATE RESTRICT,
			name		VARCHAR(256) NOT NULL,
			description VARCHAR(1024) NOT NULL,
			type		VARCHAR(16) NOT NULL,
			uploaded_at	pg_catalog.timestamptz DEFAULT now(),
			attr		JSONB,
			CONSTRAINT guid_idx UNIQUE (guid)
		)
	`)
	return err
}

// Up connect storage
func (s *Storage) Up() error {
	return s.Ping()
}

// Down shutdown storage
func (s *Storage) Down() {
	s.Close()
}

// Close connect to storage
func (s *Storage) Close() {
	s.conn.Close()
}

// Ping pong connection
func (s *Storage) Ping() error {
	return s.conn.Ping()
}
