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

	mu   sync.Mutex
	once sync.Once
}

func NewChannel[T any]() *Channel[T] {
	return &Channel[T]{
		buf:         ring.NewRing[T](),
		close:       make(chan struct{}),
		subscribers: make(map[int]chan struct{}),
	}
}

func (c *Channel[T]) Close() {
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		for id, sub := range c.subscribers {
			close(sub)
			delete(c.subscribers, id)
		}

		c.buf.Clear()
		c.buf = nil
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
	defer c.mu.Unlock()

	if c.buf == nil {
		return
	}

	c.buf.Push(msg)

	for _, sub := range c.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}
}

func (c *Channel[T]) pop() (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.buf == nil {
		var tmp T
		return tmp, false
	}

	return c.buf.Pop()
}

func (c *Channel[T]) Enqueue(data []T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.buf == nil {
		return
	}

	for _, v := range data {
		c.buf.Push(v)
	}

	for _, sub := range c.subscribers {
		if len(sub) == 0 {
			sub <- struct{}{}
		}
	}
}

func (c *Channel[T]) dequeue() []T {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.buf == nil {
		return nil
	}

	return c.buf.Dequeue()
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
		defer s.channel.mu.Unlock()

		close(s.channel.subscribers[s.id])
		delete(s.channel.subscribers, s.id)

		s.channel = nil
		s.notify = nil
	})
}

func (s *Subscriber[T]) Pop() (T, bool) {
	if s.channel == nil {
		var tmp T
		return tmp, false
	}

	return s.channel.pop()
}

func (s *Subscriber[T]) Dequeue() []T {
	if s.channel == nil {
		return nil
	}

	return s.channel.dequeue()
}
