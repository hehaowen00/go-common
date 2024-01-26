package common

import "sync"

type Bus[T any] struct {
	mu          sync.Mutex
	data        []T
	counter     int
	close       chan struct{}
	subscribers map[int]chan struct{}
}

func NewBus[T any]() *Bus[T] {
	return &Bus[T]{
		close:       make(chan struct{}),
		subscribers: make(map[int]chan struct{}),
	}
}

type Subscriber[T any] struct {
	bus    *Bus[T]
	id     int
	Notify chan struct{}
}

func (s *Subscriber[T]) Close() {
	s.bus.Unsubscribe(s.id)
}

func (s *Subscriber[T]) Closed() <-chan struct{} {
	return s.bus.close
}

func (b *Bus[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for id, sub := range b.subscribers {
		close(sub)
		delete(b.subscribers, id)
	}

	close(b.close)
}

func (b *Bus[T]) Subscribe() *Subscriber[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	notify := make(chan struct{}, 1)
	if len(b.data) > 0 {
		notify <- struct{}{}
	}

	id := b.counter
	b.subscribers[id] = notify
	b.counter++

	return &Subscriber[T]{b, id, notify}
}

func (b *Bus[T]) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	close(b.subscribers[id])
	delete(b.subscribers, id)
}

func (b *Bus[T]) Push(msg T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = append(b.data, msg)

	wg := sync.WaitGroup{}

	for _, sub := range b.subscribers {
		wg.Add(1)
		go func(sub chan struct{}) {
			if len(sub) == 0 {
				sub <- struct{}{}
			}
			wg.Done()
		}(sub)
	}

	wg.Wait()
}

func (b *Bus[T]) Pop() (T, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.data) == 0 {
		return *new(T), false
	}

	msg := b.data[0]
	b.data = b.data[1:]

	return msg, true
}

func (b *Bus[T]) Enqueue(data []T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = append(b.data, data...)

	wg := sync.WaitGroup{}

	for _, sub := range b.subscribers {
		wg.Add(1)

		go func(sub chan struct{}) {
			if len(sub) == 0 {
				sub <- struct{}{}
			}
			wg.Done()
		}(sub)
	}

	wg.Wait()
}

func (b *Bus[T]) Dequeue() ([]T, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.data) == 0 {
		return nil, false
	}

	msg := b.data
	b.data = nil

	return msg, true
}
