package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huskyrobotdog/toolbox-go/cache"
	"github.com/huskyrobotdog/toolbox-go/constants"
	"github.com/huskyrobotdog/toolbox-go/id"
	"github.com/huskyrobotdog/toolbox-go/types"
)

func Send(ctx *gin.Context, code *types.Code, data ...interface{}) {
	m := gin.H{}
	if code != nil {
		m["code"] = code.Code
		m["message"] = code.Message
	} else {
		m["code"] = 0
	}
	if len(data) > 0 {
		m["data"] = data[0]
	}
	ctx.JSON(http.StatusOK, m)
}

func RecordLoginID(ctx *gin.Context, id id.ID) {
	ctx.Set(constants.WebUserTokenID, id)
}

func LoginID(ctx *gin.Context) id.ID {
	loginId, exists := ctx.Get(constants.WebUserTokenID)
	if !exists {
		return 0
	}
	return loginId.(id.ID)
}

func TokenBindLoginID(loginId id.ID, token string, expr time.Duration) error {
	if expr <= 0 {
		return errors.New("`expr` can not less then or equals zero")
	}
	return cache.Instance.SetEX(context.Background(), constants.CacheKeyUserTokenPrefix+loginId.String(), token, expr).Err()
}

func LoginIDTokenDestory(loginId id.ID) error {
	return cache.Instance.Del(context.Background(), constants.CacheKeyUserTokenPrefix+loginId.String()).Err()
}

func LoginIDToken(loginId id.ID) (string, error) {
	return cache.Instance.Get(context.Background(), constants.CacheKeyUserTokenPrefix+loginId.String()).Result()
}

type ContextData[D any] struct {
	Ctx     *gin.Context
	Data    *D
	LoginID id.ID
}

func (curr *ContextData[D]) Send(code *types.Code, data ...interface{}) {
	Send(curr.Ctx, code, data...)
}

func (curr *ContextData[D]) RequestIP() string {
	reqIP := curr.Ctx.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

type ErrorWrap struct {
	Code    *types.Code
	AttaMsg *string
}

func assert(logic bool, code *types.Code, msg *string, v ...interface{}) {
	if logic {
		return
	}
	if msg == nil {
		panic(&ErrorWrap{
			Code:    code,
			AttaMsg: nil,
		})
	} else {
		_m := fmt.Sprintf(*msg, v...)
		panic(&ErrorWrap{
			Code:    code,
			AttaMsg: &_m,
		})
	}
}

func Assert(logic bool, code *types.Code) {
	assert(logic, code, nil)
}

func AssertWithMsg(logic bool, code *types.Code, msg string, v ...interface{}) {
	assert(logic, code, &msg, v...)
}

func AssertWithError(logic bool, code *types.Code, err error) {
	if err != nil {
		_m := err.Error()
		assert(logic, code, &_m)
	} else {
		assert(logic, code, nil)
	}
}
