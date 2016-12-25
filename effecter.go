package effecter

import (
	"image"
	"image/color"
	"math"
)

//Effect is image effect interface
type Effect interface {
	Monochrome() image.Image
	ReverseConcentration() image.Image
	FourTone() image.Image
	ChangeSizeKin() image.Image
	ChangeSizeSen() image.Image
}
type effect struct {
	inputImage image.Image
}

//NewEffect is create new effect object
func NewEffect(in image.Image) Effect {
	return &effect{in}
}

func (ef *effect) imageLoop(rect image.Rectangle, effectFunction func(x, y int) color.RGBA64) image.Image {
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// 座標(x,y)のR, G, B, α の値を取得
			col := effectFunction(x, y)
			rgba.Set(x, y, col)
		}
	}
	return rgba
}

// 逆にする。
func (ef *effect) ReverseConcentration() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		col.R = math.MaxUint16 - uint16(r)
		col.G = math.MaxUint16 - uint16(g)
		col.B = math.MaxUint16 - uint16(b)
		col.A = uint16(a)
		return col
	})
}

// モノクロにする
func (ef *effect) Monochrome() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
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
		return col
	})
}

// 4階調にする
func (ef *effect) FourTone() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
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
	})
}

var up = func(val int, rat float64) int { return int(float64(val) * rat) }
var down = func(val int, rat float64) int { return int(float64(val) / rat) }

// 最近傍方
func (ef *effect) ChangeSizeKin() image.Image {

	// とりあえず正方形にして、1/2倍にする
	ratio := 0.9
	rect := ef.inputImage.Bounds()
	width := up(rect.Size().X, ratio)
	height := width
	newRect := image.Rect(0, 0, width, height)
	yRatio := float64(height) / float64(rect.Size().Y)

	return ef.imageLoop(newRect, func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(down(x, ratio), down(y, yRatio)).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	})
}

func (ef *effect) ChangeSizeSen() image.Image {

	// とりあえず正方形にして、1/2倍にする
	ratio := 0.9
	rect := ef.inputImage.Bounds()
	width := up(rect.Size().X, ratio)
	height := width
	newRect := image.Rect(0, 0, width, height)
	yRatio := float64(height) / float64(rect.Size().Y)

	return ef.imageLoop(newRect, func(x, y int) color.RGBA64 {
		return ef.senning(x, y, ratio, yRatio)
	})
}

func (ef *effect) senning(x, y int, xRatio, yRatio float64) color.RGBA64 {
	// 比較のための４点とその位置の比を求める関数
	createParam := func(inSize, outSize int, ratio float64) (int, int, float64) {
		v1 := down(outSize, ratio)
		v2 := v1 + 1
		if v2 > inSize-1 {
			v2 = inSize - 1
		}
		v3 := float64(outSize)/ratio - float64(v1)
		return v1, v2, v3
	}
	//	計測点（パラメータ）を作る
	point := ef.inputImage.Bounds().Size()
	j1, j2, q := createParam(point.Y, y, yRatio)
	i1, i2, p := createParam(point.X, x, xRatio)
	positions := [4]struct {
		X int
		Y int
	}{
		{i1, j1},
		{i2, j1},
		{i1, j2},
		{i2, j2},
	}
	// 4点の値を取得する
	var valsR, valsG, valsB, valsA [len(positions)]float64
	for i, v := range positions {
		r, g, b, a := ef.inputImage.At(v.X, v.Y).RGBA()
		valsR[i] = float64(r)
		valsG[i] = float64(g)
		valsB[i] = float64(b)
		valsA[i] = float64(a)
	}
	//各要素の値を求める関数
	con := func(val [4]float64) uint16 {
		return uint16((1.0-q)*((1.0-p)*val[0]+p*val[1]) + q*((1-p)*val[2]+p*val[3]))
	}
	return color.RGBA64{con(valsR), con(valsG), con(valsB), con(valsA)}
}
