package server

import (
	"net/http"
	"time"

	log "gitlab.com/artemkovalevich00/iomize/pkg/log/logrus"
)

func (server *Server) configureMainHttpServer() {
	s := http.NewServeMux()
	for _, entrypoint := range server.mainConfiguration.EntryPoints {
		log.Do.Info(entrypoint.String())

		log.Do.Info("Initialize " + entrypoint.String() + " entrypoint...")
		entrypoint.Init(server.mainConfiguration.PipeLineList)
		log.Do.Info("Prepare " + entrypoint.String() + " routes...")
		routesList := entrypoint.RoutesList()
		for _, route := range routesList {
			s.HandleFunc(route.Path, route.Handler)
		}
	}
	server.mainHttpServer = &http.Server{
		Handler: s,
		Addr: server.mainConfiguration.Host + ":" + server.mainConfiguration.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (server *Server) runMainServer() {
	log.Do.Info("Server started on " + server.mainConfiguration.Host + ":" + server.mainConfiguration.Port)
	server.mainHttpServer.ListenAndServe()
}
