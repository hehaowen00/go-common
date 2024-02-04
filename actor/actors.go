package actor

import (
	"log"
	"sync"
)

const (
	RestartPolicyAlways = 0
	RestartPolicyNever  = 1
)

type IActor interface {
	Send(message *Message)
	Run(s *supervisor)
	Handle(s *supervisor, message interface{}) error
	Stop()
}

type State struct {
	state interface{}
}

type Handle func(*State, *Message) error

type supervisor struct {
	name          string
	restartPolicy int
	wg            sync.WaitGroup
	actor         IActor
	active        int
	mu            sync.Mutex
}

func (s *supervisor) panic(id int, req *Message) {
	if s.restartPolicy == RestartPolicyAlways {
		go s.actor.Run(s)
		if req != nil {
			if req.Attempts < 10 {
				s.actor.Send(req)
			} else {
				log.Printf("message failed (actor:%s) (retries:%d): %v\n", s.name, req.Attempts, req.Data)
			}
		}
	}
}

type Message struct {
	Data     interface{}
	Reply    chan interface{}
	Attempts int
}

func NewMessage(v interface{}) *Message {
	return &Message{
		Data: v,
	}
}

func NewMessageWithReply(sender string, v interface{}) (*Message, <-chan interface{}) {
	m := &Message{
		Data:     v,
		Reply:    make(chan interface{}),
		Attempts: 0,
	}

	return m, m.Reply
}

type actorStop struct{}
