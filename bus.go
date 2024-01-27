package common

import (
	"sync"
)

type Bus[T any] struct {
	buf         *Ring[T]
	counter     int
	close       chan struct{}
	subscribers map[int]chan struct{}

	mu   sync.Mutex
	once sync.Once
}

func NewBus[T any]() *Bus[T] {
	return &Bus[T]{
		buf:         NewRing[T](),
		close:       make(chan struct{}),
		subscribers: make(map[int]chan struct{}),
	}
}

type Subscriber[T any] struct {
	bus    *Bus[T]
	id     int
	notify chan struct{}
	once   sync.Once
}

func (s *Subscriber[T]) Closed() <-chan struct{} {
	return s.bus.close
}

func (s *Subscriber[T]) Notify() <-chan struct{} {
	return s.notify
}

func (s *Subscriber[T]) Unsubscribe() {
	s.once.Do(func() {
		s.bus.mu.Lock()
		defer s.bus.mu.Unlock()

		close(s.bus.subscribers[s.id])
		delete(s.bus.subscribers, s.id)

		s.bus = nil
		s.notify = nil
	})
}

func (b *Bus[T]) Close() {
	b.once.Do(func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		for id, sub := range b.subscribers {
			close(sub)
			delete(b.subscribers, id)
		}

		b.buf.Clear()
		close(b.close)
	})
}

func (s *Subscriber[T]) Pop() (T, bool) {
	return s.bus.pop()
}

func (s *Subscriber[T]) Dequeue() ([]T, bool) {
	return s.bus.dequeue()
}

func (b *Bus[T]) Subscribe() *Subscriber[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	notify := make(chan struct{}, 1)
	if b.buf.Len() > 0 {
		notify <- struct{}{}
	}

	id := b.counter
	b.subscribers[id] = notify
	b.counter++

	return &Subscriber[T]{
		bus:    b,
		id:     id,
		notify: notify,
	}
}

func (b *Bus[T]) Push(msg T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf.Push(msg)

	for _, sub := range b.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}
}

func (b *Bus[T]) pop() (T, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Pop()
}

func (b *Bus[T]) Enqueue(data []T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, v := range data {
		b.buf.Push(v)
	}

	for _, sub := range b.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}
}

func (b *Bus[T]) dequeue() ([]T, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var data []T

	msg, ok := b.buf.Pop()
	for ok {
		data = append(data, msg)
		msg, ok = b.buf.Pop()
	}

	return data, len(data) > 0
}
