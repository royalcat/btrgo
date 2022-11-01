package btrgo_test

import (
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

		btrgo.Sort(adds)
		btrgo.Sort(expectedAdds)
		if !btrgo.CompareSlices(adds, expectedAdds) {
			t.Errorf("Expected adds: %v, got %v", expectedAdds, adds)
		}

		btrgo.Sort(dels)
		btrgo.Sort(expectedDels)
		if !btrgo.CompareSlices(dels, expectedDels) {
			t.Errorf("Expected empty diff, got adds %v and dels %v", adds, dels)
		}

	}
}
