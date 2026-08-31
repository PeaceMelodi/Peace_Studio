package worker

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func buildOutputPath(jobID string, action string, srcPath string) string {
	ext := filepath.Ext(srcPath)
	filename := jobID + "-" + action + ext
	return filepath.Join("uploads", "processed", filename)
}

func ToGrayscale(jobID string, srcPath string) (string, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	img, format, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayImg.Set(x, y, color.GrayModel.Convert(originalColor))
		}
	}

	dstPath := buildOutputPath(jobID, "grayscale", srcPath)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(dstFile, grayImg)
	case "gif":
		err = gif.Encode(dstFile, grayImg, nil)
	default:
		err = jpeg.Encode(dstFile, grayImg, nil)
	}
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func ToBlur(jobID string, srcPath string) (string, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	img, format, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	blurredImg := image.NewRGBA(bounds)

	const radius = 8

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum, count uint32

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					nx, ny := x+dx, y+dy
					if nx < bounds.Min.X || nx >= bounds.Max.X || ny < bounds.Min.Y || ny >= bounds.Max.Y {
						continue
					}
					r, g, b, a := img.At(nx, ny).RGBA()
					rSum += r
					gSum += g
					bSum += b
					aSum += a
					count++
				}
			}

			blurredImg.Set(x, y, color.RGBA{
				R: uint8((rSum / count) >> 8),
				G: uint8((gSum / count) >> 8),
				B: uint8((bSum / count) >> 8),
				A: uint8((aSum / count) >> 8),
			})
		}
	}

	dstPath := buildOutputPath(jobID, "blur", srcPath)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(dstFile, blurredImg)
	case "gif":
		err = gif.Encode(dstFile, blurredImg, nil)
	default:
		err = jpeg.Encode(dstFile, blurredImg, nil)
	}
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func ToResize(jobID string, srcPath string) (string, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	img, format, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	const targetWidth = 300

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	dstPath := buildOutputPath(jobID, "resize", srcPath)

	if origWidth <= targetWidth {
		srcFile.Seek(0, 0)
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return "", err
		}
		defer dstFile.Close()
		switch strings.ToLower(format) {
		case "png":
			err = png.Encode(dstFile, img)
		case "gif":
			err = gif.Encode(dstFile, img, nil)
		default:
			err = jpeg.Encode(dstFile, img, nil)
		}
		if err != nil {
			return "", err
		}
		return dstPath, nil
	}

	targetHeight := origHeight * targetWidth / origWidth

	resizedImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := bounds.Min.X + x*origWidth/targetWidth
			srcY := bounds.Min.Y + y*origHeight/targetHeight
			resizedImg.Set(x, y, img.At(srcX, srcY))
		}
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(dstFile, resizedImg)
	case "gif":
		err = gif.Encode(dstFile, resizedImg, nil)
	default:
		err = jpeg.Encode(dstFile, resizedImg, nil)
	}
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func ToPixelate(jobID string, srcPath string) (string, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	img, format, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	pixelatedImg := image.NewRGBA(bounds)

	const blockSize = 15

	for by := bounds.Min.Y; by < bounds.Max.Y; by += blockSize {
		for bx := bounds.Min.X; bx < bounds.Max.X; bx += blockSize {
			var rSum, gSum, bSum, aSum, count uint32

			maxY := by + blockSize
			if maxY > bounds.Max.Y {
				maxY = bounds.Max.Y
			}
			maxX := bx + blockSize
			if maxX > bounds.Max.X {
				maxX = bounds.Max.X
			}

			for y := by; y < maxY; y++ {
				for x := bx; x < maxX; x++ {
					r, g, b, a := img.At(x, y).RGBA()
					rSum += r
					gSum += g
					bSum += b
					aSum += a
					count++
				}
			}

			avgColor := color.RGBA{
				R: uint8((rSum / count) >> 8),
				G: uint8((gSum / count) >> 8),
				B: uint8((bSum / count) >> 8),
				A: uint8((aSum / count) >> 8),
			}

			for y := by; y < maxY; y++ {
				for x := bx; x < maxX; x++ {
					pixelatedImg.Set(x, y, avgColor)
				}
			}
		}
	}

	dstPath := buildOutputPath(jobID, "pixelate", srcPath)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(dstFile, pixelatedImg)
	case "gif":
		err = gif.Encode(dstFile, pixelatedImg, nil)
	default:
		err = jpeg.Encode(dstFile, pixelatedImg, nil)
	}
	if err != nil {
		return "", err
	}

	return dstPath, nil
}