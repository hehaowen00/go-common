package actor

import (
	"fmt"
	"log"
	"sync"
)

type Scalar struct {
	count  int
	queue  chan interface{}
	handle Handle
	wg     sync.WaitGroup
	once   sync.Once
}

func NewScalar(handle Handle, count int) *Scalar {
	return &Scalar{
		queue:  make(chan interface{}, 5),
		handle: handle,
		count:  count,
	}
}

func (a *Scalar) Run(s *supervisor) {
	name := s.name

	var req *Message

	if a.count == 0 {
		s.wg.Add(1)
		defer func() {
			s.wg.Done()
		}()

		ctx := &Context{
			name:   name,
			system: s.sys,
		}

		for {
			select {
			case msg, ok := <-a.queue:
				if !ok {
					return
				}

				switch msg.(type) {
				case *Message:
					req = msg.(*Message)
					req.Attempts++

					a.wg.Add(1)
					go func(msg *Message) {
						defer func() {
							a.wg.Done()
							if r := recover(); r != nil {
								msg.error = r
								s.panic(0, msg)
							}
						}()

						err := a.handle.Handle(ctx, msg)
						if err != nil {
							log.Printf("[info] [actor:%s] error: %v\n", name, err)
						}
					}(req)
				}
			}
		}
	}

	for i := 0; i < a.count; i++ {
		log.Printf("[info] [system] actor start: %s:%d\n", name, i)
		a.wg.Add(1)

		var req *Message

		ctx := &Context{
			name:   fmt.Sprintf("%s:%d", name, i),
			system: s.sys,
		}

		go func(id int) {
			defer func() {
				a.wg.Done()
				if r := recover(); r != nil {
					req.error = r
					s.panic(id, req)
				}
			}()

			for {
				select {
				case msg, ok := <-a.queue:
					if !ok {
						log.Printf("[info] [actor:%s:%d] terminated\n", name, id)
						return
					}

					switch msg.(type) {
					case *Message:
						req = msg.(*Message)
						req.Attempts++
						err := a.handle.Handle(ctx, req)
						if err != nil {
							log.Printf("[info] [actor:%s] error: %v\n", name, err)
						}
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
