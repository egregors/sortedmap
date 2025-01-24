package sortedmap

import (
	"reflect"
	"slices"
	"testing"

	"github.com/egregors/sortedmap/ptr"

	"github.com/stretchr/testify/assert"
)

func TestNewFromMap(t *testing.T) {
	type args[Map interface{ ~map[K]V }, K comparable, V any] struct {
		m    Map
		less func(i, j KV[K, V]) bool
	}
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		args args[Map, K, V]
		want map[K]V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			args: args[map[string]int, string, int]{m: map[string]int{}, less: func(i, j KV[string, int]) bool {
				return i.key < j.key
			}},
			want: map[string]int{},
		},
		{
			name: "simple map – ascending keys",
			args: args[map[string]int, string, int]{m: map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, less: func(i, j KV[string, int]) bool {
				return i.key < j.key
			}},
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
			},
		},
		{
			name: "simple map – descending keys",
			args: args[map[string]int, string, int]{m: map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			},
				less: func(i, j KV[string, int]) bool {
					return i.key > j.key
				},
			},
			want: map[string]int{
				"Charlie": 25,
				"Bob":     42,
				"Alice":   30,
			},
		},
		{
			name: "simple map – ascending values",
			args: args[map[string]int, string, int]{m: map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			},
				less: func(i, j KV[string, int]) bool {
					return i.val < j.val
				},
			},
			want: map[string]int{
				"Charlie": 25,
				"Alice":   30,
				"Bob":     42,
			},
		},
		{
			name: "simple map – descending values",
			args: args[map[string]int, string, int]{m: map[string]int{
				"Bob":     42,
				"Charlie": 25,
				"Alice":   30,
			},
				less: func(i, j KV[string, int]) bool {
					return i.val > j.val
				},
			},
			want: map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromMap(tt.args.m, tt.args.less).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args[K comparable, V any] struct {
		less func(i, j KV[K, V]) bool
	}
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name        string
		args        args[K, V]
		shouldPanic bool
		want        map[K]V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			args: args[string, int]{less: func(i, j KV[string, int]) bool {
				return i.key < j.key
			}},
			want: map[string]int{},
		},
		{
			name:        "missing less function",
			args:        args[string, int]{less: nil},
			shouldPanic: true,
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.PanicsWithValue(t, "less function is required", func() {
					New[map[string]int, string, int](tt.args.less)
				})

				return
			}
			if got := New[map[string]int, string, int](tt.args.less).Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortedMap_Get(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name  string
		sm    *SortedMap[Map, K, V]
		key   K
		want  V
		want1 bool
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm: NewFromMap(map[string]int{}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			key:   "Berik the Cat",
			want:  0,
			want1: false,
		},
		{
			name: "simple map, existing key",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			key:   "Alice",
			want:  30,
			want1: true,
		},
		{
			name: "simple map, non-existing key",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			key:   "I've been waiting for you all this time",
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.sm.Get(tt.key)
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.key)
			assert.Equalf(t, tt.want1, got1, "Get(%v)", tt.key)
		})
	}
}

func TestSortedMap_Delete(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name  string
		sm    *SortedMap[Map, K, V]
		key   K
		want  map[K]V
		want1 *V
		want2 bool
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "successful deletion",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			key: "Bob",
			want: map[string]int{
				"Alice":   30,
				"Charlie": 25,
			},
			want1: ptr.To(42),
			want2: true,
		},
		{
			name: "unsuccessful deletion",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
			},
			want1: nil,
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delVal, ok := tt.sm.Delete(tt.key)
			assert.Equalf(t, tt.want1, delVal, "Delete(%v)", tt.key)
			assert.Equalf(t, tt.want2, ok, "Delete(%v)", tt.key)
			assert.Equalf(t, tt.want, tt.sm.Collect(), "Delete(%v)", tt.key)
		})
	}
}

func TestSortedMap_All(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want map[K]V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm:   New[map[string]int, string, int](func(i, j KV[string, int]) bool { return true }),
			want: map[string]int{},
		},
		{
			name: "simple map",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := make(map[string]int)
			for k, v := range tt.sm.All() {
				m[k] = v
			}
			assert.Equalf(t, tt.want, m, "All()")
		})
	}
}

func TestSortedMap_Keys(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want []K
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm:   New[map[string]int, string, int](func(i, j KV[string, int]) bool { return true }),
			want: nil,
		},
		{
			name: "simple map",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			want: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name: "simple map – descending keys",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key > j.key
			}),
			want: []string{"Charlie", "Bob", "Alice"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := slices.Collect(tt.sm.Keys())
			assert.Equalf(t, tt.want, keys, "Keys()")
		})
	}
}

func TestSortedMap_Values(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want []V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm:   New[map[string]int, string, int](func(i, j KV[string, int]) bool { return true }),
			want: nil,
		},
		{
			name: "simple map",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			want: []int{30, 42, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, slices.Collect(tt.sm.Values()), "Values()")
		})
	}
}

func TestSortedMap_Insert(t *testing.T) {
	type args[K comparable, V any] struct {
		key K
		val V
	}
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm: New[map[string]int, string, int](func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			args: args[string, int]{key: "Alice", val: 30},
			want: map[string]int{
				"Alice": 30,
			},
		},
		{
			name: "simple map, new key",
			sm: NewFromMap(map[string]int{
				"Bob":   42,
				"Alice": 30,
				"David": 35,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			args: args[string, int]{key: "Charlie", val: 25},
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
			},
		},
		{
			name: "simple map, replace value",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			args: args[string, int]{key: "Alice", val: 35},
			want: map[string]int{
				"Alice":   35,
				"Bob":     42,
				"Charlie": 25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sm.Insert(tt.args.key, tt.args.val)
		})
	}
}

func TestSortedMap_Collect(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want Map
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm:   New[map[string]int, string, int](func(i, j KV[string, int]) bool { return true }),
			want: map[string]int{},
		},
		{
			name: "simple map",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.key < j.key
			}),
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.sm.Collect(), "Collect()")
		})
	}
}
