package graphicutils

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func DecodeImage(filePath string) (img image.Image, err error) {

	if !(Exists(filePath)) {
		err = errors.New(filePath + " doesn't exist")
		return
	}

	if IsDirectory(filePath) {
		err = errors.New(filePath + " is directory")
		return
	}

	file, fileErr := os.Open(filePath)
	defer file.Close()

	if fileErr != nil {
		err = fileErr
		return
	}

	img, err = png.Decode(file)

	return
}

func GetPixelValue(img image.Image, pt image.Point) (pixel Pixel, err error) {

	size := img.Bounds()

	if pt.X < 0 || pt.X >= size.Max.X || pt.Y < 0 || pt.Y >= size.Max.Y {
		err = errors.New("invalid Point")
		return
	}

	r, g, b, a := img.At(pt.X, pt.Y).RGBA()

	pixel.R = uint8(r)
	pixel.G = uint8(g)
	pixel.B = uint8(b)
	pixel.A = uint8(a)

	return
}

func BlendPixel(fgPixel Pixel, bgPixel Pixel) (pixel Pixel) {

	fgAlpha := float64(fgPixel.A)
	bgAlpha := float64(bgPixel.A)

	r := fgPixel.R
	g := fgPixel.G
	b := fgPixel.B

	alpha := fgAlpha + ((255 - fgAlpha) * (bgAlpha / 255))

	fmt.Println(alpha)

	pixel.R = uint8(r)
	pixel.G = uint8(g)
	pixel.B = uint8(b)
	pixel.A = uint8(math.Round(alpha))

	return
}

func addColor(fgColor, fgAlpha, bgColor, bgAlpha uint8) (finalColor uint8) {
	return fgColor

	difference := int(fgColor) - int(bgColor)

	finalColor = uint8(int(fgColor) + (difference * (255 - int(fgAlpha))))

	return
}
