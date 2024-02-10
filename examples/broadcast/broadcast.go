package main

import (
	"log"

	"github.com/hehaowen00/go-common/actor"
)

var system = actor.NewSystem(
	&actor.Config{
		Name:          "broadcast",
		Actor:         broadcast,
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "actor1",
		Actor:         actor1,
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "actor2",
		Actor:         actor2,
		RestartPolicy: actor.RestartPolicyAlways,
	},
)

var broadcast = actor.NewBroadcast("actor1", "actor2")

var actor1 = actor.NewActor(
	0,
	func(sys *actor.MessageContext, state *actor.State, msg *actor.Message) error {
		s, _ := actor.GetState[int](state)

		val, err := actor.GetMessage[int](msg)
		if err != nil {
			return err
		}

		*s = *s + val

		log.Println("actor1", *s)

		return nil
	},
)

var actor2 = actor.NewActor(
	0,
	func(sys *actor.MessageContext, state *actor.State, msg *actor.Message) error {
		s, _ := actor.GetState[int](state)
		*s = *s + 1
		log.Println("actor2", *s)
		return nil
	},
)

func main() {
	system.Start()

	ctx := system.Context()

	for i := 0; i < 10; i++ {
		ctx.Send("replicated", actor.NewMessage[int](i))
	}

	system.Wait()
}
