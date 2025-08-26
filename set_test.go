package bitflags

import (
	"slices"
	"testing"
)

func TestFlags(t *testing.T) {
	type perms uint8
	const (
		read perms = 1 << iota
		write
		exec
	)

	flags := NewSet(read)
	True(t, flags.Has(read))
	False(t, flags.Has(write))
	False(t, flags.Has(exec))

	flags.Clear(read)
	False(t, flags.Has(read))

	flags.Add(read, exec)
	True(t, flags.HasAll(read, exec))
	False(t, flags.Has(write))

	all := slices.Collect(flags.All())
	Len(t, all, 2)
	Equal(t, []perms{read, exec}, all)

	flags = flags.Merge(NewSet(write))
	all = slices.Collect(flags.All())
	Len(t, all, 3)
	Equal(t, []perms{read, write, exec}, all)
}

func BenchmarkFlags(b *testing.B) {
	type perms uint8
	const (
		read perms = 1 << iota
		write
		exec
	)

	b.Run("iterator", func(b *testing.B) {
		flags := NewSet(read, write, exec)
		for b.Loop() {
			for range flags.All() {
			}
		}
	})

	b.Run("no validation", func(b *testing.B) {
		for b.Loop() {
			var flags Set[perms]
			flags.Add(read)
			flags.Has(read)
		}
	})

	b.Run("validation", func(b *testing.B) {
		ValidateFlags(true)
		defer ValidateFlags(false)
		for b.Loop() {
			var flags Set[perms]
			flags.Add(read)
			flags.Has(read)
		}
	})
}

func Equal[T comparable, S []T](t testing.TB, expected, got S) {
	if !slices.Equal(expected, got) {
		t.Fatalf("slices are not equal:\nexpected: %v\ngot: %v", expected, got)
	}
}

func Len[T any, S []T](t testing.TB, got S, size int) {
	if len(got) != size {
		t.Fatalf("len(slice) != %d: got %d", size, len(got))
	}
}

func True(t testing.TB, v bool) {
	if !v {
		t.Fatalf("expected true")
	}
}

func False(t testing.TB, v bool) {
	if v {
		t.Fatalf("expected false")
	}
}
