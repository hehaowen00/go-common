package scopedpool

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Handler[T any] interface {
	GetScope(*T) (string, error)
	Handle(*T) error
}

type Pool[T any] struct {
	handler    Handler[T]
	bufferSize int
	timeout    time.Duration
	workers    map[string]*worker[T]
	wg         sync.WaitGroup

	alive atomic.Bool
	once  sync.Once
	mu    sync.Mutex
}

func NewPool[T any](handler Handler[T], bufferSize int, timeout time.Duration) *Pool[T] {
	sup := Pool[T]{
		handler:    handler,
		bufferSize: bufferSize,
		timeout:    timeout,
		workers:    make(map[string]*worker[T]),
	}

	sup.alive.Store(true)

	return &sup
}

func (sup *Pool[T]) Push(message *T) {
	if message == nil {
		return
	}

	if !sup.alive.Load() {
		return
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()

	scope, err := sup.handler.GetScope(message)
	if err != nil {
		log.Println("[supervisor] error getting message id", err)
		return
	}

	w, ok := sup.workers[scope]
	if !ok {
		sup.wg.Add(1)
		w = &worker[T]{
			id: scope,
			f:  sup.handler.Handle,

			onClose: sup.clearWorker,
			in:      make(chan *T, sup.bufferSize),
			timeout: sup.timeout,
			wg:      &sup.wg,
		}

		go w.Run()
		sup.workers[scope] = w
	}

	w.in <- message
}

func (sup *Pool[T]) clearWorker(id string) {
	sup.mu.Lock()
	defer sup.mu.Unlock()

	delete(sup.workers, id)
}

func (sup *Pool[T]) Stop() {
	sup.once.Do(func() {
		sup.alive.Store(false)
		for _, w := range sup.workers {
			close(w.in)
		}
	})

	sup.wg.Wait()
}

type worker[T any] struct {
	id string
	f  func(*T) error

	onClose func(string)
	in      chan *T
	timeout time.Duration
	wg      *sync.WaitGroup
}

func (w *worker[T]) Run() {
	defer w.wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)

	for {
		select {
		case <-ctx.Done():
			log.Println("[supervisor] worker", w.id, "stopping")
			cancel()
			w.onClose(w.id)
			return
		case msg, ok := <-w.in:
			if !ok {
				cancel()
				return
			}

			ctx, cancel = context.WithTimeout(context.Background(), w.timeout)
			var err error

			for i := 0; i < 10; i++ {
				err = w.f(msg)
				if err == nil {
					break
				}
			}

			if err != nil {
				log.Printf("[supervisor] error processing message (%s) - %v\n", w.id, err)
			}
		}
	}
}
