package pngquant

import (
	pngquant "github.com/AKovalevich/go-pngquant"
	"image"
)

func HandlerPngquant(image *image.Image, params map[string]string) error {
	args := []string{}

	args = append(args,"-")
	args = append(args,"--speed")
	args = append(args,"6")
	args = append(args,"--quality")
	args = append(args,"65-80")
	args = append(args, "--output")
	args = append(args, "./test.png")

	compressedImage, err := pngquant.Compress(*image, args)
	if err != nil {
		return err
	}
	image = &compressedImage
	return nil
}
