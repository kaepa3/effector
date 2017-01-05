package density

import (
	"image"
	"image/color"
	"math"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
)

func ReverseDensityFunc(img image.Image, x, y int) color.RGBA64 {
	var col color.RGBA64
	r, g, b, a := img.At(x, y).RGBA()
	col.R = ex.ColorWidth - uint16(r)
	col.G = ex.ColorWidth - uint16(g)
	col.B = ex.ColorWidth - uint16(b)
	col.A = uint16(a)
	return col
}

func getMono(r, g, b uint32) uint16 {
	return uint16(float64(r)*ex.RedNTSC + float64(g)*ex.GreenNTSC + float64(b)*ex.BlueNTSC)
}

func MonochromeFunc(img image.Image, x, y int) color.RGBA64 {
	var col color.RGBA64
	r, g, b, a := img.At(x, y).RGBA()
	//それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
	mono := getMono(r, g, b)
	col.R = mono
	col.G = mono
	col.B = mono
	col.A = uint16(a)
	return col
}
func FourToneFunc(tones int) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		var col color.RGBA64
		r, g, b, a := img.At(x, y).RGBA()
		//	4階調とする
		z1 := uint16(math.MaxUint16 / (tones))
		z2 := uint16(math.MaxUint16 / (tones - 1))
		vals := []uint32{r, g, b}
		ptr := []*uint16{&col.R, &col.G, &col.B}
		//	計算する
		for i, v := range vals {
			*ptr[i] = (uint16(v) / z1) * z2
		}
		col.A = uint16(a)
		return col
	}
}

func ChangeSizeKinFunc(src image.Image, xRatio, yRatio float64) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		var col color.RGBA64
		r, g, b, a := src.At(ex.Division(x, xRatio), ex.Division(y, yRatio)).RGBA()
		col.R = uint16(r)
		col.G = uint16(g)
		col.B = uint16(b)
		col.A = uint16(a)
		return col
	}
}

func ChangeSizeSenFunc(src image.Image, xRatio, yRatio float64) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		//	計測点（パラメータ）を作る
		rect := src.Bounds()
		j1, j2, q := createSenParam(rect.Size().Y, y, yRatio)
		i1, i2, p := createSenParam(rect.Size().X, x, xRatio)
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
			r, g, b, a := src.At(v.X, v.Y).RGBA()
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
}

func createSenParam(inSize, pos int, ratio float64) (int, int, float64) {
	v1 := ex.Division(pos, ratio)
	v2 := v1 + 1
	if v2 > inSize-1 {
		v2 = inSize - 1
	}
	v3 := float64(pos)/ratio - float64(v1)
	return v1, v2, v3
}

func LinearDensityFunc(levA, levB uint16) icom.EffectFunc {
	con := func(val uint32) uint16 {
		v := math.MaxUint16 * (float64((uint16(val) - levA)) / float64((levB - levA)))
		if v > math.MaxUint16 {
			v = math.MaxUint16
		} else if v < 0 {
			v = 0
		}
		return uint16(v)
	}
	return func(img image.Image, x, y int) color.RGBA64 {
		r, g, b, a := img.At(x, y).RGBA()
		return color.RGBA64{con(r), con(g), con(b), uint16(a)}
	}
}

func UnlinearDensityFunc(gamma float64) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		var col color.RGBA64
		r, g, b, a := img.At(x, y).RGBA()
		col.R = unlinerCon(r, gamma)
		col.G = unlinerCon(g, gamma)
		col.B = unlinerCon(b, gamma)
		col.A = uint16(a)
		return col
	}
}
func unlinerCon(val uint32, gamma float64) uint16 {
	v := math.MaxUint16 * math.Pow(float64(val)/float64(math.MaxUint16), gamma)
	if v > math.MaxUint16 {
		v = math.MaxUint16
	} else if v < 0 {
		v = 0
	}
	return uint16(v)
}

func ContrastImprovementFunc(gamma float64) icom.EffectFunc {
	con := func(val uint32) uint16 {
		maxUint := float64(ex.ColorWidth)
		var v uint16
		if val > math.MaxUint16/2 {
			v = uint16((maxUint / 2.0) * (1.0 + math.Sqrt((2.0*float64(val)-maxUint)/maxUint)))
		} else {
			// 非線形濃度変換
			v = unlinerCon(val, gamma)
		}
		if v > math.MaxUint16 {
			v = math.MaxUint16
		} else if v < 0 {
			v = 0
		}
		return uint16(v)
	}
	return func(img image.Image, x, y int) color.RGBA64 {
		r, g, b, a := img.At(x, y).RGBA()
		return color.RGBA64{con(r), con(g), con(b), uint16(a)}
	}
}

func AverageHistogramFunc(img image.Image) icom.EffectFunc {
	//ヒストグラムの取得
	_, _, _, hisL := icom.MakeHistogramData(img)
	//lookupTBLの作成
	luTbl := createLookupTable(hisL, img.Bounds())

	return func(img image.Image, x, y int) color.RGBA64 {
		r, g, b, a := img.At(x, y).RGBA()
		return color.RGBA64{luTbl[r], luTbl[g], luTbl[b], uint16(a)}
	}
}
func createLookupTable(his [ex.ColorWidthAryMax]uint16, rect image.Rectangle) (table [ex.ColorWidthAryMax]uint16) {
	var sum, val uint16
	//平均画素数
	average := uint16(((rect.Size().X * rect.Size().Y) / ex.ColorWidth) + 1)
	//平均値の計算
	for i, v := range his {
		sum += v
		val += sum / average
		sum %= average
		table[i] = val
	}
	return
}
