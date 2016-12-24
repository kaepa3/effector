package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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

	outputImage := convertToMonochromeImage(inputImage) // 変換
	option := &jpeg.Options{Quality: 100}
	err = jpeg.Encode(outputFile, outputImage, option) // エンコード

	if nil != err {
		fmt.Println(err)
	}

	defer outputFile.Close()
}

func convertToMonochromeImage(inputImage image.Image) image.Image {
	rect := inputImage.Bounds()
	width := rect.Size().X
	height := rect.Size().Y
	rgba := image.NewRGBA(rect)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var col color.RGBA64
			// 座標(x,y)のR, G, B, α の値を取得
			r, g, b, a := inputImage.At(x, y).RGBA()

			// それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
			outR := float32(r) * 0.298912
			outG := float32(g) * 0.58611
			outB := float32(b) * 0.114478
			mono := uint16(outR + outG + outB)
			col.R = mono
			col.G = mono
			col.B = mono
			col.A = uint16(a)
			rgba.Set(x, y, col)
		}
	}

	return rgba.SubImage(rect)
}
