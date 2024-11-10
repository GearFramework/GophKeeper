package v1

import (
	"context"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/user"
)

// Signin authenticate user
func (a *RestAPI) Signin(ctx context.Context, req *gk.SigninRequest) (string, error) {
	usr, err := user.Repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", ErrUnauthorized
	}
	if !auth.CompareHashPassword(usr.Password, req.Password) {
		return "", ErrUnauthorized
	}
	logger.Log.Infof("authenticated user: %s, UUID: %s", usr.Username, usr.UUID)
	return usr.UUID, nil
}
