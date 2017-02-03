package effector

import (
	"image"

	"github.com/kaepa3/effector/density"
	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/histogram"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/spatial"
	"github.com/kaepa3/effector/twoval"
)

//ReverseDensity は与えられた画像の濃度を逆転する
func ReverseDensity(img image.Image) image.Image {
	return icom.SimpleEffect(img, density.ReverseDensityFunc)
}

//Monochrome は与えられた画像をモノクロにする
func Monochrome(img image.Image) image.Image {
	return icom.SimpleEffect(img, density.MonochromeFunc)
}

//FourTone は与えられた画像を4階調に変更する
func FourTone(img image.Image, tones int) image.Image {
	doFunc := density.FourToneFunc(tones)
	return icom.SimpleEffect(img, doFunc)
}

//ChangeSizeKin は与えられた画像を最小近傍方を使用して拡大・縮小する。
func ChangeSizeKin(img image.Image, xRatio, yRatio float64) image.Image {
	rect := img.Bounds()
	width := ex.Times(rect.Size().X, xRatio)
	height := ex.Times(rect.Size().Y, yRatio)
	newRect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(newRect)

	doFunc := density.ChangeSizeKinFunc(img, xRatio, yRatio)
	return icom.SimpleEffect(rgba, doFunc)
}

//ChangeSizeSen は与えられた画像を線形補間法を使用して拡大・縮小する。
func ChangeSizeSen(img image.Image, xRatio, yRatio float64) image.Image {
	rect := img.Bounds()
	width := ex.Times(rect.Size().X, xRatio)
	height := ex.Times(rect.Size().Y, yRatio)
	newRect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA64(newRect)

	doFunc := density.ChangeSizeSenFunc(img, xRatio, yRatio)
	return icom.SimpleEffect(rgba, doFunc)
}

// LinearDensity は与えられた画像を線形濃度変換を使用して濃度変換する。
func LinearDensity(img image.Image, levelA, levelB uint16) image.Image {
	doFunc := density.LinearDensityFunc(levelA, levelB)
	return icom.SimpleEffect(img, doFunc)
}

// UnlinearDensity は与えられた画像を非線形濃度変換を使用して濃度変換する。
func UnlinearDensity(img image.Image, gamma float64) image.Image {
	doFunc := density.UnlinearDensityFunc(gamma)
	return icom.SimpleEffect(img, doFunc)
}

// ContrastImprovement は与えられた画像のコントラスト改善を実施する
func ContrastImprovement(img image.Image, gamma float64) image.Image {
	doFunc := density.ContrastImprovementFunc(gamma)
	return icom.SimpleEffect(img, doFunc)
}

// AverageHistogram は与えられた画像のヒストグラム平均化を実施する。
func AverageHistogram(img image.Image) image.Image {
	doFunc := density.AverageHistogramFunc(img)
	return icom.SimpleEffect(img, doFunc)
}

// Histogram は与えられた画像のヒストグラムを作成する
func Histogram(img image.Image, title, xLabel, yLabel string) image.Image {
	return histogram.Output(img, title, xLabel, yLabel)
}

//AverageFilter は与えられた画像に対して積和演算から平均値を採用する。
func AverageFilter(img image.Image, centerWeight float64) image.Image {
	doFunc := spatial.AverageFunc(centerWeight)
	return icom.SimpleEffect(img, doFunc)
}

//MedianFilter は与えられた画像に対して積和演算から中央値を採用する。
func MedianFilter(img image.Image) image.Image {
	return icom.SimpleEffect(img, spatial.MedianFunc)
}

//PrewittFilter は与えられた画像に対してPrewittフィルターを適用する。
func PrewittFilter(img image.Image) image.Image {
	return icom.SimpleEffect(img, spatial.PrewittFunc)
}

//VirticalLineFilter は与えられた画像に対してPrewittフィルターを適用する。
func VirticalLineFilter(img image.Image, weight float64, flg bool) image.Image {
	doFunc := spatial.VirticalLineFunc(weight, flg)
	return icom.SimpleEffect(img, doFunc)
}

//HorizontalLineFilter は与えられた画像に対してPrewittフィルターを適用する。
func HorizontalLineFilter(img image.Image, weight float64, flg bool) image.Image {
	doFunc := spatial.HorizontalLineFunc(weight, flg)
	return icom.SimpleEffect(img, doFunc)
}

//LaplacianFilter はラプラシアンフィルタのかかった画像を返す。
func LaplacianFilter(img image.Image) image.Image {
	return icom.SimpleEffect(img, spatial.LaplacianFunc)
}

//SharpeningFilter は鋭角化した画像を返す
func SharpeningFilter(img image.Image) image.Image {
	return icom.SimpleEffect(img, spatial.SharpeningFunc)
}

//StaticThreshold は固定式位置で２値化する。
func StaticThreshold(img image.Image, threshold uint32) image.Image {
	doFunc := twoval.StaticThreshold(threshold)
	return icom.SimpleEffect(img, doFunc)
}

//StaticThreshold は固定式位置で２値化する。
func VariableThreshold(img image.Image, size, width int, isBlack bool) image.Image {
	doFunc := twoval.VariableThreshold(size, width, isBlack)
	return icom.SimpleEffect(img, doFunc)
}

//Thinning は固定式位置で２値化する。
func Thinning(img image.Image) image.Image {
	return twoval.Thinning(img)
}

//BoundaryTracking は固定式位置で２値化する。
func BoundaryTracking(img image.Image) image.Image {
	return twoval.BoundaryTracking(img)
}
