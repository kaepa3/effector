package density

import (
	"image"
	"image/color"
	"testing"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/testutill"
)

func Test_ReverseDensityFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		_, _, _, a := img.At(x, y).RGBA()
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, uint16(a)})
	})
	col := ReverseDensityFunc(img, 0, 0)
	if col.R != ex.ColorWidth || col.G != ex.ColorWidth || col.B != ex.ColorWidth {
		t.Error("val err")
	}
}

func Test_MonochromeFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		_, _, _, a := img.At(x, y).RGBA()
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, uint16(a)})
	})
	icom.ImageLoop(img, func(x int, y int) {
		r, g, b, _ := img.At(x, y).RGBA()
		col := MonochromeFunc(img, 0, 0)
		mono := getMono(r, g, b)
		if mono != col.R || mono != col.G || mono != col.B {
			t.Error("error")
		}
	})
}

func Test_FourToneFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		_, _, _, a := img.At(x, y).RGBA()
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, uint16(a)})
	})
	doFunc := FourToneFunc(4)
	col := doFunc(img, 0, 0)
	if col.R != 0 || col.G != 0 || col.B != 0 {
		t.Error("val err")
	}
}

func Test_ChangeSizeKinFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		_, _, _, a := img.At(x, y).RGBA()
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, uint16(a)})
	})
	doFunc := ChangeSizeKinFunc(rgba, 1, 1)
	icom.ImageLoop(rgba, func(x int, y int) {
		r, g, b, _ := img.At(x, y).RGBA()
		col := doFunc(img, 0, 0)
		if uint16(r) != col.R || uint16(g) != col.G || uint16(b) != col.B {
			t.Error("size change err")
		}
	})
}
func Test_ChangeSizeSenFunc(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(rgba, func(x int, y int) {
		_, _, _, a := img.At(x, y).RGBA()
		val := uint16(x + y)
		rgba.Set(x, y, color.RGBA64{val, val, val, uint16(a)})
	})
	doFunc := ChangeSizeSenFunc(rgba, 1, 1)
	icom.ImageLoop(rgba, func(x int, y int) {
		r, g, b, _ := img.At(x, y).RGBA()
		col := doFunc(img, 0, 0)
		if uint16(r) != col.R || uint16(g) != col.G || uint16(b) != col.B {
			t.Error("size change err")
		}
	})
}
