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
		t.Logf("Adds %v Dels %v", adds, dels)
		if len(adds) != 0 || len(dels) != 0 {
			t.Errorf("Expected empty diff, got adds %v and dels %v", adds, dels)
		}
	}
	{
		a := []int{7, 3, 2, 1}
		b := []int{2, 8, 2, 3, 4, 3}
		adds, dels := btrgo.SliceDiffSplited(a, b)
		t.Logf("Adds %v Dels %v", adds, dels)
		// TOOD t.Fail()
	}
}
