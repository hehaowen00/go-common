package channel

import (
	"sync"

	"github.com/hehaowen00/go-common/ring"
)

type Channel[T any] struct {
	buf         *ring.Ring[T]
	counter     int
	close       chan struct{}
	subscribers map[int]chan struct{}

	mu   *sync.Mutex
	once sync.Once
}

func NewChannel[T any]() *Channel[T] {
	mu := &sync.Mutex{}

	return &Channel[T]{
		buf:         ring.NewRing[T](),
		close:       make(chan struct{}),
		subscribers: make(map[int]chan struct{}),
		mu:          mu,
	}
}

func (c *Channel[T]) Close() {
	c.once.Do(func() {
		c.mu.Lock()

		for id, sub := range c.subscribers {
			close(sub)
			delete(c.subscribers, id)
		}

		c.mu.Unlock()

		close(c.close)
	})
}

func (c *Channel[T]) Subscribe() *Subscriber[T] {
	c.mu.Lock()
	defer c.mu.Unlock()

	notify := make(chan struct{}, 1)
	if c.buf.Len() > 0 {
		notify <- struct{}{}
	}

	id := c.counter
	c.subscribers[id] = notify
	c.counter++

	return &Subscriber[T]{
		channel: c,
		id:      id,
		notify:  notify,
	}
}

func (c *Channel[T]) Push(msg T) {
	c.mu.Lock()

	if c.buf == nil {
		c.mu.Unlock()
		return
	}

	c.buf.Push(msg)

	c.mu.Unlock()

	for _, sub := range c.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}
}

func (c *Channel[T]) pop() (T, bool) {
	c.mu.Lock()

	if c.buf == nil {
		var tmp T
		c.mu.Unlock()
		return tmp, false
	}

	v, ok := c.buf.Pop()
	c.mu.Unlock()

	return v, ok
}

func (c *Channel[T]) Enqueue(data []T) {
	c.mu.Lock()

	if c.buf == nil {
		c.mu.Unlock()
		return
	}

	c.buf.Enqueue(data)

	for _, sub := range c.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}

	c.mu.Unlock()
}

func (c *Channel[T]) dequeue() []T {
	c.mu.Lock()

	if c.buf == nil {
		c.mu.Unlock()
		return nil
	}

	v := c.buf.Dequeue()
	c.mu.Unlock()

	return v
}

type Subscriber[T any] struct {
	channel *Channel[T]
	id      int
	notify  chan struct{}
	once    sync.Once
}

func (s *Subscriber[T]) Closed() <-chan struct{} {
	return s.channel.close
}

func (s *Subscriber[T]) Notify() <-chan struct{} {
	return s.notify
}

func (s *Subscriber[T]) Unsubscribe() {
	s.once.Do(func() {
		s.channel.mu.Lock()

		close(s.channel.subscribers[s.id])
		delete(s.channel.subscribers, s.id)

		s.channel = nil
		s.notify = nil

		s.channel.mu.Unlock()
	})
}

func (s *Subscriber[T]) Pop() (T, bool) {
	return s.channel.pop()
}

func (s *Subscriber[T]) Dequeue() []T {
	if s.channel == nil {
		return nil
	}

	return s.channel.dequeue()
}
