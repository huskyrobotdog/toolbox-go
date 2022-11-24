package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huskyrobotdog/toolbox-go/constants"
	"github.com/huskyrobotdog/toolbox-go/id"
	"github.com/huskyrobotdog/toolbox-go/log"
	"github.com/huskyrobotdog/toolbox-go/token"
	"github.com/huskyrobotdog/toolbox-go/web"
	"go.uber.org/zap"
)

func Authentication(tokenRequired bool) func(*gin.Context) {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader(constants.WebUserToken)

		if tokenStr == "" {
			if tokenRequired {
				web.Send(ctx, constants.CodeNotLogin)
				ctx.Abort()
				return
			} else {
				ctx.Next()
				return
			}
		}

		claims, err := token.Instance.ParseToken(tokenStr)
		if err != nil {
			log.Instance.Warn("token invalid", zap.Error(err))
			web.Send(ctx, constants.CodeTokenExpired)
			ctx.Abort()
			return
		}

		loginId, err := id.Parse(claims.ID)
		if err != nil {
			log.Instance.Warn("convert id failed", zap.Error(err))
			web.Send(ctx, constants.CodeServiceError)
			ctx.Abort()
			return
		} else {
			web.RecordLoginID(ctx, loginId)
		}

		if persistentToken, err := web.LoginIDToken(loginId); err != nil {
			log.Instance.Warn("get cache token failed", zap.Error(err))
			web.Send(ctx, constants.CodeServiceError)
			ctx.Abort()
			return
		} else if persistentToken != tokenStr {
			log.Instance.Warn("token not match", zap.String("token", tokenStr), zap.String("persistentToken", persistentToken))
			web.Send(ctx, constants.CodeTokenExpired)
			ctx.Abort()
			return
		}

		if token.Instance.NeedFlush(claims) {
			newTokenStr, err := token.Instance.FlushToken(claims)
			if err != nil {
				log.Instance.Warn("token flush failed", zap.Error(err))
				web.Send(ctx, constants.CodeServiceError)
				ctx.Abort()
				return
			}
			ctx.Header(constants.WebUserToken, newTokenStr)
			web.TokenBindLoginID(loginId, tokenStr, time.Duration(claims.ExpiresAt.Unix()-time.Now().Unix()))
		}

		ctx.Next()
	}
}
