package main

import (
	"log"
	"sync"

	"github.com/hehaowen00/go-common/actor"
)

type TestActor struct {
	data []int
	sync.Mutex
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
		Name:          "multi",
		Actor:         multiActor,
		RestartPolicy: actor.RestartPolicyAlways,
	},
)

var testActor = actor.NewActor(
	0,
	func(state *actor.State, msg *actor.Message) error {
		s, _ := actor.GetState[int](state)

		val, ok := actor.GetMessage[int](msg)
		if !ok {
			panic("invalid message")
		}

		*s = *s + (val * 2)

		msg.Reply <- *s

		return nil
	})

var panicActor = actor.NewActor(
	0,
	func(state *actor.State, msg *actor.Message) error {
		panic("panic")
		return nil
	})

var multiActor = actor.NewMultiProcess(
	TestActor{},
	func(state *actor.State, msg *actor.Message) error {
		s, ok := actor.GetState[TestActor](state)
		if !ok {
			panic("invalid state")
		}

		s.Lock()
		defer s.Unlock()

		val, ok := actor.GetMessage[int](msg)
		if !ok {
			panic("invalid message")
		}

		s.data = append(s.data, val)

		log.Println("data", s.data)

		return nil
	}, 4)

func main() {
	system.Start()

	conn1, ok := system.GetConn("test")
	if !ok {
		panic("invalid conn")
	}

	msg, reply := actor.NewMessageWithReply(1)
	conn1.Send(msg)
	log.Println(<-reply)

	msg, reply = actor.NewMessageWithReply(2)
	conn1.Send(msg)
	log.Println(<-reply)

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

	msg = actor.NewMessage(nil)
	conn2.Send(msg)

	system.Wait()
}
