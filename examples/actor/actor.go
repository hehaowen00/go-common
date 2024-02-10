package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/hehaowen00/go-common/actor"
)

type Test struct {
}

var system = actor.NewSystem(
	&actor.Config{
		Name:          "test",
		Actor:         testActor,
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "panic",
		Actor:         panicActor,
		RestartPolicy: actor.RestartPolicyNever,
	},
	&actor.Config{
		Name:          "panic2",
		Actor:         panicActor,
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "multi",
		Actor:         processActor,
		RestartPolicy: actor.RestartPolicyAlways,
	},
)

var testActor = actor.NewActor(
	0,
	func(sys *actor.MessageContext, state *actor.State, msg *actor.Message) error {
		s, ok := actor.GetState[int](state)
		if !ok {
			return fmt.Errorf("invalid state")
		}

		val, err := actor.GetMessage[int](msg)
		if err != nil {
			return err
		}

		*s = *s + val

		reply, _ := actor.NewReply[int](*s)
		msg.ReplyTo(reply)

		return nil
	})

var panicActor = actor.NewActor(
	0,
	func(sys *actor.MessageContext, state *actor.State, msg *actor.Message) error {
		s, ok := actor.GetState[int](state)
		if !ok {
			return fmt.Errorf("invalid state")
		}

		if *s < 5 {
			*s++
			panic("panic")
		}

		sys.Info("actor ok")

		return nil
	})

type TestActor struct {
	data []int
	sync.Mutex
}

var processActor = actor.NewScalar(
	TestActor{},
	func(sys *actor.MessageContext, state *actor.State, msg *actor.Message) error {
		s, ok := actor.GetState[TestActor](state)
		if !ok {
			return fmt.Errorf("invalid state")
		}

		s.Lock()
		defer s.Unlock()

		val, err := actor.GetMessage[int](msg)
		if err != nil {
			return err
		}

		single := actor.NewMessage(val)
		sys.Send("test", single)

		reply, _ := actor.NewReply(0)
		msg.ReplyTo(reply)

		s.data = append(s.data, val)

		return nil
	}, 0)

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

	system.Wait()
}
