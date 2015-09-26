package cache

import ()

type Store interface {
	Get(key string, val interface{}) error
	Set(key string, val interface{}, expire uint) error
	Delete(key string) error
	Flush() error
}
