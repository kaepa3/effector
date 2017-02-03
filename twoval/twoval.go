package twoval

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"

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

//Thinning is
func Thinning(srcImg image.Image) image.Image {
	//白を出すだけの関数
	white := func(aVal uint32) color.Color {
		return color.RGBA64{ex.ColorWidth, ex.ColorWidth, ex.ColorWidth, uint16(aVal)}
	}
	//出力画像
	var workImg *image.RGBA
	for {
		workImg = image.NewRGBA(srcImg.Bounds())
		draw.Copy(workImg, image.Point{0, 0}, srcImg, srcImg.Bounds(), draw.Src, nil)
		counter := 0
		icom.ImageLoop(srcImg, func(x int, y int) {
			//判定を行う値かチェックする
			if false == isCheckValue(x, y, srcImg) {
				return
			}
			//周囲8の値を取得する
			srcAround, sum := getArounds(srcImg, x, y)
			workAround, _ := getArounds(workImg, x, y)
			//aの値を取得しておく
			_, _, _, a := workImg.At(x, y).RGBA()
			switch sum {
			case 0:
				workImg.Set(x, y, white(a))
			case 2, 3, 4, 5:
				if true == doesKeepConnect(srcAround, workAround) {
					workImg.Set(x, y, white(a))
				}
			}
			r, _, _, _ := workImg.At(x, y).RGBA()
			if r == ex.ColorWidth {
				counter++
			}
		})
		if counter == 0 {
			break
		}
		// 元画像に戻す
		srcImg = workImg
	}
	return workImg
}

func doesKeepConnect(srcAround, workAround Around) bool {
	flg := false
	if countConnect(srcAround) == 1 && countConnect(workAround) == 1 {
		var tmpArd Around
		isAll := true
		for i := 1; i < 5; i++ {
			if workAround[i] != 0 {
				tmpArd = srcAround
				tmpArd[i] = ex.ColorWidth
				if countConnect(tmpArd) != 1 {
					isAll = false
					break
				}
			}
		}
		if isAll == true {
			flg = true
		}
	}
	return flg
}
func countConnect(in Around) int {
	var count int
	//	これのおかげで一周する
	in[len(in)-1] = in[0]
	for i := 1; i < len(in); i++ {
		if in[i] == 0 && in[i-1] != 0 {
			count++
		}
	}
	return count
}

type Around [9]uint32

func getArounds(srcImg image.Image, x, y int) (Around, uint32) {
	keys := []struct {
		xKey int
		yKey int
	}{
		{x + 1, y},
		{x + 1, y - 1},
		{x, y - 1},
		{x - 1, y - 1},
		{x - 1, y},
		{x - 1, y + 1},
		{x, y + 1},
		{x + 1, y + 1},
	}
	var cnt uint32
	var out [9]uint32
	for i, c := range keys {
		r, _, _, _ := srcImg.At(c.xKey, c.yKey).RGBA()
		out[i] = r
		if r == 0 {
			cnt++
		}
	}
	return out, cnt
}

func isCheckValue(x, y int, srcImg image.Image) bool {
	r, _, _, _ := srcImg.At(x, y).RGBA()
	size := srcImg.Bounds().Size()
	if r == ex.ColorWidth {
		return false
	}
	if x == 0 || y == 0 {
		return false
	}
	if size.X == x || size.Y == y {
		return false
	}
	return true
}

func BoundaryTracking(srcImg image.Image) image.Image {
	workImg := image.NewRGBA(srcImg.Bounds())
	//キャンパスを真っ白にする。
	icom.ImageLoop(workImg, func(x int, y int) {
		workImg.Set(x, y, color.RGBA{math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8})
	})
	//探索開始
	icom.ImageLoop(srcImg, func(x int, y int) {
		if true == isCheckValue(x, y, srcImg) {
			r, _, _, _ := srcImg.At(x, y).RGBA()
			wR, _, _, _ := workImg.At(x, y).RGBA()
			//	元画像が黒かつ、ワークが白（既に走査済みの線の可能性もこれで排除）
			if r == 0 && wR != 0 {
				//外枠か内枠かを操作する
				code := getCode(x, y, srcImg)
				if code != -1 {
					imgChase(x, y, code, srcImg, workImg)
				}
			}
		}
	})
	return workImg
}

func imgChase(x, y, code int, srcImg image.Image, workImg *image.RGBA) {
	xStart := x
	yStart := y
	fhase := func(xPos, yPos, a, b int) int {
		r, _, _, _ := srcImg.At(xPos, yPos).RGBA()
		if r == 0 {
			return a
		}
		return b
	}
	var x2, y2 int
	//	探索を開始して最初の位置に戻ってくるまで
	for x2 != x || y2 != y {
		x2 = xStart
		y2 = yStart
		switch code {
		case 0: //基準点から下
			y2++
			code = fhase(x2, y2, 6, 2)
		case 2: //基準点から右
			x2++
			code = fhase(x2, y2, 0, 4)
		case 4: //基準点から上
			y2--
			code = fhase(x2, y2, 2, 6)
		case 6: //基準点から左
			x2--
			code = fhase(x2, y2, 4, 0)
		}
		r, _, _, a := srcImg.At(x2, y2).RGBA()
		if r == 0 {
			workImg.Set(x2, y2, color.RGBA{0, 0, 0, uint8(a)})
			xStart = x2
			yStart = y2
		}
	}
}

func getCode(x, y int, img image.Image) int {
	r, _, _, _ := img.At(x-1, y).RGBA()
	if r != 0 {
		return 0
	}
	r, _, _, _ = img.At(x+1, y).RGBA()
	if r != 0 {
		return 4
	}
	return -1
}
