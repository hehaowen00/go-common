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
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "panic2",
		Actor:         panicActor,
		RestartPolicy: actor.RestartPolicyNever,
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
		s, _ := actor.GetState[int](state)

		val, err := actor.GetMessage[int](msg)
		if err != nil {
			panic(fmt.Errorf("invalid message %v", err))
		}

		*s = *s + (val * 2)

		sys.ReplyTo(msg.Sender, actor.NewReply[int]("system", *s))

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

		log.Println("actor ok")

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
			panic("invalid state")
		}

		s.Lock()
		defer s.Unlock()

		val, err := actor.GetMessage[int](msg)
		if err != nil {
			panic(fmt.Errorf("invalid message %v", err))
		}

		single := actor.NewMessageWithReply(sys.Name(), val)
		sys.Send("test", single)

		reply, err := actor.GetReply[int](sys.Reply())
		log.Println("reply:", reply, err)

		s.data = append(s.data, val)

		return nil
	}, 4)

func main() {
	system.Start()

	m := system.Context()

	conn1, ok := system.GetConn("test")
	if !ok {
		panic("invalid conn")
	}

	msg := actor.NewMessageWithReply(m.Name(), 1)
	conn1.Send(msg)

	reply, err := actor.GetReply[int](m.Reply())
	log.Println("reply:", reply, err)

	msg = actor.NewMessageWithReply(m.Name(), 2)
	conn1.Send(msg)

	reply, err = actor.GetReply[int](m.Reply())
	log.Println("reply:", reply, err)

	conn2, ok := system.GetConn("panic")
	if !ok {
		panic("invalid conn")
	}

	conn3, ok := system.GetConn("multi")
	if !ok {
		panic("invalid conn")
	}

	for i := 0; i < 10; i++ {
		msg := actor.NewMessage(i)
		conn3.Send(msg)
	}

	msg = actor.NewMessage[Test](Test{})
	conn2.Send(msg)

	system.Wait()
}
