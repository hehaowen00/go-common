package broadcast_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hehaowen00/go-common/broadcast"
)

func TestBroadcast1(t *testing.T) {
	bc := broadcast.NewBroadcast[int]()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *broadcast.Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				for _, v := range s.Dequeue() {
					fmt.Println("recv", v)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(bc.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bc.Enqueue(data)

	time.Sleep(time.Second)
	bc.Close()
	wg.Wait()
}

func TestBroadcast2(t *testing.T) {
	br := broadcast.NewBroadcast[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(s *broadcast.Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				m, ok := s.Pop()
				for ok {
					fmt.Println("recv1", m)
					time.Sleep(500 * time.Millisecond)
					m, ok = s.Pop()
				}
			}
		}
	}(br.Subscribe())

	go func(s *broadcast.Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				m, ok := s.Pop()
				for ok {
					fmt.Println("recv2", m)
					time.Sleep(500 * time.Millisecond)
					m, ok = s.Pop()
				}
			}
		}
	}(br.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	br.Enqueue(data)

	time.Sleep(5 * time.Second)
	br.Close()
	wg.Wait()
}
