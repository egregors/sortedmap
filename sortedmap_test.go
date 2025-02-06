package sortedmap

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func ptrTo[T any](v T) *T {
	return &v
}

func ptrVal[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

func panicsWithValue(t *testing.T, expected interface{}, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic with value %v, but no panic occurred", expected)
		} else if r != expected {
			t.Errorf("Expected panic with value %v, but got %v", expected, r)
		}
	}()
	fn()
}

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
				return i.Key < j.Key
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
				return i.Key < j.Key
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
					return i.Key > j.Key
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
					return i.Val < j.Val
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
					return i.Val > j.Val
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

func ExampleNewFromMap() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for k, v := range sm.All() {
		fmt.Println(k, v)
	}
	// Output:
	// Alice 30
	// Bob 42
	// Charlie 25
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
				return i.Key < j.Key
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
				panicsWithValue(t, "less function is required", func() {
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

func ExampleNew() {
	sm := New[map[string]int, string, int](func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	sm.Insert("Alice", 30)
	sm.Insert("Bob", 42)
	for k, v := range sm.All() {
		fmt.Println(k, v)
	}
	// Output:
	// Alice 30
	// Bob 42
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
				return i.Key < j.Key
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
				return i.Key < j.Key
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
				return i.Key < j.Key
			}),
			key:   "I've been waiting for you all this time",
			want:  0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.sm.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get(%v) = %v, want %v", tt.key, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get(%v) = %v, want %v", tt.key, got1, tt.want1)
			}
		})
	}
}

func ExampleSortedMap_Get() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	val, ok := sm.Get("Alice")
	fmt.Println(val, ok)
	// Output:
	// 30 true
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
				return i.Key < j.Key
			}),
			key: "Bob",
			want: map[string]int{
				"Alice":   30,
				"Charlie": 25,
			},
			want1: ptrTo(42),
			want2: true,
		},
		{
			name: "unsuccessful deletion",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Key < j.Key
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
			if ptrVal(delVal) != ptrVal(tt.want1) {
				t.Errorf("Delete(%v) = %v, want %v", tt.key, delVal, tt.want1)
			}
			if ok != tt.want2 {
				t.Errorf("Delete(%v) = %v, want %v", tt.key, ok, tt.want2)
			}
			if got := tt.sm.Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete(%v) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_Delete() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	val, ok := sm.Delete("Alice")
	fmt.Println(ptrVal(val), ok)
	val, ok = sm.Delete("Alice")
	fmt.Println(ptrVal(val), ok)
	// Output:
	// 30 true
	// 0 false
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
				return i.Key < j.Key
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
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("All() = %v, want %v", m, tt.want)
			}
		})
	}
}

func ExampleSortedMap_All() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for k, v := range sm.All() {
		fmt.Println(k, v)
	}
	// Output:
	// Alice 30
	// Bob 42
	// Charlie 25
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
				return i.Key < j.Key
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
				return i.Key > j.Key
			}),
			want: []string{"Charlie", "Bob", "Alice"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := slices.Collect(tt.sm.Keys())
			if !reflect.DeepEqual(keys, tt.want) {
				t.Errorf("Keys() = %v, want %v", keys, tt.want)
			}
		})
	}
}

