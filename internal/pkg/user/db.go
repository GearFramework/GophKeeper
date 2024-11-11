package user

import (
	"context"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

var (
	sqlGetByUUID = `
        SELECT *
          FROM gks.users
         WHERE uuid = $1
    `
	sqlGetByUsername = `
        SELECT *
          FROM gks.users
         WHERE username = $1
    `
	sqlInsertUser = `
		INSERT INTO gks.users 
		       (uuid, username, password)
		VALUES ($1, $2, $3)
	`
)

// Repository of user
type Repository struct {
	Store gk.DBStorable
}

// Repo global user repository
var Repo *Repository

// NewRepository return user repository
func NewRepository(store gk.DBStorable) *Repository {
	Repo = &Repository{
		Store: store,
	}
	return Repo
}

// GetByUUID return user by UUID if exists
func (repo *Repository) GetByUUID(ctx context.Context, UUID string) (*model.User, error) {
	var user model.User
	err := repo.Store.Get(ctx, &user, sqlGetByUUID, UUID)
	return &user, err
}

// GetByUsername return user by username if exists
func (repo *Repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := repo.Store.Get(ctx, &user, sqlGetByUsername, username)
	return &user, err
}

// Create insert user to store
func (repo *Repository) Create(ctx context.Context, uuid, username, hash string) error {
	_, err := repo.Store.Insert(ctx, sqlInsertUser, uuid, username, hash)
	if err != nil {
		logger.Log.Warn(err.Error())
		return err
	}
	logger.Log.Infof("new customer UUID: %s", uuid)
	return nil
}
