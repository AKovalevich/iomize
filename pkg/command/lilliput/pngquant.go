package lilliput

import (
	"fmt"
	"os"
	"github.com/discordapp/lilliput"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": map[int]int{lilliput.JpegQuality: 85},
	".png":  map[int]int{lilliput.PngCompression: 7},
	".webp": map[int]int{lilliput.WebpQuality: 85},
}

func HandlerLilliput(imageByte []byte, params map[string]string) ([]byte, error) {
	var outputWidth int
	var outputHeight int
	var stretch bool

	decoder, err := lilliput.NewDecoder(imageByte)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		fmt.Printf("error reading image header, %s\n", err)
		os.Exit(1)
	}

	// print some basic info about the image
	fmt.Printf("image type: %s\n", decoder.Description())
	fmt.Printf("%dpx x %dpx\n", header.Width(), header.Height())

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(1024)
	defer ops.Close()

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)


	outputWidth = 50
	outputHeight = 50


	resizeMethod := lilliput.ImageOpsFit
	if stretch {
		resizeMethod = lilliput.ImageOpsResize
	}

	outputType := ".png"
	opts := &lilliput.ImageOptions{
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		FileType:             outputType,
		EncodeOptions:        EncodeOptions[outputType],
	}

	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return nil, err
	}

	return outputImg, nil
}
