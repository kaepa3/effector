package effecter

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"
	"time"
)

func createBuffImage(x, y int) image.Image {
	rect := image.Rectangle{image.Point{0, 0}, image.Point{10, 10}}
	rgba := image.NewRGBA(rect)
	var obj effect
	rand.Seed(time.Now().UnixNano())
	obj.inputImage = rgba
	obj.imageLoop(func(x, y int) {
		valR := uint16(rand.Int31n(ColorWidth))
		valG := uint16(rand.Int31n(ColorWidth))
		valB := uint16(rand.Int31n(ColorWidth))
		rgba.Set(x, y, color.RGBA64{valR, valG, valB, 0})
	})
	return rgba
}

//	imageLoopのテスト
func Test_imageLoop(t *testing.T) {
	// 10*10の画像を作る
	var buf effect
	buf.inputImage = createBuffImage(10, 10)
	count := 0
	buf.imageLoop(func(x, y int) {
		count++
	})
	if count != 100 {
		t.Errorf("call count Error:%d", count)
	}
}

//ReverseConcentrationのテスト
func Test_ReverseConcentration(t *testing.T) {
	// 10*10の画像を作る
	var obj effect
	obj.inputImage = createBuffImage(10, 10)
	buf := obj.ReverseConcentration()
	obj.imageLoop(func(x, y int) {
		r1, g1, b1, a1 := obj.inputImage.At(x, y).RGBA()
		r2, g2, b2, a2 := buf.At(x, y).RGBA()
		if a1 != a2 {
			t.Errorf("a val error %d:%d (%d:%d)", x, y, a1, a2)
		}
		ary := []uint32{r1 + r2, g1 + g2, b1 + b2}
		for _, v := range ary {
			if ColorWidth != v {
				t.Errorf("color error %d:%d (%x)%x", x, y, v, ColorWidth)
			}
		}
	})
}

//Monochromeのテスト
func Test_Monochrome(t *testing.T) {
	// 10*10の画像を作る
	var obj effect
	obj.inputImage = createBuffImage(10, 10)
	buf := obj.Monochrome()
	//	試験関数
	getMono := func(r, g, b float64) uint16 {
		return uint16(r*RedNTSC + g*GreenNTSC + b*BlueNTSC)
	}
	const ThresholdValue = ColorWidth / 100
	obj.imageLoop(func(x, y int) {
		r1, g1, b1, a1 := obj.inputImage.At(x, y).RGBA()
		r2, g2, b2, a2 := buf.At(x, y).RGBA()
		if a1 != a2 {
			t.Errorf("a val error %d:%d (%d:%d)", x, y, a1, a2)
		}
		if r2 != g2 || g2 != b2 {
			t.Errorf("a val error %d:%d:%d", r2, g2, b2)
		}
		mono := getMono(float64(r1), float64(g1), float64(b1))
		abs := math.Abs(float64(r2) - float64(mono))
		if abs >= ThresholdValue {
			t.Errorf("a val error %d:%d", r2, mono)
		}
	})
}

//FourToneのテスト
func Test_FourTone(t *testing.T) {
	// 10*10の画像を作る
	var obj effect
	obj.inputImage = createBuffImage(10, 10)
	buf := obj.FourTone()
	countList := make([]uint32, 0, 6)
	obj.imageLoop(func(x, y int) {
		r, g, b, _ := buf.At(x, y).RGBA()
		add := false
		vals := []uint32{r, g, b}
		for _, val := range vals {
			for _, v := range countList {
				if val == v {
					add = true
					break
				}
			}
			if add {
				countList = append(countList, val)
			}
		}
	})
	if len(countList) > 4 {
		t.Errorf("count error:%d", len(countList))
	}
}
