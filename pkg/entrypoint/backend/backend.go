package backend

import (
	"net/http"

	"gitlab.com/artemkovalevich00/iomize/pkg/route"
	"gitlab.com/artemkovalevich00/iomize/pkg/pipeline"
	"os"
	"image"
	"image/png"
)

const (
	DefaultApiPrefix = "api/v1"
)

type Entrypoint struct {
	Name string
	Routes []route.Route
}

func New() *Entrypoint {
	entrypoint := &Entrypoint{}
	return entrypoint
}

func (txe *Entrypoint) RoutesList() []route.Route {
	return txe.Routes
}

// Initialize entrypoint
func (txe *Entrypoint) Init(pipeLineList pipeline.PipeLineList) {
	txe.Routes = []route.Route{
		{
			Path: "/" + DefaultApiPrefix +  "/hello",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fImg1, _ := os.Open("./example.png")
				defer fImg1.Close()
				img, _, _ := image.Decode(fImg1)
				pipeLineList["pngquant"].Exec(&img)
				toimg, _ := os.Create("./example_optimized.png")
				defer toimg.Close()
				png.Encode(toimg, img)
			},
		},
	}
}

func (txe *Entrypoint) String() string {
	return txe.Name
}
