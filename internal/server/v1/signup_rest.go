package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/auth"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// NewSignupRequest return REST signup request
func NewSignupRequest(ctx *gin.Context) (*gk.SignupRequest, error) {
	var req gk.SignupRequest
	if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
		logger.Log.Errorf(
			"invalid request header: Content-Type %s\n",
			ctx.Request.Header.Get("Content-Type"),
		)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return nil, ErrSignupBadRequest
	}
	defer func() {
		if err := ctx.Request.Body.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	err := json.NewDecoder(ctx.Request.Body).Decode(&req)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, ErrSignupBadRequest
	}
	if req.Username == "" || req.Password == "" {
		return nil, ErrSignupBadRequest
	}
	return &req, nil
}

func (a *RestAPI) addSignup(gr *gin.RouterGroup) {
	gr.POST("/signup", func(ctx *gin.Context) {
		a.SignupHandler(ctx)
	})
}

// SignupHandler handle of rest signup request
func (a *RestAPI) SignupHandler(ctx *gin.Context) {
	logger.Log.Info("SignupHandler called")
	req, err := NewSignupRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	uuid, err := a.Signup(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	token, err := auth.BuildJWT(uuid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Writer.Header().Set("Authorization", token)
	ctx.Status(http.StatusCreated)
}
