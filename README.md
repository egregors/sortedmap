# 📚 sortedmap

`sortedmap` provides an effective sorted map implementation for Go.
It uses a heap to maintain order and iterators under the hood.

---

[![Build Status](https://github.com/egregors/sortedmap/workflows/build/badge.svg)](https://github.com/egregors/sortedmap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/egregors/sortedmap)](https://goreportcard.com/report/github.com/egregors/sortedmap)
[![Coverage Status](https://coveralls.io/repos/github/egregors/sortedmap/badge.svg)](https://coveralls.io/github/egregors/sortedmap)
[![godoc](https://godoc.org/github.com/egregors/sortedmap?status.svg)](https://godoc.org/github.com/egregors/sortedmap)

## Features

* 🚀 Efficient sorted map implementation
* 🔧 Customizable sorting by key or value
* 🐈 Zero dependencies
* 📦 Easy to use API (inspired by the stdlib `maps` and `slices` packages)

## Installation

To install the package, run:

```sh
go get github.com/egregors/sortedmap
```

## Usage

Here's a quick example of how to use the `sortedmap` package:

```go
package main

import (
	"fmt"

	sm "github.com/egregors/sortedmap"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// Create a new map sorted by keys
	m := sm.NewFromMap(map[string]int{
		"Bob":   31,
		"Alice": 26,
		"Eve":   84,
	}, func(i, j sm.KV[string, int]) bool {
		return i.Key < j.Key
	})

	fmt.Println(m.Collect())
	// Output: map[Alice:26 Bob:31 Eve:84]

	m.Insert("Charlie", 34)
	fmt.Println(m.Collect())
	// Output: map[Alice:26 Bob:31 Charlie:34 Eve:84]

	m.Delete("Bob")
	fmt.Println(m.Collect())
	// Output: map[Alice:26 Charlie:34 Eve:84]

	// Create a new map sorted by values
	m2 := sm.NewFromMap(map[string]Person{
		"Bob":   {"Bob", 31},
		"Alice": {"Alice", 26},
		"Eve":   {"Eve", 84},
	}, func(i, j sm.KV[string, Person]) bool {
		return i.Val.Age < j.Val.Age
	})

	fmt.Println(m2.Collect())
	// Output: map[Alice:{Alice 26} Bob:{Bob 31} Eve:{Eve 84}]

	// Create a new map sorted by values but if the values are equal, sort by keys
	m3 := sm.NewFromMap(map[string]Person{
		"Bob":   {"Bob", 26},
		"Alice": {"Alice", 26},
		"Eve":   {"Eve", 84},
	}, func(i, j sm.KV[string, Person]) bool {
		if i.Val.Age == j.Val.Age {
			return i.Key < j.Key
		}

		return i.Val.Age < j.Val.Age
	})

	fmt.Println(m3.Collect())
	// Output: map[Alice:{Alice 26} Bob:{Bob 26} Eve:{Eve 84}]
}

```

## API and Complexity

| Method          | Description                                                          | Complexity |
|-----------------|----------------------------------------------------------------------|------------|
| `New`           | Creates a new `SortedMap` with a comparison function                 | O(1)       |
| `NewFromMap`    | Creates a new `SortedMap` from an existing map with a comparison     | O(n log n) |
| `Get`           | Retrieves the value associated with a key                            | O(1)       |
| `Delete`        | Removes a key-value pair from the map                                | O(n)       |
| `All`           | Returns a sequence of all key-value pairs in the map                 | O(n log n) |
| `Keys`          | Returns a sequence of all keys in the map                            | O(n log n) |
| `Values`        | Returns a sequence of all values in the map                          | O(n log n) |
| `Insert`        | Adds or updates a key-value pair in the map                          | O(log n)   |
| `Collect`       | Returns  a regular map with an *unordered* content off the SortedMap | O(n log n) |
| `CollectAll`    | Returns a slice of key-value pairs                                   | O(n log n) |
| `CollectKeys`   | Returns a slice of the map’s keys                                    | O(n log n) |
| `CollectValues` | Returns a slice of the map's values                                  | O(n log n) |
| `Len`           | Returns length of underlying map                                     | O(1)       |

## Benchmarks

```shell
BenchmarkNew-10                         165887913           7.037 ns/op
BenchmarkNewFromMap-10                  419106              2716 ns/op
BenchmarkSortedMap_Get-10               191580795           5.327 ns/op
BenchmarkSortedMap_Delete-10            3328420             365.0 ns/op
BenchmarkSortedMap_All-10               1000000000          0.3116 ns/op
BenchmarkSortedMap_Keys-10              1000000000          0.3118 ns/op
BenchmarkSortedMap_Values-10            1000000000          0.3117 ns/op
BenchmarkSortedMap_Insert-10            6665839             182.5 ns/op
BenchmarkSortedMap_Collect-10           649450              1835 ns/op
BenchmarkSortedMap_CollectAll-10        1237276             972.4 ns/op
BenchmarkSortedMap_CollectKeys-10       1250041             964.9 ns/op
BenchmarkSortedMap_CollectValues-10     1294760             927.7 ns/op
BenchmarkSortedMap_Len-10               1000000000          0.3176 ns/op
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
