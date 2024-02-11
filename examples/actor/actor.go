package main

import (
	"log"
	"sync"

	"github.com/hehaowen00/go-common/actor"
)

type Test struct {
}

var system = actor.NewSystem(
	&actor.Config{
		Name:          "test",
		Actor:         actor.NewActor(&Actor1{}),
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "panic",
		Actor:         actor.NewActor(&PanicActor{}),
		RestartPolicy: actor.RestartPolicyNever,
	},
	&actor.Config{
		Name:          "panic2",
		Actor:         actor.NewActor(&PanicActor{}),
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "multi",
		Actor:         actor.NewScalar(&ScalarActor{}, 0),
		RestartPolicy: actor.RestartPolicyAlways,
	},
)

type Actor1 struct {
	data int
}

func (t *Actor1) Handle(ctx *actor.Context, msg *actor.Message) error {
	val, err := actor.GetMessage[int](msg)
	if err != nil {
		return err
	}

	t.data += val

	reply, _ := actor.NewReply[int](t.data)
	msg.ReplyTo(reply)

	return nil
}

type PanicActor struct {
	count int
}

func (p *PanicActor) Handle(ctx *actor.Context, msg *actor.Message) error {
	if p.count < 10 {
		p.count++
		panic("panic")
	}

	ctx.Info("actor ok")

	return nil
}

type ScalarActor struct {
	data []int
	sync.Mutex
}

func (s *ScalarActor) Handle(ctx *actor.Context, msg *actor.Message) error {
	s.Lock()
	defer s.Unlock()

	val, err := actor.GetMessage[int](msg)
	if err != nil {
		return err
	}

	single := actor.NewMessage(val)
	ctx.Send("test", single)

	reply, _ := actor.NewReply(0)
	msg.ReplyTo(reply)

	s.data = append(s.data, val)

	return nil
}

func main() {
	system.Start()

	ctx := system.Context()

	for i := 0; i < 10; i++ {
		msg, resp, _ := actor.NewMessageWithReply(i)
		ctx.Send("multi", msg)
		actor.GetReply[int](resp)
	}

	msg, resp, _ := actor.NewMessageWithReply(1)
	ctx.Send("test", msg)

	reply, err := actor.GetReply[int](resp)
	log.Println("reply:", 1, reply, err)

	msg, resp, _ = actor.NewMessageWithReply(2)
	ctx.Send("test", msg)

	reply, err = actor.GetReply[int](resp)
	log.Println("reply:", 2, reply, err)

	msg = actor.NewMessage[Test](Test{})
	ctx.Send("panic", msg)

	msg = actor.NewMessage[Test](Test{})
	ctx.Send("panic2", msg)

	system.Wait()
}
