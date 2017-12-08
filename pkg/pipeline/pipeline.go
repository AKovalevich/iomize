package pipeline

import (
	"io/ioutil"
	"sync"

	log "github.com/AKovalevich/scrabbler/log/logrus"
	"github.com/AKovalevich/iomize/pkg/profiler"

	"gopkg.in/yaml.v2"
	"github.com/AKovalevich/iomize/pkg/command/lilliput"
	"github.com/AKovalevich/iomize/pkg/command/pngquant"
	"time"
)

type PipeItem struct {
	sync.RWMutex
	Name       string
	Handler    func([]byte, map[string]string) ([]byte, error)
	Validators []func() error
}

type PipeItemInfo struct {
	sync.RWMutex
	Name   string
	Params map[string]string
}

type PipeLine struct {
	sync.RWMutex
	Name  string
	Pipes []*PipeItemInfo
}

type PipeLineList map[string]*PipeLine

// @TODO Looks like a bad solution. Need to refactor this.
var pipeItemScope = make(map[string]*PipeItem)

func (p *PipeItem) Register() {
	log.Do.Info("Register new pipeline: " + p.Name)
	pipeItemScope[p.Name] = p
}

func (pl *PipeLine) Exec(originImageByte []byte) ([]byte, error) {
	var compressedImageByte []byte
	var err error

	if len(compressedImageByte) <= 0 {
		compressedImageByte = originImageByte
	}

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
			print("\n")
			log.Do.Debug("Started processing pipe: " + pipe.Name)
			profilerTime := time.Now()
			compressedImageByte, err = pipe.Handler(compressedImageByte, e.Params)
			if err != nil {
				return nil, err
			}
			profiler.TimeTrack(profilerTime, pipe.Name)
			log.Do.Debug("Finished processing pipe: " + pipe.Name)
		}

		e.Unlock()
	}

	return compressedImageByte, nil
}

// Initialize PipeLines from configuration file
func InitPipelines(configPath string) (PipeLineList, error) {
	// Register lilliput command
	pipeTest := &PipeItem{
		Name:    "lilliput",
		Handler: lilliput.HandlerLilliput,
	}
	pipeTest.Register()

	pipePngquant := &PipeItem{
		Name:    "pngquant",
		Handler: pngquant.HandlerPngquant,
	}
	pipePngquant.Register()

	optipng := &PipeItem{
		Name:    "optipng",
		Handler: pngquant.HandlerPngquant,
	}
	optipng.Register()

	pipeLineList, err := readPipeLineFromFile(configPath)
	if err != nil {
		return pipeLineList, err
	}

	return pipeLineList, nil
}

// Read PipeLineList from configuration file
func readPipeLineFromFile(configPath string) (PipeLineList, error) {
	var pipeLineList PipeLineList

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return pipeLineList, err
	}
	yaml.Unmarshal(file, &pipeLineList)

	return pipeLineList, err
}
