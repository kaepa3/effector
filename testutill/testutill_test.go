package testutill

import (
	"image"
	"testing"
)

func Test_CreateRndImg(t *testing.T) {
	img := CreateRndImg(5, 5)
	if false == sizeCheck(img, 5, 5) {
		t.Errorf("size err(%d:%d)", img.Bounds().Size().X, img.Bounds().Size().Y)
	}
	mae := uint32(0)
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r == g && g == b && r == mae {
				t.Error("val err")
			}
			mae = r
		}
	}
}

func Test_CreateImg(t *testing.T) {
	img := CreateImg(5, 5)
	if false == sizeCheck(img, 5, 5) {
		t.Error("size err")
	}
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				t.Error("val err")
			}
		}
	}
}

func sizeCheck(img image.Image, h, w int) bool {
	rect := img.Bounds()
	if rect.Size().X != w || rect.Size().Y != h {
		return false
	}
	return true

}
