package actor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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
			replies:       make(chan *Reply, 1000),
		}
	}

	sys.registry["system"] = &supervisor{
		name:    "system",
		replies: make(chan *Reply, 1000),
	}

	return sys
}

type MessageContext struct {
	name   string
	system *System
}

func (m *MessageContext) Name() string {
	return m.name
}

func (m *MessageContext) Send(dest string, message *Message) {
	m.system.registry[dest].actor.Send(message)
}

func (m *MessageContext) ReplyTo(dest string, reply *Reply) {
	reply.receiver = dest
	dest = strings.Split(dest, ":")[0]
	for _, super := range m.system.registry {
		if super.name == dest {
			super.replies <- reply
			return
		}
	}
	m.system.registry[dest].replies <- reply
}

func (m *MessageContext) Reply() *Reply {
	dest := strings.Split(m.name, ":")[0]
	return <-m.system.registry[dest].replies
}

func (c *System) Context() *MessageContext {
	return &MessageContext{
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

	log.Println("application started")
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
		if super.name == "system" {
			continue
		}
		super.actor.Stop()
		super.wg.Wait()
	}
}
