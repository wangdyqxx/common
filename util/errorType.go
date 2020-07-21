package util

import (
	"errors"
)

var (
	ErrorNotFound = errors.New("not found")
	ErrorNil      = errors.New("is nil")
)

type Status struct {
	Code    int32
	Message string
}

func (m *Status) GetStatus() *Status {
	return (*Status)(m)
}

func New(c int32, msg string) *Status {
	return &Status{Code: int32(c), Message: msg}
}

func FromError(err error) (s *Status, ok bool) {
	if err == nil {
		return &Status{Code: 0}, true
	}
	if se, ok := err.(interface{ GetStatus() *Status }); ok {
		return se.GetStatus(), true
	}
	return New(2, err.Error()), false
}
