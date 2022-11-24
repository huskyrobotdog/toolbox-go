package constants

import "github.com/huskyrobotdog/toolbox-go/types"

var (
	CodeOK                  = &types.Code{Code: 0, Message: ""}
	CodeServiceError        = &types.Code{Code: 1, Message: "service error"}
	CodeNotLogin            = &types.Code{Code: 2, Message: "Not logged in or login timed out"}
	CodeRequestParamInvalid = &types.Code{Code: 3, Message: "Invalid request parameter"}
)

var (
	CodeTokenInvalidation = &types.Code{Code: 100, Message: "token invalidation"}
	CodeTokenIsMalformed  = &types.Code{Code: 101, Message: "malformed token"}
	CodeTokenExpired      = &types.Code{Code: 102, Message: "token expires"}
	CodeTokenNotActiveYet = &types.Code{Code: 103, Message: "token not active yet"}
)
