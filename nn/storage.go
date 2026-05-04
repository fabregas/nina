package nn

import (
	"errors"
)

type LocalStorage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Remove(key string)
	Clear()
	SetJSON(key string, value any) error
	GetJSON(key string, value any) error
	Watch(comp Component, key string, callback func(newValue string))
	unwatchAll(comp Component)
}

var NotFound = errors.New("not found")
