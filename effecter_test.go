package effecter

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

func Test_Effect(t *testing.T) {
	// ファイル読み込み
	inputImage := inputFile("sampleimage/test.png")
	if nil == inputImage {
		t.Error("error read file")
	}
	eff := NewEffect(inputImage)

	funcs := []struct {
		Action func() image.Image
		Plefex string
	}{
		{eff.Monochrome, "mono"},
		{eff.ReverseConcentration, "revcon"},
		{eff.FourTone, "fourtone"},
		{eff.ChangeSizeKin, "sizekin"},
		{eff.ChangeSizeSen, "sizesen"},
		{eff.LinearDensity, "linerden"},
		{eff.UnlinearDensity, "unden"},
		{eff.ContrastImprovement, "contrast"},
		{eff.AverageHistogram, "ave"},
	}
	for _, v := range funcs {
		// ファイル出力
		outputFile(v.Plefex, v.Action())
	}
}
func Test_OutMeta(t *testing.T) {
	eff := NewEffect(inputFile("sampleimage/test_ave.png"))
	normalFuncs := []func(){eff.Histogram}
	for _, v := range normalFuncs {
		v()
	}
}

func inputFile(path string) image.Image {
	// ファイル読み込み
	inputFile, err := os.Open(path)
	if nil != err {
		fmt.Println(err)
		return nil
	}
	// decodeの実施
	inputImage, _, err := image.Decode(inputFile)
	if nil != err {
		fmt.Println(err)
		return nil
	}
	inputFile.Close()
	return inputImage
}

func outputFile(append string, outputImage image.Image) bool {
	outputFile, err := os.Create("sampleimage/test_" + append + ".png")
	if nil != err {
		fmt.Println(err)
		return false
	}
	err = png.Encode(outputFile, outputImage) // エンコード

	if nil != err {
		fmt.Println(err)
		return false
	}
	defer outputFile.Close()
	return true
}
