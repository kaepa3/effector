package histogram

import (
	"bytes"
	"image"
	"strconv"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
)

// Histogram は与えられた画像のヒストグラムを作成する
func Output(img image.Image, title, xLabel, yLabel string) image.Image {
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
	w, err := p.WriterTo(6*vg.Inch, 5*vg.Inch, "png")
	if err != nil {
		return img
	}
	var rb = bytes.NewBuffer([]byte{})
	w.WriteTo(rb)
	hisImg, _, err := image.Decode(rb)
	if err != nil {
		return img
	}
	return hisImg
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
