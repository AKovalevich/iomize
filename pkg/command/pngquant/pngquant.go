package pngquant

import (
	pngquant "github.com/AKovalevich/go-pngquant"
)

func HandlerPngquant(imageByte []byte, params map[string]string) ([]byte, error) {
	args := []string{}

	args = append(args,"-")
	args = append(args,"--speed")
	args = append(args,"6")
	args = append(args,"--quality")
	args = append(args,"65-80")
	args = append(args, "--output")
	args = append(args, "./test.png")

	compressedImage, err := pngquant.CompressBytes(imageByte, args)
	print("OPTPNG")
	print(compressedImage)
	if err != nil {
		print(err.Error())
		return nil, err
	}
	return compressedImage, nil
}
