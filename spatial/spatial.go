package spatial

import (
	"image"
	"image/color"
)

func AverageFunc(in image.Image, x, y int) color.RGBA64 {
	// 恥はやらない。
	if x == 0 || y == 0 || x == in.Bounds().Size().X || y == in.Bounds().Size().Y {
		r, g, b, a := in.At(x, y).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
	const CenterWeight = 1.0
	var sumR, sumG, sumB uint32
	var valA uint16
	//3*3のサイズを全部足す
	for pX := x - 1; pX <= x+1; pX++ {
		for pY := y - 1; pY <= y+1; pY++ {
			r, g, b, a := in.At(pX, pY).RGBA()
			weight := 1.0
			if pX == x && pY == y {
				valA = uint16(a)
				weight = CenterWeight
			}
			sumR += uint32(float64(r) * weight)
			sumG += uint32(float64(g) * weight)
			sumB += uint32(float64(b) * weight)
		}
	}
	valR := uint16(float64(sumR) / (8 + CenterWeight))
	valG := uint16(float64(sumG) / (8 + CenterWeight))
	valB := uint16(float64(sumB) / (8 + CenterWeight))
	return color.RGBA64{valR, valG, valB, valA}
}
