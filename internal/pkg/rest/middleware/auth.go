package middleware

import (
	"errors"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Auth rest authorization middleware
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Log.Info("middleware auth start")
		UUID, err := getUUID(ctx)
		if err == nil {
			ctx.Set(gk.AuthParamName, UUID)
		}
		ctx.Next()
	}
}

func getUUID(ctx *gin.Context) (string, error) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		logger.Log.Warn("no authorization")
		return "", errors.New("no authorization")
	}
	UUID, err := auth.GetUserUUIDFromJWT(token)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", errors.New("no authorization")
	}
	return UUID, nil
}
