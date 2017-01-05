package spatial

import (
	"fmt"
	"image"
	"image/color"
	"localhost/effector/ex"
	"math"
	"sort"

	"github.com/kaepa3/effector/icom"
)

func isRect(x, y int, img image.Image) bool {
	if x == 0 || y == 0 || x == img.Bounds().Size().X || y == img.Bounds().Size().Y {
		return true
	}
	return false
}

//AverageFunc は画像XY近傍9Pxを取得し平均値を返す。
func AverageFunc(centerWeight float64) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		if isRect(x, y, img) == true {
			r, g, b, a := img.At(x, y).RGBA()
			return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
		}
		var sumR, sumG, sumB uint32
		var valA uint16
		//3*3のサイズを全部足す
		for pX := x - 1; pX <= x+1; pX++ {
			for pY := y - 1; pY <= y+1; pY++ {
				r, g, b, a := img.At(pX, pY).RGBA()
				weight := 1.0
				if pX == x && pY == y {
					valA = uint16(a)
					weight = centerWeight
				}
				sumR += uint32(float64(r) * weight)
				sumG += uint32(float64(g) * weight)
				sumB += uint32(float64(b) * weight)
			}
		}
		valR := uint16(float64(sumR) / (8 + centerWeight))
		valG := uint16(float64(sumG) / (8 + centerWeight))
		valB := uint16(float64(sumB) / (8 + centerWeight))
		return color.RGBA64{valR, valG, valB, valA}
	}
}

//Gasos は画素のソートのために使用するデータ型
type Gasos []uint16

func (s Gasos) Len() int {
	return len(s)
}

func (s Gasos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Gasos) Less(i, j int) bool {
	return s[i] < s[j]
}

//MedianFunc は画像のあるXY近傍9Pxを取得し５番目の値を返す。
func MedianFunc(img image.Image, x, y int) color.RGBA64 {
	aryR := make([]uint16, 10)
	aryG := make([]uint16, 10)
	aryB := make([]uint16, 10)
	aryA := make([]uint16, 10)
	count := 0
	//3*3のサイズを全部足す
	for pX := x - 1; pX <= x+1; pX++ {
		for pY := y - 1; pY <= y+1; pY++ {
			r, g, b, a := img.At(pX, pY).RGBA()
			aryR[count] = uint16(r)
			aryG[count] = uint16(g)
			aryB[count] = uint16(b)
			aryA[count] = uint16(a)
			count++
		}
	}
	// ソートする
	sort.Sort(Gasos(aryR))
	sort.Sort(Gasos(aryG))
	sort.Sort(Gasos(aryB))
	sort.Sort(Gasos(aryA))

	return color.RGBA64{aryR[4], aryG[4], aryB[4], aryA[4]}
}

