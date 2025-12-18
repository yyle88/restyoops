package restyoops

import (
	"github.com/go-resty/resty/v2"
	"github.com/yyle88/must"
)

// Detective wraps Config and provides a convenient API
// Detective 封装 Config 并提供便捷的 API
type Detective struct {
	cfg *Config
}

// NewDetective creates a Detective with the specified Config
// NewDetective 使用指定的 Config 创建 Detective
func NewDetective(cfg *Config) *Detective {
	return &Detective{
		cfg: must.Full(cfg),
	}
}

// Detect classifies a resty response and returns both response and oops issue
// Detect 分类 resty 响应并返回响应和 oops 问题
func (c *Detective) Detect(resp *resty.Response, respCause error) (*resty.Response, *OopsIssue) {
	oops := Detect(c.cfg, resp, respCause)
	if oops != nil {
		must.Nice(oops.Kind)
		must.Wrong(oops.Cause)
	}
	return resp, oops
}
