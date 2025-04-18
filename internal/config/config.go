package config

import (
	"os"

	"github.com/spf13/viper"
)

type RequestConfig struct {
	AllowedMethods []string `mapstructure:"allowed_methods"`
}

type ResponseConfig struct {
	Headers map[string]string `mapstructure:"headers"`
	Body string `mapstructure:"body"`
	StatusCode int `mapstructure:"status_code"`
}

type EndpointConfig struct {
}
type Endpoint struct {
	Description string `mapstructure:"description"`
	Name string `mapstructure:"name"`
	Template string `mapstructure:"template"`
	EndpointConfig *EndpointConfig `mapstructure:"config"`
	RequestConfig *RequestConfig `mapstructure:"request_config"`
	ResponseConfig *ResponseConfig `mapstructure:"response_config"`
}

type ServerConfig struct {
	Port uint16 `mapstructure:"port"`
	Ratelimit int `mapstructure:"ratelimit"`
	RatelimitType string `mapstructure:"ratelimit_type"`
}

type Server struct {
	Name string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Endpoints map[string]Endpoint `mapstructure:"endpoints"`
	ServerConfig *ServerConfig `mapstructure:"config"`
}

func ParseConfiguration(configpath string) (*Server, error) {
	config := &Server{}
	file, err := os.Open(configpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(file)
	if err != nil {	
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}