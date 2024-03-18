package btrgo_test

import (
	"slices"
	"testing"

	"github.com/royalcat/btrgo"
)

func TestSliceDiffSplited(t *testing.T) {
	{
		a := []int{1, 2, 3}
		b := []int{2, 3, 1}
		adds, dels := btrgo.SliceDiffSplited(a, b)
		if len(adds) != 0 || len(dels) != 0 {
			t.Errorf("Expected empty diff, got adds %v and dels %v", adds, dels)
		}
	}
	{
		expectedAdds := []int{4, 8}
		expectedDels := []int{7, 1}

		a := []int{7, 3, 2, 1}
		b := []int{2, 8, 2, 3, 4, 3}
		dels, adds := btrgo.SliceDiffSplited(a, b)

		slices.Sort(adds)
		slices.Sort(expectedAdds)
		if !slices.Equal(adds, expectedAdds) {
			t.Errorf("Expected adds: %v, got %v", expectedAdds, adds)
		}

		slices.Sort(dels)
		slices.Sort(expectedDels)
		if !slices.Equal(dels, expectedDels) {
			t.Errorf("Expected empty diff, got adds %v and dels %v", adds, dels)
		}

	}
}

func TestSliceUnique(t *testing.T) {
	values := []pair[[]int]{
		{[]int{1, 2, 1}, []int{1, 2}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{}, []int{}},
		{[]int{3, 2}, []int{2, 3}},
		{[]int{-1, 1, -1}, []int{-1, 1}},
	}
	for _, v := range values {
		res := btrgo.SliceUnique(v.I)
		slices.Sort(res)
		if !slices.Equal(v.J, res) {
			t.Errorf("For slice %v expected unique slice %v got %v", v.I, v.J, res)
		}
	}
}

func testSliceUnic(t *testing.T, i, r []int) {
	res := btrgo.SliceUnique(i)
	if slices.Equal(r, res) {
		t.Errorf("Expected slice %v got %v", r, res)
	}
}

type pair[V any] struct {
	I, J V
}
