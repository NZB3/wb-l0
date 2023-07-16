package storage

import "errors"

var (
	ErrDBNotExists    = errors.New("db not exists")
	ErrCacheNotExists = errors.New("cache not exists")
)
