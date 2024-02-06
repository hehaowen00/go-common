package actor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
)

type Actor struct {
	queue   chan interface{}
	replies chan *Reply
	state   *State
	handle  Handle
	once    sync.Once
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
		queue:  make(chan interface{}, 100),
		handle: handle,
	}
}

func (a *Actor) Send(message *Message) {
	a.queue <- message
}

func (a *Actor) Handle(s *supervisor, message interface{}) error {
	return a.handle(nil, a.state, message.(*Message))
}

func (a *Actor) Run(s *supervisor) {
	s.wg.Add(1)
	name := s.name
	s.status = StatusAlive

	log.Printf("actor start: %s\n", name)

	var req *Message

	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			req.errors = append(req.errors, r)
			s.panic(0, req)
		}
	}()

	sys := &MessageContext{
		name:   name,
		system: s.sys,
	}

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
				err := a.handle(sys, a.state, msg.(*Message))
				if err != nil {
					log.Println("actor error:", name, err)
				}
			default:
				panic("invalid message")
			}
		}

		req = nil
	}
}

func (a *Actor) Stop() {
	a.once.Do(func() {
		close(a.queue)
	})
}

type Context struct {
}

type Reply struct {
	receiver string
	data     []byte
}

func (r *Reply) Receiver() string {
	return r.receiver
}

func (r *Reply) Data() []byte {
	return r.data
}

func NewReply[T any](receiver string, data T) *Reply {
	buf := bytes.Buffer{}

	enc := gob.NewEncoder(&buf)

	enc.Encode(data)

	return &Reply{
		receiver: receiver,
		data:     buf.Bytes(),
	}
}

func GetReply[T any](r *Reply) (T, error) {
	if r == nil {
		var temp T
		return temp, fmt.Errorf("reply is nil")
	}

	buf := bytes.NewBuffer(r.data)

	dec := gob.NewDecoder(buf)

	var v T

	err := dec.Decode(&v)

	return v, err

}
