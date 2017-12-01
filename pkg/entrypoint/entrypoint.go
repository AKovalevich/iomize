package entrypoint

import (
	"strings"
	"fmt"

	"github.com/AKovalevich/iomize/pkg/route"
	"github.com/AKovalevich/iomize/pkg/pipeline"
	"github.com/AKovalevich/iomize/pkg/entrypoint/backend"
)

// Base Entrypoint interface
type Entrypoint interface {
	// Initialize entrypoint
	Init(pipeLineList pipeline.PipeLineList)
	// Get entrypoint name
	String() string
	// Get entrypoint routes list
	RoutesList() []route.Route
}

type EntrypointList []Entrypoint

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (e *EntrypointList) Set(value string) error {
	entrypoints := strings.Split(value, ",")

	if len(entrypoints) == 0 {
		return fmt.Errorf("bad EntryPointList format: %s", value)
	}

	for _, entrypointName := range entrypoints {
		// Try to create entrypoint
		switch entrypointName {
		case "backend":
			backendEntrypoint := backend.New()
			backendEntrypoint.Name = "backend"
			*e = append(*e, backendEntrypoint)
			break
		case "web":
		case "admin":
			break
		}
	}

	return nil
}

// Get return the EntryPoints map
func (e *EntrypointList) Get() interface{} {
	return EntrypointList(*e)
}

// SetValue sets the EntryPoints map with val
func (e *EntrypointList) SetValue(val interface{}) {
	*e = EntrypointList(val.(EntrypointList))
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (e *EntrypointList) String() string {
	var entrypoints []string
	for _, entrypoint := range *e {
		// Try to create entrypoint
		entrypoints = append(entrypoints, entrypoint.String())
	}

	return strings.Join(entrypoints, ", ")
}

// Type is type of the struct
func (e *EntrypointList) Type() string {
	return "defaultentrypoints"
}
