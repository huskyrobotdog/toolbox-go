package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/huskyrobotdog/toolbox-go/constants"
	"github.com/huskyrobotdog/toolbox-go/log"
	"github.com/huskyrobotdog/toolbox-go/web"
	"go.uber.org/zap"
)

func ErrorCode(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case *web.ErrorWrap:
				if e.AttaMsg != nil {
					log.Instance.WithOptions(zap.AddCallerSkip(3)).Warn(e.Code.Error(), zap.String("attaMsg", *e.AttaMsg))
				} else {
					log.Instance.WithOptions(zap.AddCallerSkip(3)).Warn(e.Code.Error())
				}
				web.Send(c, e.Code)
			default:
				log.Instance.WithOptions(zap.AddCallerSkip(1)).Error("unknown error", zap.Any("cause", err))
				web.Send(c, constants.CodeServiceError)
			}
			c.Abort()
		}
	}()
	c.Next()
}
