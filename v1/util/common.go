package util

import "errors"

// 定义错误类型
var (
	ErrInvalidArgument = errors.New("invalid arugment")

	ErrInterrupted = errors.New("interrupted")
)
