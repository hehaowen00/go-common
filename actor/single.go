package actor

import "log"

type Actor struct {
	queue   chan interface{}
	replies chan *Reply
	state   *State
	handle  Handle
}

type Config struct {
	Name          string
	Actor         IActor
	Max           int
	RestartPolicy int
	Retries       int
}

func NewConfig(actor IActor, restartPolicy int) *Config {
	return &Config{
		Actor:         actor,
		RestartPolicy: restartPolicy,
	}
}

func NewActor[S any](state S, handle Handle) *Actor {
	return &Actor{
		state: &State{
			state: &state,
		},
		queue:  make(chan interface{}, 5),
		handle: handle,
	}
}

func (a *Actor) Send(message *Message) {
	a.queue <- message
}

func (a *Actor) Handle(s *supervisor, message interface{}) error {
	return a.handle(a.state, message.(*Message))
}

func (a *Actor) Run(s *supervisor) {
	s.wg.Add(1)
	name := s.name

	var req *Message

	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			s.panic(0, req)
		}
	}()

	for {
		select {
		case msg, ok := <-a.queue:
			if !ok {
				log.Println("actor terminated:", name)
				return
			}

			switch msg.(type) {
			case *Message:
				req = msg.(*Message)
				req.Attempts++
				a.handle(a.state, msg.(*Message))
			case *Reply:
				rep := msg.(*Reply)
				a.queue <- rep
			default:
				panic("invalid message")
			}
		}

		req = nil
	}
}

func (a *Actor) Stop() {
	close(a.queue)
}

type Context struct {
}

type Reply struct {
	Receiver string
	Data     interface{}
}
