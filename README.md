# sortedmap

`sortedmap` provides an effective sorted map implementation for Go.
Below you will find information about the repository, its usage, and the API with complexity details.

## Features

* 🚀 Efficient sorted map implementation
* 🔧 Customizable sorting by key or value

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

	"github.com/egregors/sortedmap"
)

func main() {
	sm := sortedmap.NewFromMap(map[string]int{
		"Bob":     42,
		"Alice":   30,
		"Charlie": 25,
	}, func(i, j sortedmap.KV[string, int]) bool {
		return i.key < j.key
	})

	fmt.Println(sm.Collect())
}
```

## API and Complexity

| Method       | Description                                                             | Complexity |
|--------------|-------------------------------------------------------------------------|------------|
| `New`        | Creates a new `SortedMap` with a custom comparison function             | O(1)       |
| `NewFromMap` | Creates a new `SortedMap` from an existing map with a custom comparison | O(n log n) |
| `Get`        | Retrieves the value associated with a key                               | O(1)       |
| `Delete`     | Removes a key-value pair from the map                                   | O(n)       |
| `All`        | Returns a sequence of all key-value pairs in the map                    | O(n log n) |
| `Keys`       | Returns a sequence of all keys in the map                               | O(n log n) |
| `Values`     | Returns a sequence of all values in the map                             | O(n log n) |
| `Insert`     | Adds or updates a key-value pair in the map                             | O(log n)   |
| `Collect`    | Returns a map with the same contents as the `SortedMap`                 | O(n log n) |

## Contributing

We welcome contributions! Please see the `CONTRIBUTING.md` file for guidelines on how to contribute to this project.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
