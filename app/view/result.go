package view

// Result 统一响应结构
type Result struct {
	// 结果类型 0-无数据 1-有数据
	Result int `json:"result"`
	// 结果码
	Code int `json:"code"`
	// 结果提示
	Msg string `json:"msg"`
	// 数据
	Data interface{} `json:"data,omitempty"`
}

// Code 结果码类型
type Code int

// 结果码常量
const (
	// CodeSuccess 成功
	CodeSuccess Code = 1000
	// CodeFailed 失败
	CodeFailed Code = 1001
	// CodeInvalidBody 非法参数
	CodeInvalidBody Code = 1002
	// CodeLoseParam 参数缺失
	CodeLoseParam Code = 1003
	// CodeInvalidParams 参数验证失败
	CodeInvalidParams Code = 1004
	// CodeNoAuth 身份鉴权失败
	CodeNoAuth Code = 1006
	// CodeInvalidToken Token失效
	CodeInvalidToken Code = 1007
	// CodeSystemError 系统错误
	CodeSystemError Code = 9999
)

// 结果类型常量
const (
	// ResultNone 无数据
	ResultNone int = 0
	// ResultHave 有数据
	ResultHave int = 1
)

// ErrorResult 错误返回
func ErrorResult(code Code) *Result {
	return &Result{
		Result: ResultNone,
		Code:   int(code),
		Msg:    CodeString(code),
	}
}

// ErrorResultWithMsg 自定义错误返回
func ErrorResultWithMsg(code Code, msg string) *Result {
	return &Result{
		Result: ResultNone,
		Code:   int(code),
		Msg:    msg,
	}
}

// SuccessResult 成功返回
func SuccessResult(data interface{}) *Result {
	return &Result{
		Result: ResultHave,
		Code:   int(CodeSuccess),
		Msg:    CodeString(CodeSuccess),
		Data:   data,
	}
}

// SuccessResultWithoutData 成功返回（无数据）
func SuccessResultWithoutData() *Result {
	return &Result{
		Result: ResultNone,
		Code:   int(CodeSuccess),
		Msg:    CodeString(CodeSuccess),
	}
}

// CodeString 结果码字符串
func CodeString(code Code) string {
	switch code {
	case CodeSuccess:
		return "成功"
	case CodeFailed:
		return "失败"
	case CodeInvalidBody:
		return "非法参数"
	case CodeLoseParam:
		return "参数缺失"
	case CodeInvalidParams:
		return "参数验证失败"
	case CodeNoAuth:
		return "身份鉴权失败"
	case CodeInvalidToken:
		return "TOKEN失效"
	case CodeSystemError:
		return "系统错误"
	default:
		return "未知错误"
	}
}
