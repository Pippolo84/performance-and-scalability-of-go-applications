package insertremove

import (
	"fmt"
	"sort"
	"testing"
)

var implementations = []struct {
	descr string
	obj   InserterRemover
}{
	{
		descr: "slice",
		obj:   NewSliceImpl(),
	},
	{
		descr: "list",
		obj:   NewListImpl(),
	},
}

var testCases = []struct {
	descr       string
	input       int
	expectError bool
}{
	{
		descr:       "negative value",
		input:       -5,
		expectError: true,
	},
	{
		descr:       "single element",
		input:       1,
		expectError: false,
	},
	{
		descr:       "small number of elements",
		input:       10,
		expectError: false,
	},
	{
		descr:       "large number of elements",
		input:       25000,
		expectError: false,
	},
}

func TestInsertRemove(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.descr, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.descr, func(t *testing.T) {
					err := impl.obj.Insert(tc.input)

					if tc.expectError {
						if err == nil {
							t.Fatalf("Insert(%d) expected to fail, but returned without error", tc.input)
						}

						values := impl.obj.Values()
						if len(values) != 0 {
							t.Fatalf("after failed Insert, len(container) = %d\texpected %d", len(values), 0)
						}
					} else {
						if err != nil {
							t.Fatalf("Insert(%d) returned unexpected error", tc.input)
						}

						values := impl.obj.Values()
						if len(values) != tc.input {
							t.Fatalf("after Insert, len(container) = %d\texpected %d", len(values), tc.input)
						}

						sorted := sort.SliceIsSorted(values, func(i, j int) bool {
							return values[i] < values[j]
						})
						if !sorted {
							t.Fatalf("container is not sorted after Insert(%d)", tc.input)
						}

						impl.obj.Remove()

						values = impl.obj.Values()
						if len(values) != 0 {
							t.Fatalf("after Remove, len(container) = %d\texpected %d", len(values), 0)
						}
					}
				})
			}
		})
	}
}

func Benchmark(b *testing.B) {
	for _, impl := range implementations {
		b.Run(impl.descr, func(b *testing.B) {
			for _, n := range []int{0, 8, 16, 32, 64, 256, 1024, 4096, 16384, 32768} {
				b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						impl.obj.Insert(n)
						impl.obj.Remove()
					}
				})
			}
		})
	}
}
