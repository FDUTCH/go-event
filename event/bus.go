package event

import (
	"reflect"
	"slices"
	"sync"
)

// Bus is a thread-safe bus of events.
type Bus struct {
	mu       sync.RWMutex
	handlers map[reflect.Type]any
}

// NewBus creates new Bus.
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[reflect.Type]any),
	}
}

// Subscribe subscribes to event.
func (bus *Bus) Subscribe[T any](listener func(ev T)) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	eventType := reflect.TypeOf((*T)(nil)).Elem()
	handlers, ok := bus.handlers[eventType]
	if ok {
		bus.handlers[eventType] = append(handlers.([]func(T)), listener)
	} else {
		bus.handlers[eventType] = []func(T){listener}
	}
}

// Unsubscribe unsubscribes from event.
func (bus *Bus) Unsubscribe[T any](fn func(ev T)) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	// there a way to shoot your self in the foot, but idc.
	eventType := reflect.TypeOf((*T)(nil)).Elem()
	addr := reflect.ValueOf(fn).Pointer()
	handlers, ok := bus.handlers[eventType]
	if !ok {
		return
	}
	bus.handlers[eventType] = slices.DeleteFunc(handlers.([]func(T)), func(a func(T)) bool {
		return addr == reflect.ValueOf(a).Pointer()
	})
}

// Publish publishes event to the Bus.
func (bus *Bus) Publish[T any](event T) {
	bus.mu.RLock()
	eventType := reflect.TypeOf(event)
	rawHandlers, exists := bus.handlers[eventType]
	if !exists {
		bus.mu.RUnlock()
		return
	}

	handlers := slices.Clone(rawHandlers.([]func(T)))
	bus.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}
