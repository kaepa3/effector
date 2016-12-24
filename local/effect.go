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

// 逆にする。
func (ef *effect) ReverseConcentration() image.Image {
	rect := ef.inputImage.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var col color.RGBA64
			// 座標(x,y)のR, G, B, α の値を取得
			r, g, b, a := ef.inputImage.At(x, y).RGBA()
			//反転する
			col.R = math.MaxUint16 - uint16(r)
			col.G = math.MaxUint16 - uint16(g)
			col.B = math.MaxUint16 - uint16(b)
			col.A = uint16(a)
			rgba.Set(x, y, col)
		}
	}
	return rgba
}

// モノクロにする
func (ef *effect) Monochrome() image.Image {
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

// 4階調にする
func (ef *effect) FourTone() image.Image {
	rect := ef.inputImage.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)

	// この処理での特殊な値
	//	4階調とする
	tone := 4
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// 座標(x,y)のR, G, B, α の値を取得
			r, g, b, a := ef.inputImage.At(x, y).RGBA()
			effectFunction := fourToneColor
			rgba.Set(x, y, effectFunction(r, g, b, a, tone))
		}
	}
	return rgba
}

// カラーの場合の階調変更
func fourToneColor(r, g, b, a uint32, tone int) color.RGBA64 {
	var col color.RGBA64
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

// モノクロ画像判定
func isMonochromeImage(r, g, b, a uint32) bool {
	retVal := false
	if r == g && g == b {
		retVal = true
	}
	return retVal
}
