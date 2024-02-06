package actor

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"time"
)

const (
	RestartPolicyAlways = 0
	RestartPolicyNever  = 1
	RestartPolicyCrash  = 2

	StatusAlive   = 0
	StatusStopped = 1
)

type IActor interface {
	Send(message *Message)
	Run(s *supervisor)
	Stop()
}

type State struct {
	state interface{}
}

type Handle func(*MessageContext, *State, *Message) error

type supervisor struct {
	sys           *System
	name          string
	restartPolicy int
	wg            sync.WaitGroup
	actor         IActor
	active        int
	replies       chan *Reply
	mu            sync.Mutex
	status        int
}

func (s *supervisor) panic(id int, req *Message) {
	s.status = StatusStopped
	if s.restartPolicy == RestartPolicyAlways {
		go s.actor.Run(s)
	}

	if s.restartPolicy == RestartPolicyCrash {
		panic(req)
	}

	if req != nil {
		if req.Attempts < 10 {
			time.Sleep(time.Duration(req.Attempts) * time.Second)
			s.actor.Send(req)
		} else {
			log.Printf("message failed (actor:%s) (retries:%d): %v %+v\n", s.name, req.Attempts, req.Data, req.errors)
		}
	}
}

type Message struct {
	Data     []byte
	Sender   string
	Attempts int
	errors   []any
}

func NewMessage[T any](v T) *Message {
	buf := bytes.Buffer{}

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(v)
	if err != nil {
		panic(err)
	}

	return &Message{
		Data: buf.Bytes(),
	}
}

func NewMessageWithReply[T any](sender string, v T) *Message {
	buf := bytes.Buffer{}

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(v)
	if err != nil {
		panic(err)
	}

	m := &Message{
		Sender:   sender,
		Data:     buf.Bytes(),
		Attempts: 0,
	}

	return m
}

type actorStop struct{}
