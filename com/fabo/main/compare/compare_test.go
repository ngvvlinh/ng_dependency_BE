package compare

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Test int

type A struct {
	A int
	B string
	C *B `compare:"ignore"`
	D time.Time
	E []int
	F Test
	G map[int]string
	H interface{}
}

type B struct {
	BA int
	BB *C
}

type C struct {
	CA int
}

func TestCompare(t *testing.T) {
	now := time.Now()
	t.Run("test compare", func(t *testing.T) {
		tests := []struct {
			name      string
			item      interface{}
			otherItem interface{}
			result    bool
		}{
			{
				name:      "compare string with string",
				item:      "1000",
				otherItem: "100",
				result:    false,
			},
			{
				name:      "compare string with int",
				item:      "1000",
				otherItem: 1000,
				result:    false,
			},
			{
				name: "compare struct with int",
				item: &A{
					A: 10,
					B: "10",
					C: &B{BA: 10, BB: &C{CA: 10}},
					D: time.Now(),
					E: []int{10, 20},
					F: Test(5),
					G: map[int]string{5: "10"},
				},
				otherItem: 10000,
				result:    false,
			},
			{
				name: "compare struct with struct",
				item: &A{
					A: 10,
					B: "10",
					C: &B{BA: 10, BB: &C{CA: 10}},
					D: now,
					E: []int{10, 20},
					F: Test(5),
					G: map[int]string{5: "10"},
				},
				otherItem: &A{
					A: 10,
					B: "10",
					C: &B{BA: 10, BB: &C{CA: 10}},
					D: now,
					E: []int{10, 20},
					F: Test(5),
					G: map[int]string{5: "10"},
				},
				result: true,
			},
			{
				name: "compare struct with struct user ignore",
				item: &A{
					A: 10,
					B: "10",
					C: &B{BA: 10, BB: &C{CA: 11}},
					D: now,
					E: []int{10, 20},
					F: Test(5),
					G: map[int]string{5: "11"},
				},
				otherItem: &A{
					A: 10,
					B: "10",
					C: &B{BA: 10, BB: &C{CA: 11}},
					D: now,
					E: []int{10, 20},
					F: Test(5),
					G: map[int]string{5: "11"},
				},
				result: true,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Compare(test.item, test.otherItem)
				require.Equal(t, test.result, result)
			})
		}
	})
}
