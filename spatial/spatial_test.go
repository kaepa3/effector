package spatial

import (
	"image"
	"image/color"
	"testing"

	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/testutill"
)

func Test_AverageAndMedianFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)

	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, 1})
	})
	doFunc := AverageFunc(1.0)
	col := doFunc(rgba, 1, 1)
	if col.R != 2 || col.G != 2 || col.B != 2 {
		t.Error("val err")
	}
	col = MedianFunc(rgba, 1, 1)
	if col.R != 2 || col.G != 2 || col.B != 2 {
		t.Error("val err")
	}
}
