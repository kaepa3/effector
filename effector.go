package effector

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/spatial"
)

//ReverseConcentration は与えられた画像の濃度を逆転する
func ReverseConcentration(img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	icom.ImageLoop(img, func(x, y int) {
		var col color.RGBA64
		r, g, b, a := img.At(x, y).RGBA()
		col.R = ex.ColorWidth - uint16(r)
		col.G = ex.ColorWidth - uint16(g)
		col.B = ex.ColorWidth - uint16(b)
		col.A = uint16(a)
		rgba.Set(x, y, col)
	})
	return rgba
}

//Monochrome は与えられた画像をモノクロにする
func Monochrome(img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	getMono := func(r, g, b float64) uint16 {
		return uint16(r*ex.RedNTSC + g*ex.GreenNTSC + b*ex.BlueNTSC)
	}
	icom.ImageLoop(img, func(x, y int) {
		var col color.RGBA64
		r, g, b, a := img.At(x, y).RGBA()
		//それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
		mono := getMono(float64(r), float64(g), float64(b))
		col.R = mono
		col.G = mono
		col.B = mono
		col.A = uint16(a)
		rgba.Set(x, y, col)
	})
	return rgba
}

//FourTone は与えられた画像を4階調に変更する
func FourTone(img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())

	icom.ImageLoop(img, func(x, y int) {
		var col color.RGBA64
		r, g, b, a := img.At(x, y).RGBA()
		//	4階調とする
		tone := 4
		z1 := uint16(math.MaxUint16 / (tone))
		z2 := uint16(math.MaxUint16 / (tone - 1))
		vals := []uint32{r, g, b}
		ptr := []*uint16{&col.R, &col.G, &col.B}
		//	計算する
		for i, v := range vals {
			*ptr[i] = (uint16(v) / z1) * z2
		}
		col.A = uint16(a)
		rgba.Set(x, y, col)
	})
	return rgba
}

var up = func(val int, rat float64) int { return int(float64(val) * rat) }
var down = func(val int, rat float64) int { return int(float64(val) / rat) }

//ChangeSizeKin は与えられた画像を最小近傍方を使用して拡大・縮小する。
func ChangeSizeKin(img image.Image, xRatio, yRatio float64) image.Image {
	rect := img.Bounds()
	width := up(rect.Size().X, xRatio)
	height := up(rect.Size().Y, yRatio)
	newRect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(newRect)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := img.At(down(x, xRatio), down(y, yRatio)).RGBA()
			rgba.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return rgba
}

//ChangeSizeSen は与えられた画像を線形補間法を使用して拡大・縮小する。
func ChangeSizeSen(img image.Image, xRatio, yRatio float64) image.Image {
	rect := img.Bounds()
	width := up(rect.Size().X, xRatio)
	height := up(rect.Size().Y, yRatio)
	newRect := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA64(newRect)
	icom.ImageLoop(rgba, func(x int, y int) {
		con := senning(img, x, y, xRatio, yRatio)
		rgba.Set(x, y, con)

	})
	return rgba
}

