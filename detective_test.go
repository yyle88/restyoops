package restyoops_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/restyoops"
)

func TestDetective_Detect_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"msg":"ok"}`))
	}))
	defer server.Close()

	detective := restyoops.NewDetective(restyoops.NewConfig())
	response, oopsIssue := detective.Detect(resty.New().R().Get(server.URL))
	require.Nil(t, oopsIssue) // success returns nil
	require.NotNil(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode())
}

func TestDetective_Detect_HTTP500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	detective := restyoops.NewDetective(restyoops.NewConfig())
	response, oopsIssue := detective.Detect(resty.New().R().Get(server.URL))
	require.NotNil(t, oopsIssue)
	require.Equal(t, restyoops.KindHttp, oopsIssue.Kind)
	require.Equal(t, 500, oopsIssue.StatusCode)
	require.True(t, oopsIssue.Retryable)
	require.NotNil(t, response)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode())
}
