# bitflags

[![Go Reference](https://pkg.go.dev/badge/github.com/coxley/bitflags.svg)](https://pkg.go.dev/github.com/coxley/bitflags)

This covers a specific use-case for bit flags where:

- You want the flags themselves defined separately.
- You want to iterate over set values, and rely on exhaustive-switch linters to make you account for everything.
- You want generic support for "safety" (convenience), and support any type with
  `constraints.Unsigned` storage.

A code block is worth 1000 words:

```go
package main

import (
	"fmt"

	"github.com/coxley/bitflags"
)

type perms uint8

const (
	read perms = 1 << iota
	write
	exec
)

func main() {
	var flags bitflags.Set[perms]
	flags.Add(read, exec)    // alt: bitflags.NewSet(read, exec)
	flags.HasAll(read, exec) // true
	flags.Has(write)         // false

	// Output:
	//   read: set
	//   exec: set
	for flag := range flags.All() {
		switch flag {
		case read:
			fmt.Println("read: set")
		case write:
			fmt.Println("write: set")
		case exec:
			fmt.Println("exec: set")
		}
	}
}
```
