package optipng

import (
	optipng "github.com/AKovalevich/go-optipng"
)


func HandlerPngquant(imageByte []byte, params map[string]string) ([]byte, error) {
	args := []string{}

	args = append(args, "-")
	args = append(args, "--speed")
	args = append(args, "6")
	args = append(args, "--quality")
	args = append(args, "65-80")
	//args = append(args, "--output")
	//args = append(args, "./test.png")

	compressedImage, err := optipng.CompressBytes(imageByte, args)
	if err != nil {
		return nil, err
	}
	return compressedImage, nil
}
