package entity

import (
	"context"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type RepoEntityMock struct {
	mock.Mock
}

func (r *RepoEntityMock) GetByUUID(ctx context.Context, userUUID string) (*model.Entities, error) {
	args := r.Called(ctx, userUUID)
	return args.Get(0).(*model.Entities), args.Error(1)
}

func (r *RepoEntityMock) GetByGUID(ctx context.Context, userUUID, entityGUID string) (*model.Entity, error) {
	args := r.Called(ctx, userUUID, entityGUID)
	return args.Get(0).(*model.Entity), args.Error(1)
}

func (r *RepoEntityMock) Create(ctx context.Context, entity model.Entity) error {
	args := r.Called(ctx, entity)
	return args.Error(0)
}

func (r *RepoEntityMock) Delete(ctx context.Context, userUUID string, entity *model.Entity) error {
	args := r.Called(ctx, userUUID, entity)
	return args.Error(0)
}

func TestUserEntities(t *testing.T) {
	uuid := auth.CreateUUID("user")
	guid := auth.CreateUUID("guid")
	m := new(RepoEntityMock)
	m.On("GetByUUID", mock.Anything, mock.Anything, mock.Anything).Return(&model.Entities{model.Entity{
		UUID: uuid,
		GUID: guid,
	}}, nil)
	et, err := m.GetByUUID(context.Background(), uuid)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*et))
	m.AssertNumberOfCalls(t, "GetByUUID", 1)
}

func TestUserEntity(t *testing.T) {
	uuid := auth.CreateUUID("user")
	guid := auth.CreateUUID("guid")
	m := new(RepoEntityMock)
	m.On("GetByGUID", mock.Anything, uuid, guid).Return(&model.Entity{
		UUID: uuid,
		GUID: guid,
	}, nil)
	et, err := m.GetByGUID(context.Background(), uuid, guid)
	assert.NoError(t, err)
	assert.Equal(t, guid, et.GUID)
	m.AssertNumberOfCalls(t, "GetByGUID", 1)
}
