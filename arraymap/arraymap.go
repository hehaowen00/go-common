package arraymap

import (
	"cmp"
	"slices"
)

type ArrayMap[K cmp.Ordered, V any] struct {
	keys []K
	vals []V
}

type Entry[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

func NewArrayMap[K cmp.Ordered, V any]() *ArrayMap[K, V] {
	return &ArrayMap[K, V]{
		keys: nil,
		vals: nil,
	}
}

func (m *ArrayMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *ArrayMap[K, V]) Keys() []K {
	return m.keys
}

func (m *ArrayMap[K, V]) Values() []V {
	return m.vals
}

func (m *ArrayMap[K, V]) Clear() {
	m.keys = nil
	m.vals = nil
}

func (m *ArrayMap[K, V]) Get(key K) (V, bool) {
	i, ok := slices.BinarySearch(m.keys, key)
	if ok {
		return m.vals[i], true
	}

	return *new(V), false
}

func (m *ArrayMap[K, V]) Set(key K, val V) {
	i, ok := slices.BinarySearch(m.keys, key)
	if !ok {
		m.keys = slices.Insert(m.keys, i, key)
		m.vals = slices.Insert(m.vals, i, val)
	}
}

func (m *ArrayMap[K, V]) Delete(key K) {
	i, ok := slices.BinarySearch(m.keys, key)
	if ok {
		m.keys = slices.Delete(m.keys, i, i+1)
		m.vals = slices.Delete(m.vals, i, i+1)
	}
}

func (m *ArrayMap[K, V]) Has(key K) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *ArrayMap[K, V]) ToMap() map[K]V {
	data := make(map[K]V)

	for i := range m.keys {
		data[m.keys[i]] = m.vals[i]
	}

	return data
}

func (m *ArrayMap[K, V]) FromMap(data map[K]V) {
	for k, v := range data {
		m.Set(k, v)
	}
}

func (m *ArrayMap[K, V]) ToSlice() []Entry[K, V] {
	var entries []Entry[K, V]

	for i := range m.keys {
		entries = append(entries, Entry[K, V]{
			Key:   m.keys[i],
			Value: m.vals[i],
		})
	}

	return entries
}

func (m *ArrayMap[K, V]) FromSlice(entries []Entry[K, V]) {
	for i := range entries {
		m.Set(entries[i].Key, entries[i].Value)
	}
}
