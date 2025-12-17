package restyoops

import "time"

// Oops represents a structured HTTP operation outcome
// Oops 代表结构化的 HTTP 操作结果
type Oops struct {
	Kind        Kind          // Classification // 分类
	StatusCode  int           // HTTP status code // HTTP 状态码
	Retryable   bool          // Can be resolved via retries // 是否可通过重试解决
	WaitTime    time.Duration // Suggested wait time // 建议等待时间
	Cause       error         // Wrapped outcome // 被包装的结果
	ContentType string        // Response Content-Type // 响应 Content-Type
}

// IsSuccess checks if the operation was success
// IsSuccess 检查操作是否成功
func (o *Oops) IsSuccess() bool {
	return o.Kind == KindSuccess
}

// IsRetryable checks if a retry is recommended
// IsRetryable 检查是否建议重试
func (o *Oops) IsRetryable() bool {
	return o.Retryable
}

// NewOops creates an Oops with the specified params
// NewOops 使用指定的参数创建一个 Oops
func NewOops(kind Kind, statusCode int, retryable bool, cause error) *Oops {
	return &Oops{
		Kind:       kind,
		StatusCode: statusCode,
		Retryable:  retryable,
		Cause:      cause,
	}
}

// NewSuccess creates an Oops indicating success
// NewSuccess 创建一个表示成功的 Oops
func NewSuccess() *Oops {
	return &Oops{
		Kind:      KindSuccess,
		Retryable: false,
	}
}

// NewUnknown creates an Oops indicating unknown issue
// NewUnknown 创建一个表示未知问题的 Oops
func NewUnknown() *Oops {
	return &Oops{
		Kind:      KindUnknown,
		Retryable: false,
	}
}
