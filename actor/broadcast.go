package actor

import "log"

type Broadcast struct {
	queue  chan interface{}
	actors []string
}

func NewBroadcast(actors ...string) *Broadcast {
	broadcast := &Broadcast{
		actors: actors,
		queue:  make(chan interface{}, 100),
	}

	return broadcast
}

func (b *Broadcast) Send(message *Message) {
	b.queue <- message
}

func (b *Broadcast) Stop() {
	close(b.queue)
}

func (b *Broadcast) Run(s *supervisor) {
	s.wg.Add(1)

	log.Printf("[info] [actor:%s] start\n", s.name)

	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			s.panic(0, nil)
		}
	}()

	ctx := &Context{
		name:   s.name,
		system: s.sys,
	}

	for {
		select {
		case msg, ok := <-b.queue:
			if !ok {
				return
			}

			switch msg := msg.(type) {
			case *Message:
				for _, actor := range b.actors {
					go ctx.Send(actor, msg)
				}
			default:
				panic("invalid message")
			}
		}
	}
}
