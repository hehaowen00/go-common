package broadcast

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroadcast1(t *testing.T) {
	broadcast := NewBroadcast[int]()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *Subscriber[int]) {
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
	}(broadcast.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	broadcast.Enqueue(data)

	time.Sleep(time.Second)
	broadcast.Close()
	wg.Wait()
}

func TestBroadcast2(t *testing.T) {
	broadcast := NewBroadcast[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(s *Subscriber[int]) {
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
	}(broadcast.Subscribe())

	go func(s *Subscriber[int]) {
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
	}(broadcast.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	broadcast.Enqueue(data)

	time.Sleep(5 * time.Second)
	broadcast.Close()
	wg.Wait()
}
