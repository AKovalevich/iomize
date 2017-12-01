package backend

import (
	"net/http"

	"github.com/AKovalevich/iomize/pkg/route"
	"github.com/AKovalevich/iomize/pkg/pipeline"
	"os"
	"fmt"
	"io/ioutil"
	"log"
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
				inputBuf, err := ioutil.ReadFile("./example.png")
				compressedImage, err := pipeLineList["pngquant"].Exec(inputBuf)
				if err != nil {
					log.Panic(err)
				}
				print(compressedImage)
				err = ioutil.WriteFile("test2.png", compressedImage, 0775)
				if err != nil {
					fmt.Printf("error writing out resized image, %s\n", err)
					os.Exit(1)
				}
				fmt.Fprintf(w, "Test")
			},
		},
	}
}

func (txe *Entrypoint) String() string {
	return txe.Name
}
