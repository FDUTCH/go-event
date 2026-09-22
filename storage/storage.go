package storage

import (
	"reflect"
)

// Storage stores typed data.
type Storage struct {
	storage map[reflect.Type]any
}

// NewStorage create new storage.
func NewStorage() *Storage {
	return &Storage{storage: make(map[reflect.Type]any)}
}

// Set sets value.
func (s *Storage) Set[T any](val T) {
	s.storage[reflect.TypeOf(val)] = val
}

// Get returns the value.
func (s *Storage) Get[T any]() (T, bool) {
	var val T
	v, ok := s.storage[reflect.TypeFor[T]()]
	if !ok {
		return val, false
	}
	val, ok = v.(T)
	return val, ok
}

// GetOr returns the value or a value 'or'.
func (s *Storage) GetOr[T any](or T) T {
	val, ok := s.Get[T]()
	if !ok {
		return or
	}
	return val
}
