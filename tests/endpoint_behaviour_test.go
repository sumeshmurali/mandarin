package tests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sumeshmurali/mandarin/internal/config"
	"github.com/sumeshmurali/mandarin/internal/orchestrator"
)

func StartNewMockServer(c *config.Server) *orchestrator.Server {
	// Create a new server instance
	srv := orchestrator.NewServer()
	// Define the configuration for the server
	if c == nil {
		c = &config.Server{
			Name: "test-server",
			ServerConfig: &config.ServerConfig{
				Port: TestServerPort,
			},
			Endpoints: map[string]config.Endpoint{
				"/": {
					RequestConfig: &config.RequestConfig{
						AllowedMethods: []string{"GET"},
					},
					ResponseConfig: &config.ResponseConfig{
						Body: "Hello, World!",
					},
				},
			},
		}
	}

	// Run the server with the configuration
	go func() {
		err := srv.Run(c)
		if err != nil {
			log.Printf("Failed to run server: %v\n", err)
		}
	}()
	srv.WaitForStartup()
	// Add more tests as needed
	log.Println("Server is running successfully")
	return srv
}

func TestHello(t *testing.T) {
	srv := StartNewMockServer(nil)
	defer srv.Shutdown()

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", TestServerPort), nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response body:", err)
	}
	if string(body) != "Hello, World!" {
		t.Fatalf("Expected body 'Hello, World!', got %s", body)
	}
}

func TestRatelimitGlobal(t *testing.T) {

	ratelimit := 5
	cases := []struct {
		name   string
		config *config.Server
	}{
		{
			name: "Server Level Ratelimit",
			config: &config.Server{
				Name: "test-server",
				ServerConfig: &config.ServerConfig{
					Port:            TestServerPort,
					RatelimitConfig: &config.RatelimitConfig{Ratelimit: ratelimit, RatelimitType: "global"},
				},

				Endpoints: map[string]config.Endpoint{
					"/": {
						Template: "echo",
					},
				},
			},
		},
		{
			name: "Endpoint Level Ratelimit",
			config: &config.Server{
				Name: "test-server",
				ServerConfig: &config.ServerConfig{
					Port: TestServerPort,
				},
				Endpoints: map[string]config.Endpoint{
					"/": {
						Template: "echo",
						EndpointConfig: &config.EndpointConfig{
							RatelimitConfig: &config.RatelimitConfig{Ratelimit: ratelimit, RatelimitType: "global"},
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := StartNewMockServer(tc.config)
			defer srv.Shutdown()
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", TestServerPort), nil)
			if err != nil {
				t.Fatal("Failed to create request:", err)
			}

			// send ratelimit +1 requests to the endpoint and see if ratelimit is working
			exceededCalls := 0
			var wg sync.WaitGroup
			for i := 0; i < ratelimit+1; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					resp, err := client.Do(req)
					if err != nil {
						log.Fatal("Failed to send request:", err)
					}
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						if resp.StatusCode == http.StatusTooManyRequests {
							exceededCalls++
							return
						}
						log.Fatalf("Expected status OK, got %v", resp.Status)
					}
				}()
			}
			wg.Wait()
			if exceededCalls != 1 {
				t.Fatalf("Expected 1 exceeded call, got %d", exceededCalls)
			}
			// see if following calls are successful after ratelimit is reset
			time.Sleep(time.Second * 1)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal("Failed to send request:", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected status OK, got %v", resp.Status)
			}
		})
	}

}
