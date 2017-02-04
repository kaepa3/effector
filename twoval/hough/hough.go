package hough

import (
	"image"
	"math"

	"github.com/kaepa3/effector/icom"
	"github.com/llgcode/draw2d/draw2dimg"
)

const MaxRadian = 180

func Hough(srcImg image.Image) image.Image {
	size := srcImg.Bounds().Size()
	maxRho := math.Sqrt2 * float64(size.X)
	//var pp [MAX_RADIAN][maxRho]int
	pp := make([][]int, int(maxRho))
	for i, _ := range pp {
		pp[i] = make([]int, MaxRadian)
	}
	//投票
	icom.ImageLoop(srcImg, func(x int, y int) {
		for r := 0; r < MaxRadian; r++ {
			xPos := float64(x - size.X/2)
			yPos := float64(size.Y/2 - y)
			rad := float64(r) * math.Pi / MaxRadian
			rho := xPos*math.Cos(rad) + yPos*math.Sin(rad) + maxRho/2
			pp[int(rad)][int(rho)]++
		}
	})
	//開票
	maxPP := 0
	rad0 := 0
	rho0 := 0
	for rho, v := range pp {
		for rad := range v {
			if maxPP < v[rad] {
				maxPP = v[rad]
				rad0 = rad
				rho0 = rho - (int(maxRho) / 2)
			}
		}
	}
	return createHough(float64(rad0), float64(rho0), srcImg)
}
func createHough(rad0, rho0 float64, srcImage image.Image) image.Image {
	var xx [2]int
	var yy [2]int

	pos := 0
	size := srcImage.Bounds().Size()
	rad := rad0 * math.Pi / MaxRadian

	xPos := float64(size.X / 2)
	yPos := -xPos/math.Tan(rad) + rho0/math.Sin(rad)
	if yPos >= float64(-size.Y/2) && yPos <= float64(size.Y/2) {
		xx[pos] = int(xPos)
		yy[pos] = int(yPos)
		pos++
	}

	xPos = float64(-size.X / 2)
	yPos = -xPos/math.Tan(rad) + rho0/math.Sin(rad)
	if yPos >= float64(-size.Y/2) && yPos <= float64(size.Y/2) {
		xx[pos] = int(xPos)
		yy[pos] = int(yPos)
		pos++
	}

	yPos = float64(size.Y / 2)
	xPos = rho0/math.Cos(rad) - yPos*math.Tan(rad)
	if xPos > float64(-size.X/2) && xPos < float64(size.X/2) {
		xx[pos] = int(xPos)
		yy[pos] = int(yPos)
		pos++
	}

	yPos = float64(-size.Y / 2)
	xPos = rho0/math.Cos(rad) - yPos*math.Tan(rad)
	if xPos > float64(-size.X/2) && xPos < float64(size.X/2) {
		xx[pos] = int(xPos)
		yy[pos] = int(yPos)
	}
	xStart := float64(xx[0] + size.X/2)
	xEnd := float64(xx[1] + size.X/2)
	yStart := float64(size.X/2 + yy[0])
	yEnd := float64(size.X/2 + yy[1])

	workImg := image.NewRGBA(srcImage.Bounds())
	gc := draw2dimg.NewGraphicContext(workImg)
	gc.MoveTo(xStart, yStart)
	gc.LineTo(xEnd, yEnd)
	gc.DrawImage(workImg)
	return workImg
}
