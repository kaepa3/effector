package effector

import (
	"math"
	"testing"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/testutill"
)

//ReverseConcentrationのテスト
func Test_ReverseConcentration(t *testing.T) {
	// 10*10の画像を作る
	img := testutill.CreateImg(10, 10)
	buf := ReverseConcentration(img)
	icom.ImageLoop(img, func(x, y int) {
		r1, g1, b1, a1 := img.At(x, y).RGBA()
		r2, g2, b2, a2 := buf.At(x, y).RGBA()
		if a1 != a2 {
			t.Errorf("a val error %d:%d (%d:%d)", x, y, a1, a2)
		}
		ary := []uint32{r1 + r2, g1 + g2, b1 + b2}
		for _, v := range ary {
			if ex.ColorWidth != v {
				t.Errorf("color error %d:%d (%x)%x", x, y, v, ex.ColorWidth)
			}
		}
	})
}

//Monochromeのテスト
func Test_Monochrome(t *testing.T) {
	// 10*10の画像を作る
	img := testutill.CreateRndImg(10, 10)
	buf := Monochrome(img)
	//	試験関数
	getMono := func(r, g, b float64) uint16 {
		return uint16(r*ex.RedNTSC + g*ex.GreenNTSC + b*ex.BlueNTSC)
	}
	const ThresholdValue = ex.ColorWidth / 100
	icom.ImageLoop(img, func(x, y int) {
		r1, g1, b1, a1 := img.At(x, y).RGBA()
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
	img := testutill.CreateRndImg(10, 10)
	buf := FourTone(img)
	countList := make([]uint32, 0, 6)
	icom.ImageLoop(img, func(x, y int) {
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
			if add == false {
				countList = append(countList, val)
			}
		}
	})
	if len(countList) != 4 {
		t.Errorf("count error:%d", len(countList))
	}
}
