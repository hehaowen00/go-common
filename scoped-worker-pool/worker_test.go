package scopedworkerpool_test

import (
	"errors"
	"testing"
	"time"

	scopedworkerpool "github.com/hehaowen00/go-common/scoped-worker-pool"
)

type Message struct {
	ID string
}

type Processor struct{}

func (p *Processor) GetScope(msg *Message) (string, error) {
	if msg == nil {
		return "", errors.New("no msg found")
	}

	return msg.ID, nil
}

func (p *Processor) Handle(msg *Message) error {
	if msg == nil {
		return errors.New("no msg found")
	}

	if len(msg.ID)%2 != 0 {
		return errors.New("odd length")
	}

	return nil
}

func TestWorker(t *testing.T) {
	sup := scopedworkerpool.NewSupervisor[Message](&Processor{}, 16, time.Second*5)

	sup.Push(&Message{
		ID: "a",
	})

	sup.Push(&Message{
		ID: "ab",
	})

	sup.Push(nil)

	time.Sleep(10 * time.Second)
	sup.Stop()
}
