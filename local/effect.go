package effect

import (
	"image"
	"image/color"
	"math"
)

type Effect interface {
	Monochrome() image.Image
	ReverseConcentration() image.Image
	FourTone() image.Image
}
type effect struct {
	inputImage image.Image
}

//	NewEffect is create new effect object
func NewEffect(in image.Image) Effect {
	return &effect{in}
}

func (ef *effect) imageLoop(effectFunction func(r, g, b, a uint32) color.RGBA64) image.Image {
	rect := ef.inputImage.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// 座標(x,y)のR, G, B, α の値を取得
			rgba.Set(x, y, effectFunction(ef.inputImage.At(x, y).RGBA()))
		}
	}
	return rgba
}

// 逆にする。
func (ef *effect) ReverseConcentration() image.Image {
	return ef.imageLoop(reverseConcentrationColor)
}

// カラーの場合の階調変更
func reverseConcentrationColor(r, g, b, a uint32) color.RGBA64 {
	var col color.RGBA64
	//反転する
	col.R = math.MaxUint16 - uint16(r)
	col.G = math.MaxUint16 - uint16(g)
	col.B = math.MaxUint16 - uint16(b)
	col.A = uint16(a)
	return col
}

// モノクロにする
func (ef *effect) Monochrome() image.Image {
	return ef.imageLoop(monochromeColor)
}

// カラーの場合の階調変更
func monochromeColor(r, g, b, a uint32) color.RGBA64 {
	var col color.RGBA64
	//それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
	outR := float32(r) * 0.298912
	outG := float32(g) * 0.58611
	outB := float32(b) * 0.114478
	mono := uint16(outR + outG + outB)
	col.R = mono
	col.G = mono
	col.B = mono
	col.A = uint16(a)
	return col
}

// 4階調にする
func (ef *effect) FourTone() image.Image {
	return ef.imageLoop(fourToneColor)
}

// カラーの場合の階調変更
func fourToneColor(r, g, b, a uint32) color.RGBA64 {
	var col color.RGBA64
	//	4階調とする
	tone := 4
	z1 := uint16(math.MaxUint16 / (tone))
	z2 := uint16(math.MaxUint16 / (tone - 1))
	vals := []uint32{r, g, b}
	ptr := []*uint16{&col.R, &col.G, &col.B}
	//	計算する
	for i, v := range vals {
		*ptr[i] = (uint16(v) / z1) * z2
	}
	col.A = uint16(a)
	return col
}
