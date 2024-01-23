package common

import (
	"testing"
)

func TestQueue1(t *testing.T) {
	queue := NewQueue[int]()

	for i := 0; i < 10; i++ {
		queue.Enqueue(i)
	}

	if queue.Len() != 10 {
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		v := queue.Dequeue()

		if v.Ok() {
			if v.Unwrap() != i {
				t.Fail()
			}
		}
	}

	for i := 0; i < 10; i++ {
		queue.Enqueue(i)
	}

	if queue.Len() != 10 {
		t.Fail()
	}
}

func BenchmarkQueue1(b *testing.B) {
	b.StopTimer()

	queue := NewQueue[int]()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
}

func BenchmarkQueue2(b *testing.B) {
	b.StopTimer()

	queue := NewQueue[int]()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}
}

func BenchmarkQueue3(b *testing.B) {
	b.StopTimer()

	queue := NewQueue[int]()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}

	for i := 0; i < b.N; i++ {
		queue.Dequeue()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
}
