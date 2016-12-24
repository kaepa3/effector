package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"localhost/imaging/local"
	"os"
)

func main() {
	// ファイル読み込み
	inputFile, err := os.Open("plant4.jpg")

	if nil != err {
		fmt.Println(err)
		return
	}
	inputImage, _, err := image.Decode(inputFile)

	if nil != err {
		fmt.Println(err)
	}

	defer inputFile.Close()

	// ファイル出力
	outputFile, err := os.Create("plant4_out.jpg")
	if nil != err {
		fmt.Println(err)
	}

	outputImage := convert(inputImage) // 変換
	option := &jpeg.Options{Quality: 100}
	err = jpeg.Encode(outputFile, outputImage, option) // エンコード

	if nil != err {
		fmt.Println(err)
	}

	defer outputFile.Close()
}

func convert(inputImage image.Image) image.Image {
	eff := effect.NewEffect(inputImage)

	return eff.ConvertToMonochromeImage()
}
