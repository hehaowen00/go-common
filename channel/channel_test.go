package channel_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hehaowen00/go-common/channel"
)

func TestChannel1(t *testing.T) {
	ch := channel.NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch.Enqueue(data)

	time.Sleep(time.Second)
	ch.Close()
	wg.Wait()
}

func TestChannel2(t *testing.T) {
	ch := channel.NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	for i := 0; i < 10; i++ {
		ch.Push(i)
	}

	time.Sleep(3 * time.Second)
	ch.Close()
	wg.Wait()
}

func TestChannel3(t *testing.T) {
	ch := channel.NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch.Enqueue(data)

	time.Sleep(3 * time.Second)
	ch.Close()
	wg.Wait()
}

func TestChannel4(t *testing.T) {
	ch := channel.NewChannel[int]()

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch.Enqueue(data)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	time.Sleep(time.Second)
	ch.Close()
	wg.Wait()
}

func TestChannel5(t *testing.T) {
	ch := channel.NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for i := 0; i < 5; i++ {
			ch.Push(i)
		}
		wg.Done()
	}()

	go func() {
		for i := 5; i < 10; i++ {
			ch.Push(i)
		}
		wg.Done()
	}()

	go func(s *channel.Subscriber[int]) {
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
	}(ch.Subscribe())

	time.Sleep(2 * time.Second)
	ch.Close()
	wg.Wait()
}

func BenchmarkChannel1(b *testing.B) {
	ch := channel.NewChannel[int]()

	for i := 0; i < b.N; i++ {
		ch.Push(i)
	}
}

func BenchmarkChannel2(b *testing.B) {
	b.StopTimer()

	ch := channel.NewChannel[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ch.Push(i)
	}

	go func(s *channel.Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				v, ok := s.Pop()
				for ok {
					_ = v
					if !ok {
						break
					}
					v, ok = s.Pop()
				}
			}
		}
	}(ch.Subscribe())

	go func(s *channel.Subscriber[int]) {
		for {
			select {
			case <-s.Closed():
				wg.Done()
				return
			case <-s.Notify():
				v, ok := s.Pop()
				for ok {
					_ = v
					if !ok {
						break
					}
					v, ok = s.Pop()
				}
			}
		}
	}(ch.Subscribe())

	ch.Close()

	wg.Wait()
}

func BenchmarkChan1(b *testing.B) {
	ch := make(chan int, b.N)

	for i := 0; i < b.N; i++ {
		ch <- i
	}
}

func BenchmarkChan2(b *testing.B) {
	b.StopTimer()

	ch := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(2)

	b.StartTimer()

	go func() {
		for v := range ch {
			_ = v
		}
		wg.Done()
	}()

	go func() {
		for v := range ch {
			_ = v
		}
		wg.Done()
	}()

	for i := 0; i < b.N; i++ {
		ch <- i
	}

	close(ch)

	wg.Wait()
}
