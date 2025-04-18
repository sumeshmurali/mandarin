package orchestrator

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/sumeshmurali/mandarin/internal/config"
	prebuilttemplates "github.com/sumeshmurali/mandarin/internal/prebuilt_templates"
)

type Server struct {
}

var ErrServerFailed = errors.New("server failed")


func (s *Server) Run(config *config.Server) error {
	fmt.Printf("%+v",config)
	mux := http.NewServeMux()
	for name, endpoint := range config.Endpoints {
		if endpoint.Template != "" {
			t, err := prebuilttemplates.GetTemplate(endpoint.Template)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, t)
			continue
		}
		if endpoint.RequestConfig != nil && endpoint.ResponseConfig != nil  {
			h, err := NewHandleFuncFromConfig(endpoint)
			if err != nil {
				return err
			}
			mux.HandleFunc(name, h)
		}
		
	}
	var port uint16
	if config.ServerConfig != nil && config.ServerConfig.Port != 0 {
		port = config.ServerConfig.Port
	} else {
		port = 80
		
	}
	server := http.Server{
		Handler: mux,
		Addr:   fmt.Sprintf(":%d", port),
	}
	fmt.Printf("Starting server %s on port %d\n", config.Name, port)
	err:= server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed ",err)
		return  errors.Join(err, ErrServerFailed)
	}
	return nil
}