//PrewittFunc は画像のあるXY近傍9PxにPrewittフィルターをかけて返す。
func PrewittFunc(img image.Image, x, y int) color.RGBA64 {
	if isRect(x, y, img) == true {
		r, g, b, a := img.At(x, y).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
	// フィルタの情報の重み
	xFilter := [][]int{
		{-1, -1, -1},
		{0, 0, 0},
		{1, 1, 1},
	}
	yFilter := [][]int{
		{1, 0, -1},
		{1, 0, -1},
		{1, 0, -1},
	}
	//情報の集約（フィルタ込み）
	var rMd1, rMd2, gMd1, gMd2, bMd1, bMd2 int = 0, 0, 0, 0, 0, 0
	for pX := 0; pX < 3; pX++ {
		for pY := 0; pY < 3; pY++ {
			r, g, b, _ := img.At(x+(pX-1), y+(pY-1)).RGBA()
			rMd1 += xFilter[pX][pY] * int(r)
			rMd2 += yFilter[pX][pY] * int(r)
			gMd1 += xFilter[pX][pY] * int(g)
			gMd2 += yFilter[pX][pY] * int(r)
			bMd1 += xFilter[pX][pY] * int(b)
			bMd2 += yFilter[pX][pY] * int(r)
		}
	}
	// フィルタの絶対値取得
	getValue := func(md1, md2 int) uint16 {
		gaso := (math.Abs(float64(md1)) + math.Abs(float64(md2)))
		if gaso > ex.ColorWidth {
			gaso = ex.ColorWidth
		} else if gaso < 0 {
			gaso = 0
		}
		return uint16(gaso)
	}
	red := getValue(rMd1, rMd2)
	green := getValue(gMd1, gMd2)
	blue := getValue(bMd1, bMd2)
	_, _, _, a := img.At(x, y).RGBA()
	return color.RGBA64{red, green, blue, uint16(a)}
}

type Spatial struct {
	img     image.Image
	xFilter [3][3]float64
	yFilter [3][3]float64
	mdFunc  func(md1, md2 float64) uint16
}

//PrewittFunc は画像のあるXY近傍9PxにPrewittフィルターをかけて返す。
func VirticalLineFunc(img image.Image, x, y int) color.RGBA64 {
	xFilter := [3][3]float64{
		// {-1 / 2, 1, -1 / 2},
		// {-1 / 2, 1, -1 / 2},
		// {-1 / 2, 1, -1 / 2},
		{-0.5, 1, -0.5},
		{-0.5, 1, -0.5},
		{-0.5, 1, -0.5},
	}
	//情報の集約（フィルタ込み）
	var rMd1, gMd1, bMd1 float64 = 0, 0, 0
	for pX := 0; pX < 3; pX++ {
		for pY := 0; pY < 3; pY++ {
			r, g, b, _ := img.At(x+(pX-1), y+(pY-1)).RGBA()
			rMd1 = rMd1 + xFilter[pY][pX]*float64(r)
			gMd1 += xFilter[pY][pX] * float64(g)
			bMd1 += xFilter[pY][pX] * float64(b)
		}
	}
	mdFunc := func(md1 float64) uint16 {
		gaso := ex.ColorWidth - int(md1)
		if gaso > ex.ColorWidth {
			gaso = ex.ColorWidth
		} else if gaso < 0 {
			gaso = 0
		}
		return uint16(gaso)
	}
	_, _, _, a := img.At(x, y).RGBA()

	col := color.RGBA64{mdFunc(rMd1), mdFunc(gMd1), mdFunc(bMd1), uint16(a)}

	// フィルタの絶対値取得
	val := uint16(col.R * 1)
	fmt.Printf("(%d:%d) %d:%d:%d", x, y, col.R, col.G, col.B)
	return color.RGBA64{val, val, val, uint16(col.A)}
}

func (sp *Spatial) CreateColor(x, y int) color.RGBA64 {
	if isRect(x, y, sp.img) == true {
		r, g, b, a := sp.img.At(x, y).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
	//情報の集約（フィルタ込み）
	var rMd1, rMd2, gMd1, gMd2, bMd1, bMd2 float64 = 0, 0, 0, 0, 0, 0
	for pX := 0; pX < 3; pX++ {
		for pY := 0; pY < 3; pY++ {
			r, g, b, _ := sp.img.At(x+(pX-1), y+(pY-1)).RGBA()
			rMd1 += sp.xFilter[pX][pY] * float64(r)
			rMd2 += sp.yFilter[pX][pY] * float64(r)
			gMd1 += sp.xFilter[pX][pY] * float64(g)
			gMd2 += sp.yFilter[pX][pY] * float64(r)
			bMd1 += sp.xFilter[pX][pY] * float64(b)
			bMd2 += sp.yFilter[pX][pY] * float64(r)
		}
	}
	// フィルタの絶対値取得
	red := sp.mdFunc(rMd1, rMd2)
	green := sp.mdFunc(gMd1, gMd2)
	blue := sp.mdFunc(bMd1, bMd2)
	_, _, _, a := sp.img.At(x, y).RGBA()
	return color.RGBA64{red, green, blue, uint16(a)}
}
