package common

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBus1(t *testing.T) {
	bus := NewBus[int]()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				m, ok := s.Dequeue()
				if ok {
					for _, v := range m {
						fmt.Println("recv", v)
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}(bus.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bus.Enqueue(data)

	time.Sleep(time.Second)
	bus.Close()
	wg.Wait()
}

func TestBus2(t *testing.T) {
	bus := NewBus[int]()

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
					time.Sleep(1 * time.Second)
					m, ok = s.Pop()
				}
			}
		}
	}(bus.Subscribe())

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
					time.Sleep(1 * time.Second)
					m, ok = s.Pop()
				}
			}
		}
	}(bus.Subscribe())

	for i := 0; i < 10; i++ {
		bus.Push(i)
	}

	time.Sleep(6 * time.Second)
	bus.Close()
	wg.Wait()
}

func TestBus3(t *testing.T) {
	bus := NewBus[int]()

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
					time.Sleep(1 * time.Second)
					m, ok = s.Pop()
				}
			}
		}
	}(bus.Subscribe())

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
					time.Sleep(1 * time.Second)
					m, ok = s.Pop()
				}
			}
		}
	}(bus.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bus.Enqueue(data)

	time.Sleep(2 * time.Second)
	bus.Close()
	wg.Wait()
}

func TestBus4(t *testing.T) {
	bus := NewBus[int]()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bus.Enqueue(data)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				m, ok := s.Dequeue()
				if ok {
					for _, v := range m {
						fmt.Println("recv", v)
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}(bus.Subscribe())

	time.Sleep(time.Second)
	bus.Close()
	wg.Wait()
}

func TestBus5(t *testing.T) {
	bus := NewBus[int]()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for i := 0; i < 5; i++ {
			bus.Push(i)
		}
		wg.Done()
	}()

	go func() {
		for i := 5; i < 10; i++ {
			bus.Push(i)
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
				m, ok := s.Dequeue()
				if ok {
					for _, v := range m {
						fmt.Println("recv", v)
					}
				}
			}
		}
	}(bus.Subscribe())

	time.Sleep(2 * time.Second)
	bus.Close()
	wg.Wait()
}
