package sdk_camera

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
)

func BGRToJpeg(o io.Writer, bgr []byte, w int, h int, opt *jpeg.Options) error {
	if len(bgr) != w*h*3 {
		return fmt.Errorf("bgr input error")
	}

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	x, y := 0, 0
	for i := 0; i < len(bgr); i++ {
		if i > 0 && i%3 == 0 {
			b, g, r := bgr[i-3], bgr[i-2], bgr[i-1]
			rgba.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 0})
			if x > 0 && x%(w-1) == 0 {
				x = 0
				y++
			} else {
				x++
			}
		}
	}

	return jpeg.Encode(o, rgba, opt)
}
