package effecter

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

func Test_Effect(t *testing.T) {
	// ファイル読み込み
	inputImage := inputFile()
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
	}
	for _, v := range funcs {
		// ファイル出力
		if outputFile(v.Plefex, v.Action()) {

		}
	}
}

func inputFile() image.Image {
	// ファイル読み込み
	inputFile, err := os.Open("sampleimage/test.jpg")
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
	outputFile, err := os.Create("sampleimage/test_" + append + ".jpg")
	if nil != err {
		fmt.Println(err)
		return false
	}
	option := &jpeg.Options{Quality: 100}
	err = jpeg.Encode(outputFile, outputImage, option) // エンコード

	if nil != err {
		fmt.Println(err)
		return false
	}
	defer outputFile.Close()
	return true
}
