package twoval

import (
	"image"
	"log"
	"math"

	"golang.org/x/image/draw"

	"github.com/kaepa3/effector/icom"
	"github.com/llgcode/draw2d/draw2dimg"
)

const FormulaMax = 10

func BoundaryTracking(srcImg image.Image) image.Image {
	var coefficient [FormulaMax][FormulaMax]float64
	var constant [FormulaMax]float64
	//探索開始
	icom.ImageLoop(srcImg, func(x int, y int) {
		r, _, _, _ := srcImg.At(x, y).RGBA()
		if r == 0 {
			for i, v := range coefficient {
				constant[i] = math.Pow(x, i) * y
				for j, _ := range v {
					coefficient[i][j] += math.Pow(x, i+j)
				}
			}

		}
	})
	gauss(coefficient, constant)
	workImg := image.NewRGBA(srcImg.Bounds())
	draw.Copy(workImg, image.Point{0, 0}, srcImg, srcImg.Bounds(), draw.Src, nil)
	size := srcImg.Bounds().Size()
	var points [FormulaMax]image.Point
	for x := 0; x < size.X; x++ {
		yy := constant[0]
		for k := 0; k < FormulaMax; k++ {
			yy += constant[k] * math.Pow(x, k)
			points[k].X = x
			points[k].Y = yy
		}
	}
	workImg := image.NewRGBA(srcImage.Bounds())
	draw.Copy(workImg, image.Point{0, 0}, srcImage, srcImage.Bounds(), draw.Src, nil)
	gc := draw2dimg.NewGraphicContext(workImg)
	for _, v := range points {
		gc.MoveTo(v.X, v.Y)
		gc.LineTo(xEnd, yEnd)
	}
	gc.DrawImage(workImg)

	return workImg
}

func gauss(coefficient *[][]float64, constant *[]float64) {
	for k, _ := range coefficient {
		if false == pivot(FormulaMax+1, i, coefficient, constant) {
			log.Println("no answer")
			return
		}
		akk := coefficient[k][k]
		for i := k + 1; i < len(constant)+1; i++ {
			p1 := coefficient[i][k] / akk
			constant[i] -= p1 * constant[k]
		}
	}
	constant[FormulaMax-1] /= coefficient[FormulaMax-1][FormulaMax-1]
	for i := FormulaMax - 2; i >= 0; i-- {
		s := 0.0
		for j = i + 1; j < FormulaMax+1; j++ {
			s += coefficient[i][j] * constant[j]
			constant[i] = (constant[i] - s) / coefficient[i][i]
		}
	}
}

func pivot(idx int, coefficient *[][]float64, constant *[]float64) bool {
	piv := 0.0
	kk := idx
	for i := 0; i < FormulaMax+1; i++ {
		if piv < math.Abs(coefficient[i][idx]) {
			piv = math.Abs(coefficient[i][idx])
			kk = i
		}
	}
	if piv == 0.0 {
		return false
	}
	if kk != idx {
		ch := 0.0
		for j := k; j < FormulaMax+1; j++ {
			ch = coefficient[k][j]
			coefficient[k][j] = coefficient[kk][j]
			coefficient[kk][j] = ch
		}
		ch = constant[k]
		constant[k] = constant[kk]
		constant[kk] = ch
	}
}
