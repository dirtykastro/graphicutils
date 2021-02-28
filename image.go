package graphicutils

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
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

	mime, mimeErr := getContentTypeFromFile(file)
	if mimeErr != nil {
		return
	}

	if mime == "image/jpeg" {
		img, err = jpeg.Decode(file)
	} else if mime == "image/gif" {
		img, err = gif.Decode(file)
	} else if mime == "image/png" {
		img, err = png.Decode(file)
	}

	return
}

func getContentTypeFromFile(file *os.File) (mime string, err error) {
	fileBuffer := make([]byte, 512)
	_, readErr := file.Read(fileBuffer)

	if readErr != nil {
		return
	}

	mime = http.DetectContentType(fileBuffer)

	return
}

// Get MIME type of file
func GetContentType(path string) (mime string, err error) {
	file, fileErr := os.Open(path)
	if fileErr != nil {
		return
	}

	mime, err = getContentTypeFromFile(file)

	file.Close()

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
	if fgPixel.A == 255 {

		pixel = fgPixel
		return

	} else if fgPixel.A == 0 {

		pixel = bgPixel
		return
	}

	fgAlpha := float64(fgPixel.A)
	bgAlpha := float64(bgPixel.A)

	alpha := fgAlpha + ((255 - fgAlpha) * (bgAlpha / 255))

	pixel.R = blendColor(fgPixel.R, fgPixel.A, bgPixel.R, bgPixel.A)
	pixel.G = blendColor(fgPixel.G, fgPixel.A, bgPixel.G, bgPixel.A)
	pixel.B = blendColor(fgPixel.B, fgPixel.A, bgPixel.B, bgPixel.A)
	pixel.A = uint8(math.Round(alpha))

	return
}

func blendColor(fgColor, fgAlpha, bgColor, bgAlpha uint8) (finalColor uint8) {

	alpha := float64(fgAlpha) / 255

	fgColorFloat := float64(fgColor)
	bgColorFloat := float64(bgColor)

	out := alpha*fgColorFloat + (1-alpha)*bgColorFloat

	finalColor = uint8(math.Round(out))

	return
}

func GetImageDimensions(imagePath string) (width int, height int, err error) {
	file, file_err := os.Open(imagePath)
	if file_err != nil {
		err = file_err
		return
	}

	image, _, image_err := image.DecodeConfig(file)

	file.Close()

	if image_err != nil {
		err = image_err
		return
	}

	width = image.Width
	height = image.Height

	return
}
