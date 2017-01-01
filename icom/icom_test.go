package icom

import (
	"image"
	"image/color"
	"testing"

	"github.com/kaepa3/effector/testutill"
)

//	imageLoopのテスト
func Test_ImageLoop(t *testing.T) {
	// 10*10の画像を作る
	in := testutill.CreateImg(10, 10)
	count := 0
	ImageLoop(in, func(x, y int) {
		count++
	})
	if count != 100 {
		t.Errorf("call count Error:%d", count)
	}
}

func Test_MakeHistogramData(t *testing.T) {
	// 10*10の画像を作る
	buf := testutill.CreateImg(10, 10)
	testImg := image.NewRGBA64(buf.Bounds())
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			var col color.RGBA64
			col.R = uint16((x * 10) + y)
			col.G = uint16(col.R % 10)
			col.B = uint16(25)
			col.A = 0xffff
			testImg.Set(x, y, col)
		}
	}

	rD, gD, bD, _ := MakeHistogramData(testImg)

	for i, v := range rD {
		// R
		if i >= 100 {
			break
		}
		//R
		if v == 0 {
			t.Errorf("r err:%d(%d)", v, i)
		}

		//G
		if i < 10 {
			if gD[i] == 0 {
				t.Errorf("g err:%d(%d)", gD[i], i)
			}
		} else {
			if gD[i] != 0 {
				t.Errorf("g err:%d(%d)", gD[i], i)
			}
		}
		//B
		if i == 25 {
			if bD[i] != 100 {
				t.Errorf("b err:%d(%d)", bD[i], i)
			}
		} else {
			if bD[i] != 0 {
				t.Errorf("b err:%d(%d)", bD[i], i)
			}
		}
	}

}
