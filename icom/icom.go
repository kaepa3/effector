package icom

import (
	"image"

	"github.com/kaepa3/effector/ex"
)

func ImageLoop(img image.Image, doFunc func(x, y int)) {
	rect := img.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			doFunc(x, y)
		}
	}
}

func MakeHistogramData(img image.Image) (rD, gD, bD, lD [ex.ColorWidthAryMax]uint16) {
	ImageLoop(img, func(x, y int) {
		//画素をそれぞれ数える
		r, g, b, _ := img.At(x, y).RGBA()
		rD[r]++
		gD[g]++
		bD[b]++
		l := (3*r + 6*g + b) / 10
		lD[l]++
	})
	return
}
