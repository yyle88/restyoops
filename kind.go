// Package restyoops: Structured HTTP operation fault classification with retryable semantics
// Oops! See if restyv2 response is retryable
// Provides Kind enum, Oops struct, and Detect function to classify HTTP response outcomes
//
// restyoops: 结构化 HTTP 操作故障分类，带有可重试语义
// Oops! 检查 restyv2 响应是否可重试
// 提供 Kind 枚举、Oops 结构体和 Detect 函数来分类 HTTP 响应结果
package restyoops

// Kind represents the classification of an HTTP operation outcome
// Used to categorize outcomes into actionable groups
//
// Kind 代表 HTTP 操作结果的分类
// 用于将结果分类为可操作的组
type Kind string

const (
	// KindUnknown indicates unclassified issues
	// KindUnknown 表示未分类的问题
	KindUnknown Kind = "UNKNOWN"

	// KindNetwork indicates network-level issues like timeout, DNS, TCP, TLS
	// Outcomes: connection reset, deadline exceeded, no such host
	// KindNetwork 表示网络层问题，如超时、DNS、TCP、TLS
	// 结果：连接重置、截止时间超时、无此主机
	KindNetwork Kind = "NETWORK"

	// KindHttp indicates HTTP status code issues (4xx/5xx)
	// Outcomes: 429 rate limit, 500 internal, 502/503/504 gateway issues
	// KindHttp 表示 HTTP 状态码问题（4xx/5xx）
	// 结果：429 限流、500 内部问题、502/503/504 网关问题
	KindHttp Kind = "HTTP"

	// KindParse indicates response parsing issues
	// Outcomes: JSON unmarshal failed, unexpected content type
	// KindParse 表示响应解析问题
	// 结果：JSON 反序列化失败、意外的内容类型
	KindParse Kind = "PARSE"

	// KindBlock indicates request was blocked (captcha, WAF, login redirect)
	// Outcomes: 403 with HTML, 200 with captcha page, redirect to login
	// KindBlock 表示请求被阻止（验证码、WAF、登录重定向）
	// 结果：403 带 HTML、200 带验证码页面、重定向到登录
	KindBlock Kind = "BLOCK"

	// KindBusiness indicates business logic issues (HTTP 200 but business code != 0)
	// Outcomes: rate limited, insufficient balance, invalid params
	// KindBusiness 表示业务逻辑问题（HTTP 200 但业务码 != 0）
	// 结果：限流、余额不足、参数无效
	KindBusiness Kind = "BUSINESS"
)

// String returns the string representation of Kind
// String 返回 Kind 的字符串表示
func (k Kind) String() string {
	return string(k)
}

// IsUnknown checks if Kind indicates unknown issues
// IsUnknown 检查 Kind 是否表示未知问题
func (k Kind) IsUnknown() bool {
	return k == KindUnknown
}

// IsNetwork checks if Kind indicates network issues
// IsNetwork 检查 Kind 是否表示网络问题
func (k Kind) IsNetwork() bool {
	return k == KindNetwork
}

// IsHttp checks if Kind indicates HTTP status issues
// IsHttp 检查 Kind 是否表示 HTTP 状态问题
func (k Kind) IsHttp() bool {
	return k == KindHttp
}

// IsParse checks if Kind indicates parsing issues
// IsParse 检查 Kind 是否表示解析问题
func (k Kind) IsParse() bool {
	return k == KindParse
}

// IsBlock checks if Kind indicates blocked/captcha issues
// IsBlock 检查 Kind 是否表示被阻止/验证码问题
func (k Kind) IsBlock() bool {
	return k == KindBlock
}

// IsBusiness checks if Kind indicates business logic issues
// IsBusiness 检查 Kind 是否表示业务逻辑问题
func (k Kind) IsBusiness() bool {
	return k == KindBusiness
}
