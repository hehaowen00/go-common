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

func (r *Broadcast) Send(message *Message) {
	r.queue <- message
}

func (r *Broadcast) Stop() {
	close(r.queue)
}

func (r *Broadcast) Run(s *supervisor) {
	s.wg.Add(1)

	log.Printf("[info] [actor:%s] start\n", s.name)

	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			s.panic(0, nil)
		}
	}()

	sys := &MessageContext{
		name:   "replicated",
		system: s.sys,
	}

	for {
		select {
		case msg, ok := <-r.queue:
			if !ok {
				return
			}

			switch msg := msg.(type) {
			case *Message:
				for _, actor := range r.actors {
					go sys.Send(actor, msg)
				}
			default:
				panic("invalid message")
			}
		}
	}
}
