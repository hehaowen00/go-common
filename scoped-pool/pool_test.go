package scopedpool_test

import (
	"errors"
	"testing"
	"time"

	scopedpool "github.com/hehaowen00/go-common/scoped-pool"
)

type message struct {
	ID string
}

type processor struct{}

func (p *processor) GetScope(msg *message) (string, error) {
	if msg == nil {
		return "", errors.New("no msg found")
	}

	return msg.ID, nil
}

func (p *processor) Handle(msg *message) error {
	if msg == nil {
		return errors.New("no msg found")
	}

	if len(msg.ID)%2 != 0 {
		return errors.New("odd length")
	}

	return nil
}

func TestWorker(t *testing.T) {
	pool := scopedpool.NewPool[message](&processor{}, 16, time.Second*5)

	pool.Push(&message{
		ID: "a",
	})

	pool.Push(&message{
		ID: "ab",
	})

	pool.Push(nil)

	time.Sleep(10 * time.Second)
	pool.Stop()
}
