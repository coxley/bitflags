package bitflags

import (
	"fmt"
	"iter"

	"golang.org/x/exp/constraints"
)

// Set is a helper for operating over a custom bit flag
//
// The zero value is ready to use
type Set[T constraints.Unsigned] struct {
	current T
}

// NewSet returns a set with the provided flags set
func NewSet[T constraints.Unsigned](flag T, extra ...T) Set[T] {
	var bf Set[T]
	bf.Add(flag, extra...)
	return bf
}

func (f Set[T]) Empty() bool {
	return f.current == 0
}

// Merge flags together and return as a new value
func (f Set[T]) Merge(flags Set[T]) Set[T] {
	return Set[T]{f.current | flags.current}
}

// Add one or more flags to the set
//
// Panics if any value is not on a bit-boundary to protect against misuse.
func (f *Set[T]) Add(flag T, extra ...T) {
	f.checkPower(flag)
	f.current |= flag
	for _, flag := range extra {
		f.checkPower(flag)
		f.current |= flag
	}
}

// Clear one or more flags from the set
//
// Panics if any value is not on a bit-boundary to protect against misuse.
func (f *Set[T]) Clear(flag T, extra ...T) {
	f.checkPower(flag)
	f.current &^= flag

	for _, flag := range extra {
		f.checkPower(flag)
		f.current &^= flag
	}
}

// Has the flag been set?
//
// Panics if any value is not on a bit-boundary to protect against misuse.
func (f Set[T]) Has(flag T) bool {
	f.checkPower(flag)
	return f.current&flag == flag
}

// HasAll returns true if all provided flags are set
//
// Panics if any value is not on a bit-boundary to protect against misuse.
func (f Set[T]) HasAll(flag T, extra ...T) bool {
	if !f.Has(flag) {
		return false
	}

	for _, flag := range extra {
		if !f.Has(flag) {
			return false
		}
	}
	return true
}

// All produces each flag that is currently set
func (f Set[T]) All() iter.Seq[T] {
	if f.current == 0 {
		return func(yield func(T) bool) {}
	}

	return func(yield func(T) bool) {
		set := f.current

		// Loop until all bits have been cleared
		for set != 0 {
			// Produce the least-significant bit
			if !yield(set & -set) {
				return
			}
			// Clear the least-significant bit
			set &= set - 1
		}
	}
}

// checkPower panics on unexpected values for safety (if enabled)
func (f Set[T]) checkPower(flag T) {
	if checkPower != nil {
		checkPower(uint(flag))
	}
}

var checkPower func(flag uint)

// ValidateFlags will panic when flag values are not on bit boundaries when enabled.
func ValidateFlags(enabled bool) {
	if enabled {
		checkPower = func(flag uint) {
			if flag == 0 || flag&(flag-1) != 0 {
				panic(fmt.Sprintf("'%d' doesn't fall on a bit-boundary", flag))
			}
		}
	} else {
		checkPower = nil
	}
}
