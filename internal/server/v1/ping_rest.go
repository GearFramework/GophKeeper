package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *RestAPI) addPing(gr *gin.RouterGroup) {
	gr.GET("/ping", func(ctx *gin.Context) {
		a.PingHandler(ctx)
	})
}

// PingHanlder handler of ping request
func (a *RestAPI) PingHandler(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "Content-Type: text/html", []byte("Pong"))
}
