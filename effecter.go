package effecter

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

//Effect is image effect interface
type Effect interface {
	Monochrome() image.Image
	ReverseConcentration() image.Image
	FourTone() image.Image
	ChangeSizeKin() image.Image
	ChangeSizeSen() image.Image
	LinearDensity() image.Image
	UnlinearDensity() image.Image
	ContrastImprovement() image.Image
	AverageHistogram() image.Image
	Histogram(title, xLabel, yLabel, output string)
}
type effect struct {
	inputImage image.Image
}

//NewEffect is create new effect object
func NewEffect(in image.Image) Effect {
	return &effect{in}
}

func (ef *effect) imageLoop(rect image.Rectangle, effectFunction func(x, y int) color.RGBA64) image.Image {
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// 座標(x,y)のR, G, B, α の値を取得
			col := effectFunction(x, y)
			rgba.Set(x, y, col)
		}
	}
	return rgba
}

// 逆にする。
func (ef *effect) ReverseConcentration() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		col.R = math.MaxUint16 - uint16(r)
		col.G = math.MaxUint16 - uint16(g)
		col.B = math.MaxUint16 - uint16(b)
		col.A = uint16(a)
		return col
	})
}

// モノクロにする
func (ef *effect) Monochrome() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		//それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
		outR := float32(r) * 0.298912
		outG := float32(g) * 0.58611
		outB := float32(b) * 0.114478
		mono := uint16(outR + outG + outB)
		col.R = mono
		col.G = mono
		col.B = mono
		col.A = uint16(a)
		return col
	})
}

// 4階調にする
func (ef *effect) FourTone() image.Image {
	var col color.RGBA64
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
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
		return col
	})
}

var up = func(val int, rat float64) int { return int(float64(val) * rat) }
var down = func(val int, rat float64) int { return int(float64(val) / rat) }

// 最近傍方
func (ef *effect) ChangeSizeKin() image.Image {

	// とりあえず正方形にして、1/2倍にする
	ratio := 0.9
	rect := ef.inputImage.Bounds()
	width := up(rect.Size().X, ratio)
	height := width
	newRect := image.Rect(0, 0, width, height)
	yRatio := float64(height) / float64(rect.Size().Y)

	return ef.imageLoop(newRect, func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(down(x, ratio), down(y, yRatio)).RGBA()
		return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	})
}

func (ef *effect) ChangeSizeSen() image.Image {

	// とりあえず正方形にして、1/2倍にする
	ratio := 0.9
	rect := ef.inputImage.Bounds()
	width := up(rect.Size().X, ratio)
	height := width
	newRect := image.Rect(0, 0, width, height)
	yRatio := float64(height) / float64(rect.Size().Y)

	return ef.imageLoop(newRect, func(x, y int) color.RGBA64 {
		return ef.senning(x, y, ratio, yRatio)
	})
}

func (ef *effect) senning(x, y int, xRatio, yRatio float64) color.RGBA64 {
	// 比較のための４点とその位置の比を求める関数
	createParam := func(inSize, outSize int, ratio float64) (int, int, float64) {
		v1 := down(outSize, ratio)
		v2 := v1 + 1
		if v2 > inSize-1 {
			v2 = inSize - 1
		}
		v3 := float64(outSize)/ratio - float64(v1)
		return v1, v2, v3
	}
	//	計測点（パラメータ）を作る
	point := ef.inputImage.Bounds().Size()
	j1, j2, q := createParam(point.Y, y, yRatio)
	i1, i2, p := createParam(point.X, x, xRatio)
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
		r, g, b, a := ef.inputImage.At(v.X, v.Y).RGBA()
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

const RGBAMax = 3
const IndexR = 0
const IndexG = 1
const IndexB = 2
const ColorWidth = math.MaxUint16 + 1

// 最近傍方
func (ef *effect) Histogram(title, xLabel, yLabel, output string) {
	//	グラフの準備
	p, _ := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel
	//プロットの準備
	hisR, hisG, hisB, hisL := ef.makeHistogramData()

	// 各値グラフのプロット
	addPlotter(p, hisR, 0)
	addPlotter(p, hisG, 1)
	addPlotter(p, hisB, 2)
	addPlotter(p, hisL, 3)

	// 出力
	p.Save(6*vg.Inch, 5*vg.Inch, output)
	return
}

func addPlotter(p *plot.Plot, data [ColorWidth]uint16, key int) {
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

func (ef *effect) makeHistogramData() (rD, gD, bD, lD [ColorWidth]uint16) {
	ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		//画素をそれぞれ数える
		r, g, b, _ := ef.inputImage.At(x, y).RGBA()
		rD[r]++
		gD[g]++
		bD[b]++
		l := (3*r + 6*g + b) / 10
		lD[l]++
		return color.RGBA64{0, 0, 0, 0}
	})
	return
}

// 線形濃度変換
func (ef *effect) LinearDensity() image.Image {
	con := func(val uint32) uint16 {
		var levelA uint16 = 0x10
		var levelB uint16 = 0xFF00

		v := math.MaxUint16 * (float64((uint16(val) - levelA)) / float64((levelB - levelA)))
		if v > math.MaxUint16 {
			v = math.MaxUint16
		} else if v < 0 {
			v = 0
		}
		return uint16(v)
	}
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		return color.RGBA64{con(r), con(g), con(b), uint16(a)}
	})
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

// 非線形濃度変換
func (ef *effect) UnlinearDensity() image.Image {

	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		return color.RGBA64{unlinerCon(r, 0.5), unlinerCon(g, 0.5), unlinerCon(b, 0.5), uint16(a)}
	})
}

// コントラスト改善
func (ef *effect) ContrastImprovement() image.Image {
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
	return ef.imageLoop(ef.inputImage.Bounds(), func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		return color.RGBA64{con(r), con(g), con(b), uint16(a)}
	})
}

// ヒストグラム平均化
func (ef *effect) AverageHistogram() image.Image {
	rect := ef.inputImage.Bounds()
	//ヒストグラムの取得
	_, _, _, hisL := ef.makeHistogramData()
	//lookupTBLの作成
	luTbl := createLookupTable(hisL, rect)
	//描画
	var buf effect
	buf.inputImage = ef.imageLoop(rect, func(x, y int) color.RGBA64 {
		r, g, b, a := ef.inputImage.At(x, y).RGBA()
		return color.RGBA64{luTbl[r], luTbl[g], luTbl[b], uint16(a)}
	})
	buf.Histogram("ave histgram", "x", "y", "sampleimage/hist_ave.png")
	return buf.inputImage
}
func createLookupTable(his [ColorWidth]uint16, rect image.Rectangle) (table [ColorWidth]uint16) {
	var sum, val uint16
	//平均画素数
	average := uint16(((rect.Size().X * rect.Size().Y) / ColorWidth) + 1)
	//平均値の計算
	for i, v := range his {
		sum += v
		val += sum / average
		sum %= average
		table[i] = val
	}
	return
}
