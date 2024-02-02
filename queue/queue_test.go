package queue_test

import (
	"testing"

	"github.com/hehaowen00/go-common/queue"
)

func TestQueue1(t *testing.T) {
	q := queue.NewQueue[int]()

	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}

	if q.Len() != 10 {
		t.FailNow()
	}

	for i := 0; i < 10; i++ {
		v, ok := q.Dequeue()
		if !ok {
			t.FailNow()
		} else if v != i {
			t.FailNow()
		}
	}

	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}

	if q.Len() != 10 {
		t.FailNow()
	}
}

func BenchmarkQueue1(b *testing.B) {
	b.StopTimer()

	q := queue.NewQueue[int]()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkQueue2(b *testing.B) {
	b.StopTimer()

	q := queue.NewQueue[int]()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}
}

func BenchmarkQueue3(b *testing.B) {
	b.StopTimer()

	q := queue.NewQueue[int]()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < b.N; i++ {
		q.Dequeue()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}
