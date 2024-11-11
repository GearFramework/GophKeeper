package v1

import (
	"context"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/user"
)

// Signup register new user
func (a *RestAPI) Signup(ctx context.Context, req *gk.SignupRequest) (string, error) {
	_, err := user.Repo.GetByUsername(ctx, req.Username)
	if err == nil {
		return "", ErrUserAlreadyExists
	}
	hash, err := auth.CreateHashPassword(req.Password)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", ErrInternalServerError
	}
	uuid := auth.CreateUUID(req.Username)
	err = user.Repo.Create(ctx, uuid, req.Username, hash)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", ErrInternalServerError
	}
	return uuid, nil
}
