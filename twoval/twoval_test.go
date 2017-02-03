package twoval

import (
	"image"
	"image/color"
	"log"
	"testing"

	"github.com/kaepa3/effector/ex"
	"github.com/kaepa3/effector/icom"
	"github.com/kaepa3/effector/testutill"
)

func Test_StaticThreshold(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	icom.ImageLoop(img, func(x int, y int) {
		val := uint16(x)
		rgba.Set(x, y, color.RGBA64{val, val, val, 1})
	})
	doFunc := StaticThreshold(2)
	col := doFunc(rgba, 1, 1)
	if col.R != 0 || col.G != 0 || col.B != 0 {
		t.Error("val err")
	}
	col = doFunc(rgba, 2, 1)
	if col.R != ex.ColorWidth || col.G != ex.ColorWidth || col.B != ex.ColorWidth {
		t.Error("val err")
	}
}

func Test_getVariableRange(t *testing.T) {
	a, b := getVariableRange(1, 5, 2)
	if a != 0 || b != 2 {
		t.Errorf("range err%d:%d", a, b)
	}
}

func Test_VariableThreshold(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	rgba := image.NewRGBA64(img.Bounds())
	//	make test image
	icom.ImageLoop(img, func(x int, y int) {
		rgba.Set(x, y, color.RGBA64{0, 0, 0, 1})
	})

	doFunc := VariableThreshold(10, 5, false)
	col := doFunc(rgba, 1, 1)
	if col.R != 0 || col.G != 0 || col.B != 0 {
		t.Error("val err")
	}
}
func createImage(x, y int) *image.RGBA {
	rect := image.Rectangle{image.Point{0, 0}, image.Point{x, y}}
	return image.NewRGBA(rect)
}
func createTestImage(testData [][]uint) *image.RGBA {
	allLen := len(testData)
	oneLen := len(testData[0])
	srcImg := createImage(oneLen, allLen)
	pt := srcImg.Bounds().Size()
	for x := 0; x < pt.X; x++ {
		for y := 0; y < pt.Y; y++ {
			if testData[y][x] != 0 {
				srcImg.Set(x, y, color.RGBA{1, 1, 1, 1})
			}
		}
	}
	return srcImg
}

func Test_getArounds(t *testing.T) {
	val := [][]uint{
		{0, 0, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	srcImg := createTestImage(val)
	around, sum := getArounds(srcImg, 1, 1)
	expected := []uint32{1, 1, 0, 0, 1, 1, 1, 1, 0}
	for i, v := range expected {
		if around[i] != 0 && v == 0 {
			t.Errorf("error->%d[%d:%d]", i, v, around[i])
		}
		if around[i] == 0 && v != 0 {
			t.Errorf("error->%d[%d:%d]", i, v, around[i])
		}
	}
	if sum != 2 {
		t.Errorf("connect error->%d", sum)
	}
}
func Test_getArounds2(t *testing.T) {
	val := [][]uint{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	srcImg := createTestImage(val)
	_, sum := getArounds(srcImg, 1, 1)
	if sum != 0 {
		t.Errorf("connect error->%d", sum)
	}
}
func Test_isCheckValue(t *testing.T) {
	val := [][]uint{
		{0, 0, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	srcImg := createTestImage(val)
	if true == isCheckValue(0, 1, srcImg) {
		t.Errorf("x.error")
	} else if true == isCheckValue(1, 0, srcImg) {
		t.Errorf("y.error")
	} else if false == isCheckValue(1, 1, srcImg) {
		r, _, _, _ := srcImg.At(1, 1).RGBA()
		t.Errorf("true error%d", r)
	}
	srcImg.Set(1, 1, color.RGBA64{ex.ColorWidth, 10, 10, 10})
	if true == isCheckValue(1, 1, srcImg) {
		t.Error("false error")
	}
}

func Test_countConnect(t *testing.T) {
	val := [][]uint{
		{0, 0, 0},
		{1, 0, 1},
		{1, 1, 1},
	}
	srcImg := createTestImage(val)
	in, sum := getArounds(srcImg, 1, 1)
	if sum != 3 {
		t.Errorf("sum Error%d", sum)
	}
	cnt := countConnect(in)
	if cnt != 1 {
		t.Errorf("count Error:%d", cnt)
	}
}

func Test_Thinning(t *testing.T) {
	cases := []struct {
		in  [][]uint
		out [][]uint32
	}{{
		[][]uint{
			{1, 1, 1},
			{1, 0, 1},
			{1, 1, 1},
		},
		[][]uint32{
			{1, 1, 1},
			{1, 1, 1},
			{1, 1, 1},
		},
	}, {
		[][]uint{
			{0, 1, 1},
			{0, 0, 1},
			{1, 1, 1},
		},
		[][]uint32{
			{0, 1, 1},
			{0, 1, 1},
			{1, 1, 1},
		},
	}}

	for i, v := range cases {
		img := createTestImage(v.in)
		dst := Thinning(img)
		pt := dst.Bounds().Size()
		for x := 0; x < pt.X; x++ {
			for y := 0; y < pt.Y; y++ {
				r, _, _, _ := dst.At(x, y).RGBA()
				if 0 == v.out[y][x] && r != 0 {
					log.Printf("err test case:%d", i)
					t.Errorf("error [%d,%d]expect %d, actual %d", x, y, r, v.out[y][x])
				} else if 0 != v.out[y][x] && r == 0 {
					log.Printf("err test case:%d", i)
					t.Errorf("error [%d,%d]expect %d, actual %d", x, y, r, v.out[y][x])
				}
			}
		}
	}
}

func Test_BoundaryTracking(t *testing.T) {
	cases := []struct {
		in  [][]uint
		out [][]uint32
	}{{
		[][]uint{
			{1, 1, 1, 0, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
		},
		[][]uint32{
			{1, 1, 1, 0, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 0, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1},
		},
	}}

	for _, v := range cases {
		img := createTestImage(v.in)
		dst := BoundaryTracking(img)
		icom.ImageLoop(img, func(x int, y int) {
			r, _, _, _ := dst.At(x, y).RGBA()
			if 0 == v.out[y][x] && r != 0 {
				t.Errorf("error [%d,%d]expect %d, actual %d", x, y, r, v.out[y][x])
			} else if 0 != v.out[y][x] && r == 0 {
				t.Errorf("error [%d,%d]expect %d, actual %d", x, y, r, v.out[y][x])
			}
		})
	}
}
