package actor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

type System struct {
	registry map[string]*supervisor
	wg       sync.WaitGroup
}

func NewSystem(configs ...*Config) *System {
	ctx := &System{
		registry: make(map[string]*supervisor),
	}

	for _, config := range configs {
		ctx.registry[config.Name] = &supervisor{
			name:          config.Name,
			restartPolicy: config.RestartPolicy,
			actor:         config.Actor,
		}
	}

	return ctx
}

func (c *System) restart(name string) {
	log.Println("actor restart:", name)

	super, ok := c.registry[name]
	if !ok {
		return
	}
	go super.actor.Run(super)
}

func (c *System) Start() {
	for _, super := range c.registry {
		go super.actor.Run(super)
	}

	log.Println("application started")
}

func (c *System) GetConn(name string) (IActor, bool) {
	a, ok := c.registry[name]
	return a.actor, ok
}

func (c *System) Wait() {
	log.Println("waiting for interrupt")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Print("\r")
	log.Println("application stopped")

	for _, super := range c.registry {
		super.actor.Stop()
		super.wg.Wait()
	}
}
