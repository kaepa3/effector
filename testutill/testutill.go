package testutill

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/kaepa3/effector/ex"
)

func CreateBuffImage(x, y int, flg bool) image.Image {
	rect := image.Rectangle{image.Point{0, 0}, image.Point{x, y}}
	rgba := image.NewRGBA(rect)
	if flg == true {
		rand.Seed(time.Now().UnixNano())
		for x := 0; x < 5; x++ {
			for y := 0; y < 5; y++ {
				valR := uint16(rand.Int31n(ex.ColorWidth))
				valG := uint16(rand.Int31n(ex.ColorWidth))
				valB := uint16(rand.Int31n(ex.ColorWidth))
				rgba.Set(x, y, color.RGBA64{valR, valG, valB, 1})
			}
		}
	}
	return rgba
}

func CreateRndImg(x, y int) image.Image {
	return CreateBuffImage(x, y, true)
}

func CreateImg(x, y int) image.Image {
	return CreateBuffImage(x, y, false)
}
