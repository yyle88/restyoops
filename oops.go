package restyoops

import (
	"errors"
	"time"

	"github.com/yyle88/must"
)

// Oops represents a structured HTTP operation outcome
// Oops 代表结构化的 HTTP 操作结果
type Oops struct {
	Kind        Kind          // Classification // 分类
	StatusCode  int           // HTTP status code // HTTP 状态码
	ContentType string        // Response Content-Type // 响应 Content-Type
	Cause       error         // Wrapped outcome // 被包装的结果
	Retryable   bool          // Can be resolved via retries // 是否可通过重试解决
	WaitTime    time.Duration // Suggested wait time // 建议等待时间
}

// IsRetryable checks if retrying is recommended
// IsRetryable 检查是否建议重试
func (o *Oops) IsRetryable() bool {
	return o.Retryable
}

// NewOops creates an Oops with the specified params
// NewOops 使用指定的参数创建一个 Oops
func NewOops(kind Kind, statusCode int, cause error, retryable bool) *Oops {
	must.Nice(kind)
	must.In(kind, []Kind{KindUnknown, KindNetwork, KindHttp, KindParse, KindBlock, KindBusiness})
	return &Oops{
		Kind:        kind,
		StatusCode:  statusCode,
		ContentType: "",
		Cause:       must.Cause(cause),
		Retryable:   retryable,
		WaitTime:    0,
	}
}

// WithWaitTime sets the suggested wait time and returns the Oops
// WithWaitTime 设置建议等待时间并返回 Oops
func (o *Oops) WithWaitTime(d time.Duration) *Oops {
	o.WaitTime = d
	return o
}

// WithContentType sets the content type and returns the Oops
// WithContentType 设置内容类型并返回 Oops
func (o *Oops) WithContentType(contentType string) *Oops {
	o.ContentType = contentType
	return o
}

// NewUnknown creates an Oops indicating unknown issue
// NewUnknown 创建一个表示未知问题的 Oops
func NewUnknown() *Oops {
	return NewOops(KindUnknown, 0, errors.New(string(KindUnknown)), false)
}
