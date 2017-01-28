package twoval

import (
	"image"
	"image/color"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
)

func StaticThreshold(th uint32) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		r, _, _, a := img.At(x, y).RGBA()
		val := ex.ColorWidth
		if r < th {
			val = 0
		}
		return color.RGBA64{uint16(val), uint16(val), uint16(val), uint16(a)}
	}
}

type VariableData struct {
	yPos     int
	grayBase int
}

//VariableThreshold is
func VariableThreshold(size, width int, isBlack bool) icom.EffectFunc {
	data := VariableData{-1, 0}
	return func(img image.Image, x, y int) color.RGBA64 {
		//	y軸が変わったら初期化するだけ
		createGrayBase(&data, y, isBlack)
		point := img.Bounds().Size()
		//平均を取る画素の範囲を作成する
		y1, y2 := getVariableRange(y, point.Y, size)
		x1, x2 := getVariableRange(x, point.X, size)
		ave := getAverage(img, x1, x2, y1, y2)
		var th int
		// 前回の変化した種類によって基準値を変化させる。
		if data.grayBase == ex.ColorWidth {
			th = ave - width
		} else {
			th = ave + width
		}
		//平均値より明るいか？
		r, _, _, a := img.At(x, y).RGBA()
		v := uint16(0)
		if r >= uint32(th) {
			v = ex.ColorWidth
		}
		data.grayBase = int(v)
		return color.RGBA64{v, v, v, uint16(a)}
	}
}
func createGrayBase(data *VariableData, yPos int, isBlack bool) {
	if data.yPos != yPos {
		data.grayBase = 0
		if isBlack == true {
			data.grayBase = ex.ColorWidth
		}
	}
	data.yPos = yPos
}

func getAverage(img image.Image, xFrom, xTo, yFrom, yTo int) int {
	var sum, count int64
	for i := xFrom; i <= xTo; i++ {
		for j := yFrom; j <= yTo; j++ {
			r, _, _, _ := img.At(i, j).RGBA()
			sum += int64(r)
			count++
		}
	}
	return int(sum / count)
}

func getVariableRange(pos, max, siz int) (int, int) {
	v1 := pos - siz/2
	v2 := pos + siz/2
	//	最初から現在値＋規定値の半分だけ
	if pos < siz/2 {
		//	pos-siz/2最初突破の場合
		v1 = 0
	} else if (max - 1 - pos) < siz/2 {
		//	pos〜最大オーバーの場合
		v2 = max - 1
	}
	return v1, v2
}
