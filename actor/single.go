package actor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
)

type Actor struct {
	queue  chan interface{}
	handle Handle
	once   sync.Once
}

type Config struct {
	Name          string
	Actor         actorInterface
	Max           int
	RestartPolicy int
	Retries       int
}

func NewActor(handle Handle) *Actor {
	return &Actor{
		queue:  make(chan interface{}, 100),
		handle: handle,
	}
}

func (a *Actor) Send(message *Message) {
	a.queue <- message
}

func (a *Actor) Run(s *supervisor) {
	s.wg.Add(1)
	name := s.name
	s.status = StatusAlive

	log.Printf("[info] [actor:%s] start\n", name)

	var req *Message

	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			req.error = r
			s.panic(0, req)
		}
	}()

	sys := &Context{
		name:   name,
		system: s.sys,
	}

	for {
		select {
		case msg, ok := <-a.queue:
			if !ok {
				log.Printf("[info] [actor:%s] terminated\n", name)
				return
			}

			switch msg.(type) {
			case *Message:
				req = msg.(*Message)
				req.Attempts++
				err := a.handle.Handle(sys, msg.(*Message))
				if err != nil {
					log.Println("[info] [system] actor error:", name, err)
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

type Reply struct {
	data []byte
	err  error
}

func (r *Reply) Data() []byte {
	return r.data
}

func NewReply[T any](data T) (*Reply, error) {
	buf := bytes.Buffer{}

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	reply := &Reply{
		data: buf.Bytes(),
	}

	return reply, nil
}

func NewReplyWithError(err error) *Reply {
	return &Reply{
		err: err,
	}
}

func GetReply[T any](resp *Resp) (T, error) {
	r := resp.resp()

	if r.err != nil {
		var tmp T
		return tmp, r.err
	}

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
