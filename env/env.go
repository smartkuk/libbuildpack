package os

import (
	"os"
)

type Env interface {
	Get(key string) string
	Set(key, value string) error
}

type real struct{}

func New() Env {
	return &real{}
}

func (e *real) Get(key string) string {
	return os.Getenv(key)
}

func (e *real) Set(key, value string) error {
	return os.Setenv(key, value)
}

type mock struct {
	data map[string]string
}

func Mock() Env {
	return &mock{data: make(map[string]string)}
}

func (e *mock) Get(key string) string {
	return e.data[key]
}

func (e *mock) Set(key, value string) error {
	e.data[key] = value
	return nil
}
