package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/huskyrobotdog/toolbox-go/constants"
	"github.com/huskyrobotdog/toolbox-go/log"
	"github.com/huskyrobotdog/toolbox-go/web"
	"go.uber.org/zap"
)

type AnalyzeParamType int

const (
	AnalyzeParamTypeJson = iota
	AnalyzeParamTypePath
	AnalyzeParamTypeNone
)

func AnalyzeParam[P any](_type AnalyzeParamType, call func(*web.ContextData[P])) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if _type == AnalyzeParamTypeNone {
			call(&web.ContextData[P]{
				Ctx:     ctx,
				Data:    nil,
				LoginID: web.LoginID(ctx),
			})
			return
		}

		var param P
		switch _type {
		case AnalyzeParamTypeJson:
			if err := ctx.ShouldBindJSON(&param); err != nil {
				log.Instance.Warn("parse json failed", zap.Error(err))
				web.Send(ctx, constants.CodeRequestParamInvalid)
				return
			}
		case AnalyzeParamTypePath:
			if err := ctx.ShouldBindQuery(&param); err != nil {
				log.Instance.Warn("parse param failed", zap.Error(err))
				web.Send(ctx, constants.CodeRequestParamInvalid)
				return
			}
		default:
		}
		call(&web.ContextData[P]{
			Ctx:     ctx,
			Data:    &param,
			LoginID: web.LoginID(ctx),
		})
	}
}
