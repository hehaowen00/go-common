package ring

import (
	"testing"
)

func TestRing1(t *testing.T) {
	r := NewRing[int]()

	for i := 0; i < 10; i++ {
		r.Push(i)
	}

	for i := 0; i < 10; i++ {
		v, ok := r.Pop()

		if !ok {
			t.FailNow()
		}

		if v != i {
			t.FailNow()
		}
	}
}

func TestRing2(t *testing.T) {
	r := NewRing[int]()

	for i := 0; i < 10; i++ {
		r.Push(i)
	}

	data := r.Dequeue()

	for i := 0; i < 10; i++ {
		if data[i] != i {
			t.FailNow()
		}
	}
}

func BenchmarkRing1(b *testing.B) {
	b.StopTimer()

	r := NewRing[int]()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

func BenchmarkRing2(b *testing.B) {
	b.StopTimer()

	r := NewRing[int]()

	for i := 0; i < b.N; i++ {
		r.Push(i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		r.Pop()
	}
}

func BenchmarkRing3(b *testing.B) {
	b.StopTimer()

	r := NewRing[int]()

	for i := 0; i < b.N; i++ {
		r.Push(i)
	}

	for i := 0; i < b.N; i++ {
		r.Pop()
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}
