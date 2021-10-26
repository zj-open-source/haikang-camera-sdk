package sdk_camera

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
)

import (
	"bytes"
	"golang.org/x/image/bmp"
)

// BGRToJpeg BGR转jpeg
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

// GrayToJpeg 灰度图转jpeg
func GrayToJpeg(o io.Writer, gray []byte, w int, h int, opt *jpeg.Options) error {
	if len(gray) != w*h {
		return fmt.Errorf("gray input error")
	}

	grayImage := image.NewGray(image.Rect(0, 0, w, h))

	for i, x, y := 0, 0, 0; i < len(gray); i++ {
		grayImage.Set(x, y, color.Gray{Y: gray[i]})
		switch {
		case x == 0:
			x++
		case x%(w-1) != 0:
			x++
		case x%(w-1) == 0:
			x = 0
			y++
		}
	}

	return bmp.Encode(o, grayImage)
}

// ImageToJpeg image对象转jpeg
func ImageToJpeg(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
