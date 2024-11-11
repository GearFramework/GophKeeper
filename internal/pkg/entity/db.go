package entity

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

var (
	// ErrEmptyUserEntities if not found entities for user
	ErrEmptyUserEntities = errors.New("empty user entities")
)

var (
	sqlGetByUUID = `
        SELECT guid,
               name,
               description,
               type,
               uploaded_at,
               attr
          FROM gks.entities
         WHERE uuid = $1
    `
	sqlGetByGUID = `
        SELECT guid,
               name,
               description,
               type,
               uploaded_at,
               attr
          FROM gks.entities
         WHERE uuid = $1 
           AND guid = $2
    `
	sqlInsertEntity = `
		INSERT INTO gks.entities 
		       (guid, uuid, name, description, type, attr)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	sqlDeleteEntity = `
		DELETE FROM gks.entities
         WHERE uuid = $1 
           AND guid = $2
    `
)

// Repository of entities
type Repository struct {
	Store gk.DBStorable
}

// Repo global entities repository
var Repo *Repository

// NewRepository return entities repository
func NewRepository(store gk.DBStorable) *Repository {
	Repo = &Repository{
		Store: store,
	}
	return Repo
}

// GetByUUID return entities by user UUID
func (repo *Repository) GetByUUID(ctx context.Context, userUUID string) (*model.Entities, error) {
	rows, err := repo.Store.Find(ctx, sqlGetByUUID, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmptyUserEntities
		}
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	var entities model.Entities
	var guid string
	var name string
	var description string
	var uploadedAt time.Time
	var typeEntity model.EntityType
	var attr model.MetaData
	for rows.Next() {
		err := rows.Scan(&guid, &name, &description, &typeEntity, &uploadedAt, &attr)
		if err != nil {
			logger.Log.Error(err.Error())
			return nil, err
		}
		//if err := json.Unmarshal([]byte(attr), &meta); err != nil {
		//	logger.Log.Error(err.Error())
		//	continue
		//}
		entities = append(entities, model.Entity{
			GUID:        guid,
			UUID:        userUUID,
			Name:        name,
			Description: description,
			Type:        typeEntity,
			UploadedAt:  uploadedAt,
			Attr:        attr,
		})
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return &entities, err
}

// GetByGUID return entity by user UUID and entity GUID
func (repo *Repository) GetByGUID(ctx context.Context, userUUID, entityGUID string) (*model.Entity, error) {
	entity := model.Entity{}
	err := repo.Store.Get(ctx, &entity, sqlGetByGUID, userUUID, entityGUID)
	return &entity, err
}

// Create insert user entity to storage
func (repo *Repository) Create(ctx context.Context, entity model.Entity) error {
	_, err := repo.Store.Insert(ctx, sqlInsertEntity,
		entity.GUID,
		entity.UUID,
		entity.Name,
		entity.Description,
		entity.Type,
		entity.Attr,
	)
	if err != nil {
		logger.Log.Warn(err.Error())
		return err
	}
	logger.Log.Infof("new entity GUID: %s", entity.GUID)
	return nil
}

// Delete user entity from storage
func (repo *Repository) Delete(ctx context.Context, userUUID string, entity *model.Entity) error {
	logger.Log.Warnf("delete entity %s by user %s", entity.GUID, userUUID)
	return repo.Store.Delete(ctx, sqlDeleteEntity, userUUID, entity.GUID)
}
