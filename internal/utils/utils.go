package utils

import "errors"

// ErrorsAs extracts error of type T from err
// ErrorsAs 从 err 中提取类型 T 的错误
func ErrorsAs[T any](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}
