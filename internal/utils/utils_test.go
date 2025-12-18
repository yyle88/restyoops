package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/restyoops/internal/utils"
)

// errorTypeA is a test struct implementing error interface
// errorTypeA 是实现 error 接口的测试结构体
type errorTypeA struct{ msg string }

// Error returns the error message
// Error 返回错误消息
func (e *errorTypeA) Error() string { return e.msg }

// errorTypeB is a test struct implementing error interface
// errorTypeB 是实现 error 接口的测试结构体
type errorTypeB struct{ msg string }

// Error returns the error message
// Error 返回错误消息
func (e *errorTypeB) Error() string { return e.msg }

// TestErrorsAs tests ErrorsAs extracts typed errors and returns false on type mismatch
// TestErrorsAs 测试 ErrorsAs 提取类型化错误并在类型不匹配时返回 false
func TestErrorsAs(t *testing.T) {
	errA := &errorTypeA{msg: "aaa"}

	got, ok := utils.ErrorsAs[*errorTypeA](errA)
	require.True(t, ok)
	require.Equal(t, errA, got)

	_, ok = utils.ErrorsAs[*errorTypeB](errA)
	require.False(t, ok)
}
