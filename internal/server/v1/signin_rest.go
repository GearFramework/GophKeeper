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

// NewSigninRequest return REST signup request
func NewSigninRequest(ctx *gin.Context) (*gk.SigninRequest, error) {
	var req gk.SigninRequest
	if !strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
		logger.Log.Errorf(
			"invalid request header: Content-Type %s\n",
			ctx.Request.Header.Get("Content-Type"),
		)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return nil, ErrSigninBadRequest
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
		return nil, ErrSigninBadRequest
	}
	return &req, nil
}

func (a *RestAPI) addSignin(gr *gin.RouterGroup) {
	gr.POST("/signin", func(ctx *gin.Context) {
		a.SigninHandler(ctx)
	})
}

// SignupHandler handle of rest signin request
func (a *RestAPI) SigninHandler(ctx *gin.Context) {
	logger.Log.Info("SigninHandler called")
	req, err := NewSigninRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	UUID, err := a.Signin(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.AbortWithStatus(err.(*ServerError).GetHTTPStatus())
		return
	}
	token, err := auth.BuildJWT(UUID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Data(http.StatusOK, "text/plain", []byte(token))
}
