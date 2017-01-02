package spatial

import (
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/kaepa3/effector/icom"
)

func AverageFunc(centerWeight float64) icom.EffectFunc {
	return func(img image.Image, x, y int) color.RGBA64 {
		if x == 0 || y == 0 || x == img.Bounds().Size().X || y == img.Bounds().Size().Y {
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

var counter uint

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
	sort.Sort(Gasos(aryR))
	sort.Sort(Gasos(aryG))
	sort.Sort(Gasos(aryB))
	sort.Sort(Gasos(aryA))
	if counter < 10 {
		log.Print(aryR)
	}
	counter++
	return color.RGBA64{aryR[4], aryG[4], aryB[4], aryA[4]}
}
