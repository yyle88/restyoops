package restyoops

import "time"

// StatusOption holds retryable and wait time settings
// StatusOption 保存可重试和等待时间设置
type StatusOption struct {
	Retryable bool
	WaitTime  time.Duration
}

// KindOption holds retryable and wait time settings
// KindOption 保存可重试和等待时间设置
type KindOption struct {
	Retryable bool
	WaitTime  time.Duration
}

// ContentCheckFunc checks content and returns Oops if matched, nil otherwise
// ContentCheckFunc 检查内容，匹配时返回 Oops，否则返回 nil
type ContentCheckFunc func(contentType string, content []byte) *Oops

// Config holds customizable detection settings
// Config 保存可自定义的检测设置
type Config struct {
	StatusOptions map[int]*StatusOption
	KindOptions   map[Kind]*KindOption
	DefaultWait   time.Duration            // default wait time // 默认等待时间
	ContentChecks map[int]ContentCheckFunc // custom content checks // 自定义内容检查
}

// NewConfig creates a Config with sensible defaults
// NewConfig 创建带有合理默认值的 Config
func NewConfig() *Config {
	return &Config{
		StatusOptions: make(map[int]*StatusOption),
		KindOptions:   make(map[Kind]*KindOption),
		DefaultWait:   time.Second, // 1s default
		ContentChecks: make(map[int]ContentCheckFunc),
	}
}

// WithStatusRetryable sets retryable and wait time based on status code
// WithStatusRetryable 基于状态码设置可重试和等待时间
func (c *Config) WithStatusRetryable(statusCode int, retryable bool, waitTime time.Duration) *Config {
	c.StatusOptions[statusCode] = &StatusOption{
		Retryable: retryable,
		WaitTime:  waitTime,
	}
	return c
}

// WithKindRetryable sets retryable and wait time based on Kind
// WithKindRetryable 基于 Kind 设置可重试和等待时间
func (c *Config) WithKindRetryable(kind Kind, retryable bool, waitTime time.Duration) *Config {
	c.KindOptions[kind] = &KindOption{
		Retryable: retryable,
		WaitTime:  waitTime,
	}
	return c
}

// WithDefaultWait sets the default wait time
// WithDefaultWait 设置默认等待时间
func (c *Config) WithDefaultWait(d time.Duration) *Config {
	c.DefaultWait = d
	return c
}

// WithContentCheck adds a custom content check
// WithContentCheck 添加自定义内容检查
func (c *Config) WithContentCheck(statusCode int, check ContentCheckFunc) *Config {
	c.ContentChecks[statusCode] = check
	return c
}
