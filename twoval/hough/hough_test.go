package hough

import (
	"testing"

	"github.com/kaepa3/effector/testutill"
)

func Test_Hough(t *testing.T) {
	img := testutill.CreateImg(5, 5)
	Hough(img)

}
