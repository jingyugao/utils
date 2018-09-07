package utils_test

import (
	"testing"

	"github.com/jingyugao/utils"
)

func TestSet1(t *testing.T) {

	bitmap := utils.NewBitmap(100)
	for i := uint32(0); i < 100; i++ {
		bingo := bitmap.Test(i)
		if !bingo {
			t.Log("pass")
		} else {
			t.Error("failed")
		}

		bitmap.Set(i)
		bingo = bitmap.Test(i)
		if bingo {
			t.Log("pass")
		} else {
			t.Error("failed")
		}
	}

}
func TestUnSet(t *testing.T) {

	bitmap := utils.NewBitmap(100)
	for i := uint32(0); i < 100; i++ {
		bitmap.Set(i)
		bingo := bitmap.Test(i)
		if bingo {
			t.Log("pass")
		} else {
			t.Error("failed")
		}
		bitmap.UnSet(i)
		bingo = bitmap.Test(i)
		if !bingo {
			t.Log("pass")
		} else {
			t.Error("failed")
		}
	}
}
