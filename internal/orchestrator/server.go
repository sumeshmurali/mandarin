package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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
	if config.ServerConfig != nil && config.ServerConfig.RatelimitConfig != nil {
		globalRl = ratelimiter.NewRateLimiter(config.ServerConfig.RatelimitConfig)
	}
	var ratelimitMiddleWare = ratelimiter.RatelimitedHandlerMiddleWareCurry(globalRl)
	var delay int
	if config.ServerConfig != nil {
		delay = config.ServerConfig.Delay
	}
	var delayMiddleWare = DelayMiddleWare(delay)
	mux := http.NewServeMux()
	for name, endpoint := range config.Endpoints {
		rlMiddleWare := ratelimitMiddleWare
		delayMW := delayMiddleWare

		if endpoint.EndpointConfig != nil {
			if endpoint.EndpointConfig.RatelimitConfig != nil {
				rl := ratelimiter.NewRateLimiter(endpoint.EndpointConfig.RatelimitConfig)
				rlMiddleWare = ratelimiter.RatelimitedHandlerMiddleWareCurry(rl)
			}
			if endpoint.EndpointConfig.Delay != 0 {
				delayMW = DelayMiddleWare(endpoint.EndpointConfig.Delay)
			}
		}

		if endpoint.Template != "" {
			t, err := prebuilttemplates.GetTemplate(endpoint.Template)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, rlMiddleWare(delayMW(t)))
			continue
		}
		if endpoint.RequestConfig != nil && endpoint.ResponseConfig != nil {
			h, err := NewHandleFuncFromConfig(endpoint)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, rlMiddleWare(delayMW(h)))
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

func DelayMiddleWare(delay int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if delay > 0 {
				// Simulate delay
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
			next(w, r)
		}
	}
}
