package actor

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	RestartPolicyAlways = 0
	RestartPolicyNever  = 1

	StatusAlive   = 0
	StatusStopped = 1
)

var (
	ErrDidNotReply = errors.New("did not reply")
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
	mu            sync.Mutex
	status        int
}

func (s *supervisor) panic(id int, req *Message) {
	s.status = StatusStopped
	if s.restartPolicy == RestartPolicyAlways {
		go s.actor.Run(s)
	}

	if s.restartPolicy == RestartPolicyNever {
		log.Printf("[err] [actor:%s] message failed: %v %+v\n", s.name, req.Data, req.error)
	}

	if req != nil {
		if req.Attempts < 10 {
			time.Sleep(time.Duration(req.Attempts) * time.Second)
			s.actor.Send(req)
		} else {
			log.Printf("message failed (actor:%s) (retries:%d): %v %+v\n", s.name, req.Attempts, req.Data, req.error)
			req.ReplyTo(NewReplyWithError(fmt.Errorf("message failed: %+v", req.error)))
		}
	}
}

type Message struct {
	Data     []byte
	Sender   string
	Attempts int
	error    any
	reply    chan *Reply
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

type Resp struct {
	reply chan *Reply
}

func (r *Resp) resp() *Reply {
	resp, ok := <-r.reply
	if !ok {
		return &Reply{
			err: fmt.Errorf("reply closed"),
		}
	}
	return resp
}

func NewMessageWithReply[T any](v T) (*Message, *Resp, error) {
	buf := bytes.Buffer{}

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(v)
	if err != nil {
		return nil, nil, err
	}

	m := &Message{
		Data:     buf.Bytes(),
		Attempts: 0,
		reply:    make(chan *Reply, 1),
	}

	resp := &Resp{
		reply: m.reply,
	}

	return m, resp, nil
}

func (m *Message) WantsReply() bool {
	return m.reply != nil
}

func (m *Message) ReplyTo(reply *Reply) {
	if m.reply == nil {
		return
	}

	if reply == nil {
		m.reply <- &Reply{
			err: ErrDidNotReply,
		}
		close(m.reply)
		return
	}

	m.reply <- reply
}