func ExampleSortedMap_Keys() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for k := range sm.Keys() {
		fmt.Println(k)
	}
	// Output:
	// Alice
	// Bob
	// Charlie
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
				return i.Key < j.Key
			}),
			want: []int{30, 42, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Collect(tt.sm.Values()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_Values() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for v := range sm.Values() {
		fmt.Println(v)
	}
	// Output:
	// 30
	// 42
	// 25
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
				return i.Key < j.Key
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
				return i.Key < j.Key
			}),
			args: args[string, int]{key: "Charlie", val: 25},
			want: map[string]int{
				"Alice":   30,
				"Bob":     42,
				"Charlie": 25,
				"David":   35,
			},
		},
		{
			name: "simple map, replace value",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Key < j.Key
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
			if got := tt.sm.Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert(%v, %v) = %v, want %v", tt.args.key, tt.args.val, got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_Insert() {
	sm := New[map[string]int, string, int](func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	sm.Insert("Alice", 30)
	sm.Insert("Bob", 42)
	for k, v := range sm.All() {
		fmt.Println(k, v)
	}
	// Output:
	// Alice 30
	// Bob 42
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
				return i.Key < j.Key
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
			if got := tt.sm.Collect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_Collect() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	fmt.Println(sm.Collect())
	// Unordered output:
	// map[Alice:30 Bob:42 Charlie:25]
}

func TestSortedMap_CollectAll(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want []KV[K, V]
	}
	tests := []testCase[map[int]string, int, string]{
		{
			name: "empty map",
			sm: NewFromMap(map[int]string{}, func(i, j KV[int, string]) bool {
				return i.Val < j.Val
			}),
			want: []KV[int, string]{},
		},
		{
			name: "map with 5 elements",
			sm: NewFromMap(map[int]string{
				1: "one",
				3: "three",
				2: "two",
				5: "five",
				4: "four",
			}, func(i, j KV[int, string]) bool {
				return i.Key < j.Key
			}),
			want: []KV[int, string]{
				{1, "one"},
				{2, "two"},
				{3, "three"},
				{4, "four"},
				{5, "five"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.CollectAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_CollectAll() {
	sm := NewFromMap(map[int]string{
		1: "one",
		3: "three",
		2: "two",
		5: "five",
		4: "four",
	}, func(i, j KV[int, string]) bool {
		return i.Key < j.Key
	})
	fmt.Println(sm.CollectAll())
	// Output:
	// [{1 one} {2 two} {3 three} {4 four} {5 five}]
}

func TestSortedMap_CollectKeys(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want []K
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm: NewFromMap(map[string]int{}, func(i, j KV[string, int]) bool {
				return i.Val < j.Val
			}),
			want: []string{},
		},
		{
			name: "map with 5 elements sorted by value",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Val < j.Val
			}),
			want: []string{"Charlie", "Alice", "Bob"},
		},
		{
			name: "map with 5 elements sorted by key",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Key < j.Key
			}),
			want: []string{"Alice", "Bob", "Charlie"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.CollectKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_CollectKeys() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	fmt.Println(sm.CollectKeys())
	// Output:
	// [Alice Bob Charlie]
}

func TestSortedMap_CollectValues(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want []V
	}
	tests := []testCase[map[string]int, string, int]{
		{
			name: "empty map",
			sm: NewFromMap(map[string]int{}, func(i, j KV[string, int]) bool {
				return i.Val < j.Val
			}),
			want: []int{},
		},
		{
			name: "map with 5 elements sorted by value",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Val < j.Val
			}),
			want: []int{25, 30, 42},
		},
		{
			name: "map with 5 elements sorted by key",
			sm: NewFromMap(map[string]int{
				"Bob":     42,
				"Alice":   30,
				"Charlie": 25,
			}, func(i, j KV[string, int]) bool {
				return i.Key < j.Key
			}),
			want: []int{30, 42, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.CollectValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSortedMap_CollectValues() {
	sm := NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	fmt.Println(sm.CollectValues())
	// Output:
	// [30 42 25]
}

func TestSortedMap_Len(t *testing.T) {
	type testCase[Map interface{ ~map[K]V }, K comparable, V any] struct {
		name string
		sm   *SortedMap[Map, K, V]
		want int
	}
	tests := []testCase[map[int]int, int, int]{
		{
			name: "empty map",
			sm: NewFromMap(map[int]int{}, func(i, j KV[int, int]) bool {
				return i.Val < j.Val
			}),
			want: 0,
		},
		{
			name: "map with 5 elements",
			sm: NewFromMap(map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}, func(i, j KV[int, int]) bool {
				return i.Val < j.Val
			}),
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

var benchMap = map[string]int{
	"Alice":   30,
	"Bob":     42,
	"Charlie": 25,
	"David":   35,
	"Eve":     20,
	"Frank":   40,
	"Grace":   45,
	"Heidi":   50,
	"Ivan":    55,
	"Judy":    60,
	"Kevin":   65,
	"Lucy":    70,
	"Mary":    75,
	"Nancy":   80,
	"Oliver":  85,
	"Peter":   90,
	"Quincy":  95,
	"Roger":   100,
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New[map[int]int, int, int](func(i, j KV[int, int]) bool {
			return i.Key < j.Key
		})
	}
}

func BenchmarkNewFromMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFromMap(benchMap, func(i, j KV[string, int]) bool {
			return i.Key < j.Key
		})
	}
}

func BenchmarkSortedMap_Get(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Get("Roger")
	}
}

func BenchmarkSortedMap_Delete(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Delete("Bob")
	}
}

func BenchmarkSortedMap_All(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.All()
	}
}

func BenchmarkSortedMap_Keys(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Keys()
	}
}

func BenchmarkSortedMap_Values(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Values()
	}
}

func BenchmarkSortedMap_Insert(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Insert("Berik", 42)
	}
}

func BenchmarkSortedMap_Collect(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Collect()
	}
}

func BenchmarkSortedMap_CollectAll(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.CollectAll()
	}
}

func BenchmarkSortedMap_CollectKeys(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.CollectKeys()
	}
}

func BenchmarkSortedMap_CollectValues(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.CollectValues()
	}
}

func BenchmarkSortedMap_Len(b *testing.B) {
	sm := NewFromMap(benchMap, func(i, j KV[string, int]) bool {
		return i.Key < j.Key
	})
	for i := 0; i < b.N; i++ {
		sm.Len()
	}
}
