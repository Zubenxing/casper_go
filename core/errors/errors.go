package errors

// 错误码定义
const (
	// 成功
	SUCCESS = 0

	// 通用错误 1000-1999
	ErrBadRequest     = 1000 // 请求参数错误
	ErrUnauthorized   = 1001 // 未授权
	ErrForbidden      = 1002 // 禁止访问
	ErrNotFound       = 1003 // 资源不存在
	ErrServerError    = 1004 // 服务器内部错误
	ErrTooManyRequest = 1005 // 请求过于频繁

	// 认证相关错误 2000-2999
	ErrInvalidCredentials = 2000 // 用户名或密码错误
	ErrInvalidToken       = 2001 // 无效的令牌
	ErrTokenExpired       = 2002 // 令牌已过期
	ErrTokenRevoked       = 2003 // 令牌已失效
	ErrUserDisabled       = 2004 // 用户已被禁用
	ErrUserNotFound       = 2005 // 用户不存在

	// 证书相关错误 3000-3999
	ErrCertURLExists   = 3000 // URL 已存在
	ErrCertNotFound    = 3001 // 证书记录不存在
	ErrCertCheckFailed = 3002 // 证书检查失败
	ErrCertInvalid     = 3003 // 证书无效
	ErrCertExpired     = 3004 // 证书已过期

	// 数据库相关错误 4000-4999
	ErrDatabaseConnection = 4000 // 数据库连接失败
	ErrDatabaseQuery      = 4001 // 数据库查询失败
	ErrDatabaseInsert     = 4002 // 数据插入失败
	ErrDatabaseUpdate     = 4003 // 数据更新失败
	ErrDatabaseDelete     = 4004 // 数据删除失败
)

// 错误码对应的消息
var ErrorMessages = map[int]string{
	SUCCESS:           "success",
	ErrBadRequest:     "请求参数错误",
	ErrUnauthorized:   "未授权",
	ErrForbidden:      "禁止访问",
	ErrNotFound:       "资源不存在",
	ErrServerError:    "服务器内部错误",
	ErrTooManyRequest: "请求过于频繁，请稍后再试",

	ErrInvalidCredentials: "用户名或密码错误",
	ErrInvalidToken:       "无效的令牌",
	ErrTokenExpired:       "令牌已过期",
	ErrTokenRevoked:       "令牌已失效",
	ErrUserDisabled:       "用户已被禁用",
	ErrUserNotFound:       "用户不存在",

	ErrCertURLExists:   "该 URL 已存在监控列表",
	ErrCertNotFound:    "证书记录不存在",
	ErrCertCheckFailed: "证书检查失败",
	ErrCertInvalid:     "证书无效",
	ErrCertExpired:     "证书已过期",

	ErrDatabaseConnection: "数据库连接失败",
	ErrDatabaseQuery:      "数据库查询失败",
	ErrDatabaseInsert:     "数据插入失败",
	ErrDatabaseUpdate:     "数据更新失败",
	ErrDatabaseDelete:     "数据删除失败",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Error 自定义错误结构
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// New 创建新的错误
func New(code int, message ...string) *Error {
	msg := GetMessage(code)
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}
