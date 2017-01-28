package twoval

import (
	"image"
	"image/color"
	"testing"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/testutill"
)

func Test_StaticThreshold(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		val := uint16(x)
		rgba.Set(x, y, color.RGBA64{val, val, val, 1})
	})
	doFunc := StaticThreshold(2)
	col := doFunc(rgba, 1, 1)
	if col.R != 0 || col.G != 0 || col.B != 0 {
		t.Error("val err")
	}
	col = doFunc(rgba, 2, 1)
	if col.R != ex.ColorWidth || col.G != ex.ColorWidth || col.B != ex.ColorWidth {
		t.Error("val err")
	}
}

func Test_getVariableRange(t *testing.T) {
	a, b := getVariableRange(1, 5, 2)
	if a != 0 || b != 2 {
		t.Errorf("range err%d:%d", a, b)
	}
}

func Test_VariableThreshold(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	//	make test image
	icom.ImageLoop(img, func(x int, y int) {
		rgba.Set(x, y, color.RGBA64{0, 0, 0, 1})
	})

	doFunc := VariableThreshold(10, 5, false)
	col := doFunc(rgba, 1, 1)
	if col.R != 0 || col.G != 0 || col.B != 0 {
		t.Error("val err")
	}
}
