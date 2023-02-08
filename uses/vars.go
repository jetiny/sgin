package uses

import "jetiny/sgin/utils"

const (
	gErrAppCodeInvalid    utils.ClientErrorCode = "client.noAppCode"   // 没有appCode
	gErrAuthLoginRequired utils.ClientErrorCode = "auth.loginRequired" // 需要登录
	gErrAuthTokenExpired  utils.ClientErrorCode = "auth.tokenExpired"  // 访问令牌已过期
)
