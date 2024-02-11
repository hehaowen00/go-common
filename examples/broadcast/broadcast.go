package main

import (
	"log"

	"github.com/hehaowen00/go-common/actor"
)

var system = actor.NewSystem(
	&actor.Config{
		Name:          "broadcast",
		Actor:         actor.NewBroadcast("actor1", "actor2"),
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "actor1",
		Actor:         actor.NewActor(&Actor1{}),
		RestartPolicy: actor.RestartPolicyAlways,
	},
	&actor.Config{
		Name:          "actor2",
		Actor:         actor.NewActor(&Actor2{}),
		RestartPolicy: actor.RestartPolicyAlways,
	},
)

type Actor1 struct {
	sum int
}

func (a *Actor1) Handle(ctx *actor.Context, msg *actor.Message) error {
	val, err := actor.GetMessage[int](msg)
	if err != nil {
		return err
	}

	a.sum += val
	log.Println("actor1", a.sum)

	return nil
}

type Actor2 struct {
	count int
}

func (a *Actor2) Handle(ctx *actor.Context, msg *actor.Message) error {
	_, err := actor.GetMessage[int](msg)
	if err != nil {
		return err
	}

	a.count++
	log.Println("actor2", a.count)

	return nil
}

func main() {
	system.Start()

	ctx := system.Context()

	for i := 0; i < 10; i++ {
		ctx.Send("broadcast", actor.NewMessage[int](i))
	}

	system.Wait()
}
