package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannel1(t *testing.T) {
	channel := NewChannel[int]()

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
	}(channel.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	channel.Enqueue(data)

	time.Sleep(time.Second)
	channel.Close()
	wg.Wait()
}

func TestChannel2(t *testing.T) {
	channel := NewChannel[int]()

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
	}(channel.Subscribe())

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
	}(channel.Subscribe())

	for i := 0; i < 10; i++ {
		channel.Push(i)
	}

	time.Sleep(3 * time.Second)
	channel.Close()
	wg.Wait()
}

func TestChannel3(t *testing.T) {
	channel := NewChannel[int]()

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
	}(channel.Subscribe())

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
	}(channel.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	channel.Enqueue(data)

	time.Sleep(3 * time.Second)
	channel.Close()
	wg.Wait()
}

func TestChannel4(t *testing.T) {
	channel := NewChannel[int]()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	channel.Enqueue(data)

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
				time.Sleep(1 * time.Second)
			}
		}
	}(channel.Subscribe())

	time.Sleep(time.Second)
	channel.Close()
	wg.Wait()
}

func TestChannel5(t *testing.T) {
	channel := NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for i := 0; i < 5; i++ {
			channel.Push(i)
		}
		wg.Done()
	}()

	go func() {
		for i := 5; i < 10; i++ {
			channel.Push(i)
		}
		wg.Done()
	}()

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
			}
		}
	}(channel.Subscribe())

	time.Sleep(2 * time.Second)
	channel.Close()
	wg.Wait()
}
