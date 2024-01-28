package broadcast

import (
	"sync"

	"github.com/hehaowen00/go-common/ring"
)

type Broadcast[T any] struct {
	counter     int
	subscribers map[int]*Subscriber[T]

	mu   sync.Mutex
	once sync.Once
}

func NewBroadcast[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		subscribers: make(map[int]*Subscriber[T]),
	}
}

func (b *Broadcast[T]) Close() {
	b.once.Do(func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		for id, s := range b.subscribers {
			s.mu.Lock()
			s.buf.Clear()
			s.mu.Unlock()

			close(s.close)
			delete(b.subscribers, id)
		}
	})
}

func (b *Broadcast[T]) Subscribe() *Subscriber[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	s := Subscriber[T]{
		broadcast: b,
		buf:       ring.NewRing[T](),
		close:     make(chan struct{}),
		notify:    make(chan struct{}, 1),
	}

	b.subscribers[b.counter] = &s
	b.counter++

	return &s
}

func (b *Broadcast[T]) Push(item T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, s := range b.subscribers {
		s.mu.Lock()
		s.buf.Push(item)
		s.mu.Unlock()

		if len(s.notify) == 0 {
			s.notify <- struct{}{}
		}
	}
}

func (b *Broadcast[T]) Enqueue(data []T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, s := range b.subscribers {
		s.mu.Lock()

		for _, v := range data {
			s.buf.Push(v)
		}

		s.mu.Unlock()

		if len(s.notify) == 0 {
			s.notify <- struct{}{}
		}
	}
}

type Subscriber[T any] struct {
	broadcast *Broadcast[T]
	buf       *ring.Ring[T]
	close     chan struct{}
	notify    chan struct{}
	mu        sync.Mutex
	once      sync.Once
}

func (s *Subscriber[T]) Closed() <-chan struct{} {
	return s.close
}

func (s *Subscriber[T]) Notify() <-chan struct{} {
	return s.notify
}

func (s *Subscriber[T]) Unsubscribe() {
	s.once.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		close(s.notify)
		close(s.close)

		s.broadcast = nil
		s.buf.Clear()
		s.buf = nil
		s.notify = nil
	})
}

func (s *Subscriber[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.buf.Pop()
}

func (s *Subscriber[T]) Dequeue() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.buf.Dequeue()
}
