package effector

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

func Test_Effect(t *testing.T) {
	// ファイル読み込み
	in := inputFile("sampleimage/test.jpg")
	if nil == in {
		t.Error("error read file")
	}

	outputFile("mono", Monochrome(in))
	outputFile("revcon", ReverseConcentration(in))
	outputFile("fourtone", FourTone(in))
	outputFile("linerden", LinearDensity(in, 0x10, 0xFF00))
	outputFile("unden", UnlinearDensity(in))
	outputFile("contrast", ContrastImprovement(in))
	outputFile("ave", AverageHistogram(in))

}

func Test_SizeChange(t *testing.T) {
	img := inputFile("sampleimage/test.jpg")
	outputFile("sizekin", ChangeSizeKin(img, 0.8, 0.8))
	outputFile("sizesen", ChangeSizeSen(img, 0.7, 0.2))
}

func Test_OutMeta(t *testing.T) {
	img := inputFile("sampleimage/test.jpg")
	Histogram(img, "test.png hist", "nodo", "dosu", "sampleimage/hist_org.png")
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
