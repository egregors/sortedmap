package sortedmap

import (
	"container/heap"
	"iter"
)

// SortedMap is a map-like struct that keeps sorted by key or value.
// It uses a heap to maintain the order.
type SortedMap[Map ~map[K]V, K comparable, V any] struct {
	m Map
	h *kvHeap[K, V]
}

// New creates a new SortedMap with `less` as the comparison function
// The complexity is O(1)
func New[Map ~map[K]V, K comparable, V any](less func(i, j KV[K, V]) bool) *SortedMap[Map, K, V] {
	if less == nil {
		panic("less function is required")
	}

	return &SortedMap[Map, K, V]{
		m: make(Map),
		h: newKvHeap(less),
	}
}

// NewFromMap creates a new SortedMap with `less` as the comparison function and populates it with the contents of `m`.
// The complexity is O(n log n) where n = len(m).
func NewFromMap[Map ~map[K]V, K comparable, V any](m Map, less func(i, j KV[K, V]) bool) *SortedMap[Map, K, V] {
	sm := New[Map, K, V](less)
	for k, v := range m {
		sm.Insert(k, v)
	}

	return sm
}

// Get returns the value associated with the key and a boolean indicating if the key exists in the map
// The complexity is O(1)
func (sm *SortedMap[Map, K, V]) Get(key K) (V, bool) {
	val, exists := sm.m[key]

	return val, exists
}

// Delete removes the key from the map and returns the value associated with the key and a boolean indicating
// if the key existed in the map.
// The complexity is O(n) where n = len(sm.h.xs)
func (sm *SortedMap[Map, K, V]) Delete(key K) (val *V, existed bool) {
	delete(sm.m, key)
	// TODO: in order to remove the element from the heap, we need to full scan it with O(n).
	// 	probably we can use a map to store the index of the element in the heap, but the problem
	//  is that element indexes will constantly change as we remove or add elements.
	for i, el := range sm.h.xs {
		if el.Key == key {
			val = &el.Val
			heap.Remove(sm.h, i)

			return val, true
		}
	}

	return (*V)(nil), false
}

// All returns a sequence of key-value pairs in the map
func (sm *SortedMap[Map, K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		tempHeap := *sm.h
		for tempHeap.Len() > 0 {
			el := heap.Pop(&tempHeap).(KV[K, V])
			if !yield(el.Key, el.Val) {
				return
			}
		}
	}
}

// Keys returns a sequence of keys in the map
func (sm *SortedMap[Map, K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		tempHeap := *sm.h
		for tempHeap.Len() > 0 {
			el := heap.Pop(&tempHeap).(KV[K, V])
			if !yield(el.Key) {
				return
			}
		}
	}
}

// Values returns a sequence of values in the map
func (sm *SortedMap[Map, K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		tempHeap := *sm.h
		for tempHeap.Len() > 0 {
			el := heap.Pop(&tempHeap).(KV[K, V])
			if !yield(el.Val) {
				return
			}
		}
	}
}

// Insert adds a key-value pair to the map. If the key already exists, the value is updated.
func (sm *SortedMap[Map, K, V]) Insert(key K, val V) {
	if _, exists := sm.m[key]; exists {
		sm.Delete(key)
	}
	sm.m[key] = val
	heap.Push(sm.h, KV[K, V]{key, val})
}

// Collect returns a map with the same contents as the SortedMap
func (sm *SortedMap[Map, K, V]) Collect() Map {
	m := make(Map)
	for key, val := range sm.All() {
		m[key] = val
	}

	return m
}

// Len returns length of underlying map
func (sm *SortedMap[Map, K, V]) Len() int {
	return len(sm.m)
}
