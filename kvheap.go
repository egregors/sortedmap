package sortedmap

// KV is a key-value pair.
type KV[K comparable, V any] struct {
	Key K
	Val V
}

type kvHeap[K comparable, V any] struct {
	xs     []KV[K, V]
	lessFn func(i, j KV[K, V]) bool
}

func newKvHeap[K comparable, V any](less func(i, j KV[K, V]) bool) *kvHeap[K, V] {
	return &kvHeap[K, V]{
		xs:     []KV[K, V]{},
		lessFn: less,
	}
}

func (k *kvHeap[K, V]) Len() int           { return len(k.xs) }
func (k *kvHeap[K, V]) Swap(i, j int)      { k.xs[i], k.xs[j] = k.xs[j], k.xs[i] }
func (k *kvHeap[K, V]) Less(i, j int) bool { return k.lessFn(k.xs[i], k.xs[j]) }
func (k *kvHeap[K, V]) Push(x any)         { k.xs = append(k.xs, x.(KV[K, V])) }
func (k *kvHeap[K, V]) Pop() any {
	n := len(k.xs)
	if n == 0 {
		return nil
	}
	x := k.xs[n-1]
	k.xs = k.xs[:n-1]

	return x
}
