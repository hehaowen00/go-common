package actor

import (
	"log"
	"sync"
)

type Scalar struct {
	count  int
	queue  chan interface{}
	state  *State
	handle Handle
	wg     sync.WaitGroup
	once   sync.Once
}

func NewScalar[T any](state T, handle Handle, count int) *Scalar {
	return &Scalar{
		state: &State{
			state: &state,
		},
		queue:  make(chan interface{}, 5),
		handle: handle,
		count:  count,
	}
}

func (a *Scalar) Handle(s *supervisor, message interface{}) error {
	return a.handle(nil, a.state, message.(*Message))
}

func (a *Scalar) Run(s *supervisor) {
	name := s.name

	for i := 0; i < a.count; i++ {
		log.Printf("actor start: %s:%d\n", name, i)
		a.wg.Add(1)

		var req *Message

		sys := &MessageContext{
			name:   name,
			system: s.sys,
		}

		go func(id int) {
			defer func() {
				a.wg.Done()
				if r := recover(); r != nil {
					req.errors = append(req.errors, r)
					s.panic(id, req)
				}
			}()

			for {
				select {
				case msg, ok := <-a.queue:
					if !ok {
						log.Printf("actor terminated: %s:%d\n", name, id)
						return
					}

					switch msg.(type) {
					case *Message:
						req = msg.(*Message)
						a.handle(sys, a.state, req)
						req.Attempts++
					default:
						panic("invalid message")
					}
				}

				req = nil
			}
		}(i)
	}
}

func (a *Scalar) Send(message *Message) {
	a.queue <- message
}

func (a *Scalar) Stop() {
	a.once.Do(func() {
		close(a.queue)
	})
	a.wg.Wait()
}
