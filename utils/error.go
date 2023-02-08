package utils

import (
	"net/http"
	"strconv"
)

// 错误
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Stack      error  `json:"-"`
}

func NewError(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       strconv.Itoa(statusCode),
		Message:    message,
	}
}

func (e *Error) Clone() *Error {
	return &Error{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    e.Message,
	}
}

func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) WithStack(stack error) *Error {
	e.Stack = stack
	return e
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GetError() error {
	return e
}

func (e *Error) GinError() (int, error) {
	return e.StatusCode, e
}

func (e *Error) JsonError() (int, interface{}) {
	return e.StatusCode, e
}

func (e *Error) Panic() {
	if e != nil {
		panic(e)
	}
}

func (e *Error) PanicError(err error) {
	if err != nil && e != nil {
		e.Stack = err
		e.Panic()
	}
}

var (
	// 200
	Success = NewError(http.StatusOK, http.StatusText(http.StatusOK))
	// 500
	InternalServerError = NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	NotImplemented      = NewError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	ServiceUnavailable  = NewError(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable)) // 访问了服务器不存在的资源
	// 400
	BadRequest   = NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))     // 客户端请求错误
	Unauthorized = NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)) // 需要认证
	NotFound     = NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))         // 服务器不可访问
	Forbidden    = NewError(http.StatusForbidden, http.StatusText(http.StatusForbidden))       // 禁止访问
)

// 客户端错误
type ClientErrorCode string

func (code ClientErrorCode) Panic() {
	ClientError(code).Panic()
}

func (code ClientErrorCode) Error() *Error {
	return ClientError(code)
}

func (code ClientErrorCode) PanicError(err error) {
	if err != nil {
		code.PanicError(err)
	}
}

func ClientError(code ClientErrorCode) *Error {
	return BadRequest.Clone().WithCode(string(code)).WithMessage(string(code))
}

// 服务端错误
type ServerErrorCode string

func ServerError(code ServerErrorCode) *Error {
	return InternalServerError.Clone().WithCode(string(code)).WithMessage(string(code))
}

func (code ServerErrorCode) Panic() {
	ServerError(code).Panic()
}

func (code ServerErrorCode) Error() *Error {
	return ServerError(code)
}

func (code ServerErrorCode) PanicError(err error) {
	if err != nil {
		code.PanicError(err)
	}
}

var (
	ErrInputParamsInvalid      ClientErrorCode = "input.paramsInvalid"      // 参数错误
	ErrDataDuplicate           ServerErrorCode = "data.duplicate"           // 数据重复
	ErrDataExists              ServerErrorCode = "data.exists"              // 数据已存在
	ErrDataNoFound             ServerErrorCode = "data.notFound"            // 数据不存在
	ErrAuthAccessTokenInvalid  ClientErrorCode = "auth.accessTokenInvalid"  // 访问令牌错误
	ErrAuthLoginFail           ServerErrorCode = "auth.loginFail"           // 登录失败
	ErrAuthRefreshTokenInvalid ServerErrorCode = "auth.refreshTokenInvalid" // 刷新令牌失败
)

func EnsureNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func EnsureFoundNoError(found bool, err error) {
	if err != nil {
		panic(err)
	}
	if !found {
		NotFound.Panic()
	}
}

func EnsureCountNoError(count int64, err error) {
	if err != nil {
		panic(err)
	}
	if count <= 0 {
		ErrDataNoFound.Panic()
	}
}

func EnsureIgnoreFirstNoError(count any, err error) {
	if err != nil {
		panic(err)
	}
}
