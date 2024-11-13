package user

import (
	"context"
	"testing"

	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepoUserMock struct {
	mock.Mock
}

func (r *RepoUserMock) GetByUUID(ctx context.Context, UUID string) (*model.User, error) {
	args := r.Called(ctx, UUID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *RepoUserMock) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := r.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *RepoUserMock) Create(ctx context.Context, uuid, username, hash string) error {
	args := r.Called(ctx, uuid, username)
	return args.Error(0)
}

func TestUser(t *testing.T) {
	uuid := auth.CreateUUID("user")
	m := new(RepoUserMock)
	m.On("GetByUUID", mock.Anything, mock.Anything).Return(&model.User{
		UUID:     uuid,
		Username: "denis",
	}, nil)
	u, err := m.GetByUUID(context.Background(), uuid)
	assert.NoError(t, err)
	assert.Equal(t, "denis", u.Username)
	assert.Equal(t, uuid, u.UUID)
	m.AssertNumberOfCalls(t, "GetByUUID", 1)
}