func senning(img image.Image, x, y int, xRatio, yRatio float64) color.RGBA64 {
	// 比較のための４点とその位置の比を求める関数
	createParam := func(inSize, pos int, ratio float64) (int, int, float64) {
		v1 := down(pos, ratio)
		v2 := v1 + 1
		if v2 > inSize-1 {
			v2 = inSize - 1
		}
		v3 := float64(pos)/ratio - float64(v1)
		return v1, v2, v3
	}
	//	計測点（パラメータ）を作る
	rect := img.Bounds()
	j1, j2, q := createParam(rect.Size().Y, y, yRatio)
	i1, i2, p := createParam(rect.Size().X, x, xRatio)
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
		r, g, b, a := img.At(v.X, v.Y).RGBA()
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

// LinearDensity は与えられた画像を線形濃度変換を使用して濃度変換する。
func LinearDensity(img image.Image, levelA, levelB uint16) image.Image {
	con := func(val uint32) uint16 {
		v := math.MaxUint16 * (float64((uint16(val) - levelA)) / float64((levelB - levelA)))
		if v > math.MaxUint16 {
			v = math.MaxUint16
		} else if v < 0 {
			v = 0
		}
		return uint16(v)
	}
	rgba := image.NewRGBA(img.Bounds())
	icom.ImageLoop(img, func(x, y int) {
		r, g, b, a := img.At(x, y).RGBA()

		rgba.Set(x, y, color.RGBA64{con(r), con(g), con(b), uint16(a)})
	})
	return rgba
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

// UnlinearDensity は与えられた画像を非線形濃度変換を使用して濃度変換する。
func UnlinearDensity(img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	icom.ImageLoop(img, func(x, y int) {
		r, g, b, a := img.At(x, y).RGBA()
		rgba.Set(x, y, color.RGBA64{unlinerCon(r, 0.5), unlinerCon(g, 0.5), unlinerCon(b, 0.5), uint16(a)})
	})
	return rgba
}

// ContrastImprovement は与えられた画像のコントラスト改善を実施する
func ContrastImprovement(img image.Image) image.Image {
	con := func(val uint32) uint16 {
		maxUint := float64(math.MaxUint16)
		var v uint16
		if val > math.MaxUint16/2 {
			v = uint16((maxUint / 2.0) * (1.0 + math.Sqrt((2.0*float64(val)-maxUint)/maxUint)))
		} else {
			// 非線形濃度変換
			v = unlinerCon(val, 2)
		}
		if v > math.MaxUint16 {
			v = math.MaxUint16
		} else if v < 0 {
			v = 0
		}
		return uint16(v)
	}
	rgba := image.NewRGBA(img.Bounds())
	icom.ImageLoop(img, func(x, y int) {
		r, g, b, a := img.At(x, y).RGBA()
		rgba.Set(x, y, color.RGBA64{con(r), con(g), con(b), uint16(a)})
	})
	return rgba
}

// AverageHistogram は与えられた画像のヒストグラム平均化を実施する。
func AverageHistogram(img image.Image) image.Image {
	rect := img.Bounds()
	//ヒストグラムの取得
	_, _, _, hisL := icom.MakeHistogramData(img)
	//lookupTBLの作成
	luTbl := createLookupTable(hisL, rect)
	//描画
	rgba := image.NewRGBA(img.Bounds())
	icom.ImageLoop(img, func(x, y int) {
		r, g, b, a := img.At(x, y).RGBA()
		rgba.Set(x, y,
			color.RGBA64{luTbl[r], luTbl[g], luTbl[b], uint16(a)})
	})
	Histogram(rgba, "ave histgram", "x", "y", "sampleimage/hist_ave.png")
	return rgba
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

// Histogram は与えられた画像のヒストグラムを作成する
func Histogram(img image.Image, title, xLabel, yLabel, output string) {
	//	グラフの準備
	p, _ := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel
	//プロットの準備
	hisR, hisG, hisB, hisL := icom.MakeHistogramData(img)

	// 各値グラフのプロット
	addPlotter(p, hisR, 0)
	addPlotter(p, hisG, 1)
	addPlotter(p, hisB, 2)
	addPlotter(p, hisL, 3)

	// 出力
	p.Save(6*vg.Inch, 5*vg.Inch, output)
	return
}
func addPlotter(p *plot.Plot, data [ex.ColorWidthAryMax]uint16, key int) {
	var line plotter.XYer
	plots := make(plotter.XYs, len(data))
	line = plots
	for i, v := range data {
		plots[i].X = float64(i)
		plots[i].Y = float64(v)
	}
	graph, _, _ := plotter.NewLinePoints(line)
	graph.Color = plotutil.Color(key)
	p.Add(graph)
	p.Legend.Add("line:"+strconv.Itoa(key), graph)
	return
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
