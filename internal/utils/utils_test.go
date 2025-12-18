package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/restyoops/internal/utils"
)

type errorTypeA struct{ msg string }

func (e *errorTypeA) Error() string { return e.msg }

type errorTypeB struct{ msg string }

func (e *errorTypeB) Error() string { return e.msg }

func TestErrorsAs(t *testing.T) {
	errA := &errorTypeA{msg: "aaa"}

	got, ok := utils.ErrorsAs[*errorTypeA](errA)
	require.True(t, ok)
	require.Equal(t, errA, got)

	_, ok = utils.ErrorsAs[*errorTypeB](errA)
	require.False(t, ok)
}
