package actor

import (
	"log"
	"sync"
)

type MultiProcess struct {
	count  int
	queue  chan interface{}
	state  *State
	handle Handle
	wg     sync.WaitGroup
}

func NewMultiProcess[T any](state T, handle Handle, count int) *MultiProcess {
	return &MultiProcess{
		state: &State{
			state: &state,
		},
		queue:  make(chan interface{}, 5),
		handle: handle,
		count:  count,
	}
}

func (a *MultiProcess) Handle(s *supervisor, message interface{}) error {
	return a.handle(a.state, message.(*Message))
}

func (a *MultiProcess) Run(s *supervisor) {
	name := s.name

	for i := 0; i < a.count; i++ {
		log.Printf("actor start: %s:%d\n", name, i)
		a.wg.Add(1)
		var req *Message

		go func(id int) {
			defer func() {
				a.wg.Done()
				if r := recover(); r != nil {
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
						a.handle(a.state, req)
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

func (a *MultiProcess) Send(message *Message) {
	a.queue <- message
}

func (a *MultiProcess) Stop() {
	close(a.queue)
	a.wg.Wait()
}
