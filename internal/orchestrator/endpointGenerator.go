package orchestrator

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/sumeshmurali/mandarin/internal/config"
)

func NewHandleFuncFromConfig(endpoint config.Endpoint) (http.HandlerFunc, error) {

	if endpoint.RequestConfig == nil || endpoint.ResponseConfig == nil {
		return nil, fmt.Errorf("request or response config not provided for endpoint %s", endpoint.Name)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if len(endpoint.RequestConfig.AllowedMethods) != 0 && !slices.Contains(endpoint.RequestConfig.AllowedMethods, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		} else {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
		for k, v := range endpoint.ResponseConfig.Headers {
			w.Header().Set(k, v)
		}
		if endpoint.ResponseConfig.StatusCode != 0 {
			w.WriteHeader(endpoint.ResponseConfig.StatusCode)
		}
		fmt.Fprint(w, endpoint.ResponseConfig.Body)
	}, nil
}
