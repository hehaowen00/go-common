package common

import (
	"testing"
)

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
