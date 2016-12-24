package effect

import (
	"image"
	"image/color"
)

type Effect interface {
	ConvertToMonochromeImage() image.Image
}
type effect struct {
	inputImage image.Image
}

//	NewEffect is create new effect object
func NewEffect(in image.Image) Effect {
	return &effect{in}
}

// func (u color.RGBA64) ReverseConcentration(r, g, b, a uint32) color.RGBA64 {
// 	var col color.RGBA64
// 	outR := float32(r) * 0.298912
// 	outG := float32(g) * 0.58611
// 	outB := float32(b) * 0.114478
// 	mono := uint16(outR + outG + outB)
// 	col.R = mono
// 	col.G = mono
// 	col.B = mono
// 	col.A = uint16(a)
// 	return col
// }

func (ef *effect) ConvertToMonochromeImage() image.Image {
	rect := ef.inputImage.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var col color.RGBA64
			// 座標(x,y)のR, G, B, α の値を取得
			r, g, b, a := ef.inputImage.At(x, y).RGBA()
			//それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
			outR := float32(r) * 0.298912
			outG := float32(g) * 0.58611
			outB := float32(b) * 0.114478
			mono := uint16(outR + outG + outB)
			col.R = mono
			col.G = mono
			col.B = mono
			col.A = uint16(a)
			rgba.Set(x, y, col)
		}
	}
	return rgba
}
