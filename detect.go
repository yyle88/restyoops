package restyoops

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yyle88/must"
	"github.com/yyle88/restyoops/internal/utils"
)

// Detect classifies a resty response
// Detect 分类 resty 响应
func Detect(cfg *Config, resp *resty.Response, respCause error) *Oops {
	if respCause != nil {
		return detectNetworkOops(cfg, respCause)
	}

	must.Full(resp)

	statusCode := resp.StatusCode()
	contentType := resp.Header().Get("Content-Type")
	content := resp.Body()

	// Run custom content check
	// 运行自定义内容检查
	if check, ok := cfg.ContentChecks[statusCode]; ok {
		if oops := check(contentType, content); oops != nil {
			return oops
		}
	}

	// Check HTTP status code
	// 检查 HTTP 状态码
	if statusCode >= 400 {
		return detectDefaultHttpOops(cfg, statusCode, contentType)
	}

	// Success
	// 成功
	return &Oops{
		Kind:        KindSuccess,
		StatusCode:  statusCode,
		Retryable:   false,
		ContentType: contentType,
	}
}

// detectNetworkOops classifies network-level issues
// detectNetworkOops 分类网络层问题
func detectNetworkOops(cfg *Config, respCause error) *Oops {
	var kind Kind
	var defaultRetryable bool

	// Check specific types first (more specific before general)
	// 先检查具体类型（具体的在通用的前面）
	if errors.Is(respCause, context.DeadlineExceeded) || errors.Is(respCause, context.Canceled) {
		kind = KindNetwork
		defaultRetryable = true
	} else if dnsErr, ok := utils.ErrorsAs[*net.DNSError](respCause); ok {
		kind = KindNetwork
		defaultRetryable = !dnsErr.IsNotFound
	} else if _, ok := utils.ErrorsAs[*net.OpError](respCause); ok {
		kind = KindNetwork
		defaultRetryable = true
	} else if _, ok := utils.ErrorsAs[*url.Error](respCause); ok {
		kind = KindNetwork
		defaultRetryable = true
	} else if netErr, ok := utils.ErrorsAs[net.Error](respCause); ok {
		kind = KindNetwork
		defaultRetryable = netErr.Timeout()
	} else {
		kind = KindUnknown
		defaultRetryable = false
	}

	retryable, waitTime := applyOption(cfg, kind, 0, defaultRetryable)
	return &Oops{
		Kind:      kind,
		Cause:     respCause,
		Retryable: retryable,
		WaitTime:  waitTime,
	}
}

// detectDefaultHttpOops classifies HTTP status code issues
// detectDefaultHttpOops 分类 HTTP 状态码问题
func detectDefaultHttpOops(cfg *Config, statusCode int, contentType string) *Oops {
	var defaultRetryable bool
	switch statusCode {
	case http.StatusTooManyRequests, http.StatusRequestTimeout: // 429, 408
		defaultRetryable = true
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout: // 502, 503, 504
		defaultRetryable = true
	case http.StatusInternalServerError: // 500
		defaultRetryable = true
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden: // 400, 401, 403
		defaultRetryable = false
	case http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity: // 404, 409, 422
		defaultRetryable = false
	default:
		defaultRetryable = statusCode >= 500
	}

	retryable, waitTime := applyOption(cfg, KindHttp, statusCode, defaultRetryable)
	return &Oops{
		Kind:        KindHttp,
		StatusCode:  statusCode,
		Retryable:   retryable,
		WaitTime:    waitTime,
		ContentType: contentType,
	}
}

// applyOption applies config overrides and returns (retryable, waitTime)
// applyOption 应用配置覆盖并返回 (retryable, waitTime)
func applyOption(cfg *Config, kind Kind, statusCode int, defaultRetryable bool) (bool, time.Duration) {
	must.Full(cfg)
	if statusCode > 0 {
		if opt, ok := cfg.StatusOptions[statusCode]; ok {
			waitTime := opt.WaitTime
			if waitTime == 0 {
				waitTime = cfg.DefaultWait
			}
			return opt.Retryable, waitTime
		}
	}

	if opt, ok := cfg.KindOptions[kind]; ok {
		waitTime := opt.WaitTime
		if waitTime == 0 {
			waitTime = cfg.DefaultWait
		}
		return opt.Retryable, waitTime
	}

	return defaultRetryable, cfg.DefaultWait
}
