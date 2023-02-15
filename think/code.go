package think

import (
	"fmt"
	"net/http"
)

//type Code uint32
//
//const (
//	Code_Min              Code = 100000
//	Code_Success          Code = 100200 // 成功
//	Code_SuccessAction    Code = 100201 // 成功行为
//	Code_Biz              Code = 100208 // 业务码
//	Code_AlterError       Code = 102400 // 简单错误
//	Code_ParamError       Code = 101400 // 参数错误
//	Code_NotFound         Code = 100404 // 记录错误
//	Code_Repeat           Code = 100405 // 重复操作, 表示已存在
//	Code_UnDone           Code = 103400 // 记录错误
//	Code_Forbidden        Code = 100403 // 禁止访问
//	Code_SignError        Code = 101403 // 签名错误
//	Code_Unauthorized     Code = 100401 // 未认证
//	Code_TooManyRequests  Code = 100429 // 请求过于频繁
//	Code_SystemSpaceError Code = 100500 // 系统空间错误  不外抛
//	Code_PanicError       Code = 100502 // 系统崩溃错误
//	Code_Ignore           Code = 101500 // 忽略
//	Code_Undefined        Code = 102500 // 未定义
//	Code_TimeOut          Code = 102504 // 超时
//	Code_Exception        Code = 103500 // 异常
//	Code_TypeError        Code = 104500 // 类型错误
//	Code_Unavailable      Code = 105500 // 不可达
//)

func (c Code) HttpCode() int {
	switch c {
	case Code_Success, Code_Biz:
		return http.StatusOK
	case Code_SuccessAction:
		return http.StatusCreated
	case Code_ParamError, Code_NotFound, Code_AlterError, Code_UnImpl, Code_Repeat:
		return http.StatusBadRequest

	case Code_Forbidden, Code_SignError:
		return http.StatusForbidden

	case Code_Unauthorized:
		return http.StatusUnauthorized

	case Code_TooManyRequests:
		return http.StatusTooManyRequests

	case Code_SystemSpaceError, Code_Ignore, Code_Undefined, Code_Exception, Code_TypeError, Code_Unavailable:
		return http.StatusInternalServerError

	case Code_TimeOut:
		return http.StatusGatewayTimeout

	default:
		return http.StatusInternalServerError

	}
}

func (c Code) ToString() string {
	switch c {
	case Code_Success, Code_SuccessAction:
		return "成功"
	case Code_ParamError:
		return "参数错误"
	case Code_NotFound:
		return "404"
	case Code_Repeat:
		return "重复操作"
	case Code_AlterError:
		return "提示错误"
	case Code_Biz:
		return "业务码"
	case Code_UnImpl:
		return "未实现"
	case Code_Forbidden:
		return "禁止访问"
	case Code_SignError:
		return "签名错误"
	case Code_Unauthorized:
		return "未认证"
	case Code_TooManyRequests:
		return "反问过于频繁"
	case Code_SystemSpaceError:
		return "系统空间错误"
	case Code_Ignore:
		return "忽略"
	case Code_Undefined:
		return "未定义"
	case Code_TimeOut:
		return "超时"
	case Code_Exception:
		return "异常"
	case Code_TypeError:
		return "错误类型"
	case Code_Unavailable:
		return "不可达"
	case Code_PanicError:
		return "服务内部错误"
	default:
		return fmt.Sprintf("未定义： %d ", c)
	}
}
