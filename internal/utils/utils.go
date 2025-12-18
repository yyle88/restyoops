// Package utils provides utilities supporting the restyoops package
// Implements generic support functions used across different modules
//
// utils 提供支持 restyoops 包的工具
// 实现在不同模块间使用的泛型辅助函数
package utils

import "errors"

// ErrorsAs extracts error of type T from err
// ErrorsAs 从 err 中提取类型 T 的错误
func ErrorsAs[T any](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}
