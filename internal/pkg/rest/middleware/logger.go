package middleware

import (
	"time"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Logger middleware logger handle
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		logger.Log.Infof("%s request: %s",
			ctx.Request.Method,
			ctx.Request.RequestURI,
		)
		ctx.Next()
		duration := gk.GetDurationInMilliseconds(start)
		logger.Log.Infof("%s response from: %s; status: %d; size: %d | duration: %.4f ms",
			ctx.Request.Method,
			ctx.Request.RequestURI,
			ctx.Writer.Status(),
			ctx.Writer.Size(),
			duration,
		)
	}
}
