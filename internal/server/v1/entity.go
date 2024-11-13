package v1

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/entity"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
)

// ListEntities return user entities
func (a *RestAPI) ListEntities(ctx context.Context, usr *model.User) (*model.Entities, error) {
	entities, err := entity.Repo.GetByUUID(ctx, usr.UUID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, ErrInternalServerError
	}
	return entities, nil
}

// CreateEntity make meta-data of user entity
func (a *RestAPI) CreateEntity(ctx context.Context, usr *model.User, req *gk.UploadEntityRequest) (*model.Entity, error) {
	if !req.Type.Is() {
		return nil, ErrUnsupportedEntityType
	}
	e := model.Entity{
		UUID:        usr.UUID,
		GUID:        auth.CreateUUID(req.Name),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		UploadedAt:  time.Now(),
		Attr:        req.Attr,
	}
	if err := entity.Repo.Create(ctx, e); err != nil {
		logger.Log.Error(err.Error())
		return nil, ErrInternalServerError
	}
	return &e, nil
}

// UploadEntity receive user binary file and store to server space
func (a *RestAPI) UploadEntity(ctx context.Context, usr *model.User, guid string, src io.Reader) error {
	ent, err := entity.Repo.GetByGUID(ctx, usr.UUID, guid)
	if err != nil {
		return ErrEntityNotFound
	}
	if ent.Type != model.BinaryData {
		return ErrEntityInvalidType
	}
	dst := ent.GetFilename(a.conf.MountPath)
	destFile, err := os.Create(dst)
	if err != nil {
		logger.Log.Error(err.Error())
		return ErrInternalServerError
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	if _, err := io.Copy(destFile, src); err != nil {
		logger.Log.Error(err)
		return ErrInternalServerError
	}
	return nil
}

// GetEntity return meta-data of user entity
func (a *RestAPI) GetEntity(ctx context.Context, usr *model.User, guid string) (*model.Entity, error) {
	ent, err := entity.Repo.GetByGUID(ctx, usr.UUID, guid)
	fmt.Println(ent)
	if err != nil {
		logger.Log.Warn(err.Error())
		return nil, ErrEntityNotFound
	}
	return ent, nil
}

// DownloadEntity return mime-type and entity file reader
func (a *RestAPI) DownloadEntity(ctx context.Context, usr *model.User, guid string) (*model.Entity, io.Reader, error) {
	ent, err := entity.Repo.GetByGUID(ctx, usr.UUID, guid)
	if err != nil {
		return nil, nil, ErrEntityNotFound
	}
	if ent.Type != model.BinaryData {
		return nil, nil, ErrEntityInvalidType
	}
	f, err := os.Open(ent.GetFilename(a.conf.MountPath))
	if err != nil {
		return nil, nil, ErrEntityNotFound
	}
	return ent, f, nil
}

// DeleteEntity deleting user entity
func (a *RestAPI) DeleteEntity(ctx context.Context, usr *model.User, guid string) error {
	ent, err := entity.Repo.GetByGUID(ctx, usr.UUID, guid)
	if err != nil {
		return ErrEntityNotFound
	}
	tx, err := entity.Repo.Store.Begin(ctx)
	if err != nil {
		logger.Log.Errorf("store transaction error: %s", err.Error())
		return ErrInternalServerError
	}
	if err = entity.Repo.Delete(ctx, usr.UUID, ent); err != nil {
		return ErrEntityDeleted
	}
	if ent.Type == model.BinaryData {
		filepath := ent.GetFilename(a.conf.MountPath)
		if err := os.Remove(filepath); err != nil {
			logger.Log.Error(err.Error())
			if err := tx.Rollback(); err != nil {
				logger.Log.Error(err.Error())
			}
			return ErrInternalServerError
		}
	}
	if err := tx.Commit(); err != nil {
		logger.Log.Error(err.Error())
		return ErrInternalServerError
	}
	return nil
}
