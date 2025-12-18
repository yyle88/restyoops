package restyoops_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/restyoops"
)

// TestDetect_Success tests Detect returns nil on HTTP 200 success
// TestDetect_Success 测试 Detect 在 HTTP 200 成功时返回 nil
func TestDetect_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"msg":"ok"}`))
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.Nil(t, oops) // success returns nil
}

// TestDetect_HTTP500 tests Detect classifies HTTP 500 as retryable
// TestDetect_HTTP500 测试 Detect 将 HTTP 500 分类为可重试
func TestDetect_HTTP500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.Equal(t, restyoops.KindHttp, oops.Kind)
	require.Equal(t, 500, oops.StatusCode)
	require.True(t, oops.Retryable)
}

// TestDetect_HTTP429 tests Detect classifies HTTP 429 as retryable
// TestDetect_HTTP429 测试 Detect 将 HTTP 429 分类为可重试
func TestDetect_HTTP429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.Equal(t, restyoops.KindHttp, oops.Kind)
	require.Equal(t, 429, oops.StatusCode)
	require.True(t, oops.Retryable)
}

// TestDetect_HTTP404 tests Detect classifies HTTP 404 as not retryable
// TestDetect_HTTP404 测试 Detect 将 HTTP 404 分类为不可重试
func TestDetect_HTTP404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.Equal(t, restyoops.KindHttp, oops.Kind)
	require.Equal(t, 404, oops.StatusCode)
	require.False(t, oops.Retryable)
}

// TestDetect_NetworkTimeout tests Detect classifies timeout as retryable network issue
// TestDetect_NetworkTimeout 测试 Detect 将超时分类为可重试的网络问题
func TestDetect_NetworkTimeout(t *testing.T) {
	oops := restyoops.Detect(restyoops.NewConfig(), nil, context.DeadlineExceeded)
	require.Equal(t, restyoops.KindNetwork, oops.Kind)
	require.True(t, oops.Retryable)
}

// TestDetect_NetworkCanceled tests Detect classifies canceled as retryable network issue
// TestDetect_NetworkCanceled 测试 Detect 将取消分类为可重试的网络问题
func TestDetect_NetworkCanceled(t *testing.T) {
	oops := restyoops.Detect(restyoops.NewConfig(), nil, context.Canceled)
	require.Equal(t, restyoops.KindNetwork, oops.Kind)
	require.True(t, oops.Retryable)
}

// TestDetect_UnknownError tests Detect classifies unknown errors as not retryable
// TestDetect_UnknownError 测试 Detect 将未知错误分类为不可重试
func TestDetect_UnknownError(t *testing.T) {
	oops := restyoops.Detect(restyoops.NewConfig(), nil, errors.New("some unknown issue"))
	require.Equal(t, restyoops.KindUnknown, oops.Kind)
	require.False(t, oops.Retryable)
}

// TestConfig_Override403Retryable tests Config can override HTTP 403 to be retryable
// TestConfig_Override403Retryable 测试 Config 可以覆盖 HTTP 403 为可重试
func TestConfig_Override403Retryable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	// Default: 403 not retryable
	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.Equal(t, restyoops.KindHttp, oops.Kind)
	require.False(t, oops.Retryable)

	// Config override: 403 retryable with 2s wait
	cfg := restyoops.NewConfig().WithStatusRetryable(403, true, 2*time.Second)
	oops = restyoops.Detect(cfg, resp, err)
	require.Equal(t, restyoops.KindHttp, oops.Kind)
	require.True(t, oops.Retryable)
	require.Equal(t, 2*time.Second, oops.WaitTime)
}

// TestConfig_Override500NotRetryable tests Config can override HTTP 500 to not be retryable
// TestConfig_Override500NotRetryable 测试 Config 可以覆盖 HTTP 500 为不可重试
func TestConfig_Override500NotRetryable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)

	// Default: 500 retryable
	oops := restyoops.Detect(restyoops.NewConfig(), resp, err)
	require.True(t, oops.Retryable)

	// Config override: 500 not retryable
	cfg := restyoops.NewConfig().WithStatusRetryable(500, false, 0)
	oops = restyoops.Detect(cfg, resp, err)
	require.False(t, oops.Retryable)
}
