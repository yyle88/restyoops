package restyoops

import (
	"github.com/go-resty/resty/v2"
)

// Detective wraps Config and provides a convenient API
// Detective 封装 Config 并提供便捷的 API
type Detective struct {
	cfg *Config
}

func NewDetective(cfg *Config) *Detective {
	return &Detective{
		cfg: cfg,
	}
}

func (c *Detective) Detect(resp *resty.Response, respCause error) (*resty.Response, *OopsIssue) {
	oops := Detect(c.cfg, resp, respCause)
	return resp, oops
}
