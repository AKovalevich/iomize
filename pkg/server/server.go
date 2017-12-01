package server

import (
	"os/signal"
	"net/http"
	"context"
	"time"
	"sync"
	"os"

	"github.com/AKovalevich/iomize/pkg/config"
)

// Server is the reverse-proxy/load-balancer engine
type Server struct {
	mainConfiguration *config.MainConfiguration
	signals							chan os.Signal
	stopChan						chan bool
	mainHttpServer					*http.Server
}

func NewServer(config *config.MainConfiguration) Server {
	var server Server
	server.mainConfiguration = config

	// Configure signals
	server.configureSignals()
	go server.listenSignals()

	// Configure main scrabbler HTTP server
	server.configureMainHttpServer()

	return server
}

func (server *Server) Serve() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		go server.runMainServer()
	}()
	wg.Wait()
}

func (server *Server) Stop() {
	var entryPoints []*http.Server
	entryPoints = append(entryPoints, server.mainHttpServer)
	var wg sync.WaitGroup
	for _, v := range entryPoints {
		wg.Add(1)
		go func(srv *http.Server) {
			defer wg.Done()
			graceTimeOut := time.Duration(server.mainConfiguration.GraceTimeOut)
			ctx, cancel := context.WithTimeout(context.Background(), graceTimeOut)
			if err := srv.Shutdown(ctx); err != nil {
				srv.Close()
			}
			cancel()
		}(v)
	}

	wg.Wait()
	server.Close()
}

func (server * Server) Close() {
	// Close Web UI HTTP server
	signal.Stop(server.signals)
	close(server.signals)
	os.Exit(1)
}
