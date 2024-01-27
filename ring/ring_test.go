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

func TestRing3(t *testing.T) {
	r := NewRing[int]()

	for i := 0; i < 10; i++ {
		r.Push(i)
	}

	for i := 0; i < 7; i++ {
		v, ok := r.Pop()
		if !ok {
			t.FailNow()
		}

		if v != i {
			t.FailNow()
		}
	}

	for i := 0; i < 13; i++ {
		r.Push(i)
	}

	data := r.Dequeue()

	var expected []int

	for i := 7; i < 10; i++ {
		expected = append(expected, i)
	}

	for i := 0; i < 13; i++ {
		expected = append(expected, i)
	}

	if len(data) != len(expected) {
		t.Fail()
	}

	for i := 0; i < len(data); i++ {
		if data[i] != expected[i] {
			t.Fail()
		}
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
