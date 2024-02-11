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
	sys := &System{
		registry: make(map[string]*supervisor),
	}

	for _, config := range configs {
		sys.registry[config.Name] = &supervisor{
			sys:           sys,
			name:          config.Name,
			restartPolicy: config.RestartPolicy,
			actor:         config.Actor,
		}
	}

	sys.registry["system"] = &supervisor{
		name: "system",
	}

	return sys
}

type Context struct {
	name   string
	system *System
}

func (m *Context) Name() string {
	return m.name
}

func (m *Context) Send(dest string, message *Message) {
	m.system.registry[dest].actor.Send(message)
}

func (m *Context) Info(format string, v ...interface{}) {
	v = append([]interface{}{m.name}, v...)
	log.Printf("[info] [actor:%s] "+format, v...)
}

func (c *System) Context() *Context {
	return &Context{
		name:   "system",
		system: c,
	}
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
		if super.name == "system" {
			continue
		}
		go super.actor.Run(super)
	}

	log.Println("[info] [system] application started")
}

func (c *System) Stop() {
	for _, super := range c.registry {
		if super.name == "system" {
			continue
		}

		super.actor.Stop()
		super.wg.Wait()
	}
}

func (c *System) GetConn(name string) (actorInterface, bool) {
	a, ok := c.registry[name]
	return a.actor, ok
}

func (c *System) Wait() {
	log.Println("[info] [system] waiting for interrupt")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Print("\r")
	log.Println("[info] [system] application stopped")

	for _, super := range c.registry {
		if super.name == "system" {
			continue
		}
		super.actor.Stop()
		super.wg.Wait()
	}
}
