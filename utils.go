package kash

import "fmt"

type ErrorCode uint32

const (
	ErrorCacheNotFound ErrorCode = 1

)

type Error struct {
	ErrMsg string
	ErrCode ErrorCode
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrCode, e.ErrMsg)
}