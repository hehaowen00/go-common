package arraymap

import (
	"fmt"
	"testing"
)

func TestArrayMap1(t *testing.T) {
	table := NewArrayMap[int, string]()

	for i := 0; i < 1000; i++ {
		table.Set(i, fmt.Sprintf("%d", i))
	}

	for i := 0; i < 1000; i++ {
		v, ok := table.Get(i)
		if ok {
			if v != fmt.Sprintf("%d", i) {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	}
}

func BenchmarkArrayMap1(b *testing.B) {
	b.StopTimer()

	table := NewArrayMap[int, string]()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		table.Set(i, fmt.Sprintf("%d", i))
	}

	for i := 0; i < b.N; i++ {
		table.Get(i)
	}
}

func BenchmarkMap1(b *testing.B) {
	b.StopTimer()

	table := map[int]string{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		table[i] = fmt.Sprintf("%d", i)
	}

	for i := 0; i < b.N; i++ {
		_ = table[i]
	}
}
