package pipeline

import (
	"sync"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"gitlab.com/artemkovalevich00/iomize/pkg/command/test"
	log "github.com/AKovalevich/scrabbler/log/logrus"

	"image"
	"gitlab.com/artemkovalevich00/iomize/pkg/command/pngquant"
)

type PipeItem struct {
	sync.RWMutex
	Name string
	Handler func(*image.Image, map[string]string) error
	Validators []func() error
}

type PipeItemInfo struct {
	sync.RWMutex
	Name string
	Params map[string]string
}

type PipeLine struct {
	sync.RWMutex
	Name string
	Pipes []*PipeItemInfo
}

type PipeLineList map[string]*PipeLine

// @TODO Looks like a bad solution. Need to refactor this.
var pipeItemScope = make(map[string]*PipeItem)

func (p *PipeItem) Register() {
	log.Do.Info("Register new pipeline: " + p.Name)
	pipeItemScope[p.Name] = p
}

func (pl *PipeLine) Exec(image *image.Image) error {
	for _, e := range pl.Pipes {
		e.Lock()
		// Run validation handlers
		//for _, v := range e.Validate {
		//	err = v()
		//	if err != nil {
		//		break
		//	}
		//}

		if pipe, ok := pipeItemScope[e.Name]; ok {
			err := pipe.Handler(image, e.Params)
			if err != nil {
				return err
			}
		}

		e.Unlock()
	}

	return nil
}

// Initialize PipeLines from configuration file
func InitPipelines(configPath string) (PipeLineList, error) {
	// Register test command
	pipeTest := &PipeItem{
		Name: "test",
		Handler: test.HandlerTest,
	}
	pipeTest.Register()

	pipePngquant := &PipeItem{
		Name: "pngquant",
		Handler: pngquant.HandlerPngquant,
	}
	pipePngquant.Register()

	pipeLineList, err := readPipeLineFromFile(configPath)
	if err != nil {
		return pipeLineList, err
	}

	return pipeLineList, nil
}

// Read PipeLineList from configuration file
func readPipeLineFromFile(configPath string) (PipeLineList, error)  {
	var pipeLineList PipeLineList

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return pipeLineList, err
	}
	yaml.Unmarshal(file, &pipeLineList)

	return pipeLineList, err
}
