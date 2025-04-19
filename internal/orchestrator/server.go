package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/sumeshmurali/mandarin/internal/config"
	prebuilttemplates "github.com/sumeshmurali/mandarin/internal/prebuilt_templates"
	"github.com/sumeshmurali/mandarin/internal/ratelimiter"
)

type Server struct {
	server   *http.Server
	starting chan struct{}
}

var ErrServerFailed = errors.New("server failed")

func NewServer() *Server {
	return &Server{
		starting: make(chan struct{}),
		server:   &http.Server{},
	}
}

func (s *Server) Run(config *config.Server) error {
	var globalRl ratelimiter.Ratelimiter
	if config.ServerConfig.RatelimitConfig != nil {
		globalRl = ratelimiter.NewRateLimiter(config.ServerConfig.RatelimitConfig)
	}
	var ratelimitMiddleWare = ratelimiter.RatelimitedHandlerMiddleWareCurry(globalRl)

	mux := http.NewServeMux()
	for name, endpoint := range config.Endpoints {
		rlMiddleWare := ratelimitMiddleWare
		if endpoint.RatelimitConfig != nil {
			rl := ratelimiter.NewRateLimiter(endpoint.RatelimitConfig)
			rlMiddleWare = ratelimiter.RatelimitedHandlerMiddleWareCurry(rl)
		}
		if endpoint.Template != "" {
			t, err := prebuilttemplates.GetTemplate(endpoint.Template)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, rlMiddleWare(t))
			continue
		}
		if endpoint.RequestConfig != nil && endpoint.ResponseConfig != nil {
			h, err := NewHandleFuncFromConfig(endpoint)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, rlMiddleWare(h))
		}

	}
	var port uint16
	if config.ServerConfig != nil && config.ServerConfig.Port != 0 {
		port = config.ServerConfig.Port
	} else {
		port = 80

	}
	s.server.Handler = mux
	s.server.Addr = fmt.Sprintf(":%d", port)

	fmt.Printf("Starting server %s on port %d\n", config.Name, port)
	close(s.starting)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed ", err)
		return errors.Join(err, ErrServerFailed)
	}
	return nil
}

func (s *Server) Shutdown() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
}

func (s *Server) WaitForStartup() {
	<-s.starting
}
