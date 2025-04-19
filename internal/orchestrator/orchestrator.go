package orchestrator

import (
	"fmt"

	"github.com/sumeshmurali/mandarin/internal/config"
)

func Run(config *config.Server) {
	srv := NewServer()
	err := srv.Run(config)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
