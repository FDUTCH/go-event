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

// Subscribe subscribes to event with minimal priority.
func (bus *Bus) Subscribe[T any](listener func(ev T)) {
	bus.SubscribeWithPriority(listener, 0)
}

// SubscribeWithPriority subscribes to event with priority passed.
func (bus *Bus) SubscribeWithPriority[T any](listener func(ev T), priority uint32) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	eventType := reflect.TypeFor[T]()
	handlers, ok := bus.handlers[eventType]
	if ok {
		queue := handlers.(*priorityQueue[T])
		queue.add(listener, priority)
	} else {
		q := new(priorityQueue[T])
		q.add(listener, priority)
		bus.handlers[eventType] = q
	}
}

// Unsubscribe unsubscribes from event.
func (bus *Bus) Unsubscribe[T any](fn func(ev T)) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	eventType := reflect.TypeFor[T]()
	queue, ok := bus.handlers[eventType]
	if !ok {
		return
	}
	queue.(*priorityQueue[T]).delete(fn)

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

	handlers := slices.Clone(rawHandlers.(*priorityQueue[T]).callers)
	bus.mu.RUnlock()

	for i := range handlers {
		handlers[i].fn(event)
	}
}

type priorityQueue[T any] struct {
	callers       []caller[T]
	globalCounter uint32
}

func (q *priorityQueue[T]) add(fn func(T), priority uint32) {
	toInsert := caller[T]{
		priority: (uint64(priority) << 32) + uint64(q.globalCounter),
		fn:       fn,
	}
	if len(q.callers) == 0 {
		q.callers = []caller[T]{toInsert}
		return
	}

	for idx, ca := range q.callers {
		if ca.priority < toInsert.priority {
			q.callers = slices.Insert(q.callers, idx, toInsert)
			break
		}
	}

	q.globalCounter++ // don't know if I really should keep this.
}

func (q *priorityQueue[T]) delete(fn func(T)) {
	addr := reflect.ValueOf(fn).Pointer()
	q.callers = slices.DeleteFunc(q.callers, func(a caller[T]) bool {
		return addr == reflect.ValueOf(a.fn).Pointer()
	})
}

type caller[T any] struct {
	priority uint64
	fn       func(T)
}
