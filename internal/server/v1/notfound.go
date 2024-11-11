package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler handle if route is not exists
func (a *RestAPI) NotFoundHandler(ctx *gin.Context) {
	ctx.Status(http.StatusNotFound)
}
