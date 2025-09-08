package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Endpoint struct {
	Method string
	URL    string
	Path   string
}

func ParseOpenAPISpec(specPath string, baseURLOverride string) ([]Endpoint, string, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, "", err
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, "", err
	}

	baseURL := baseURLOverride
	if baseURL == "" {
		if servers, ok := spec["servers"].([]interface{}); ok && len(servers) > 0 {
			if srv, ok := servers[0].(map[string]interface{}); ok {
				if url, ok := srv["url"].(string); ok && url != "" {
					baseURL = url
				}
			}
		}
	}
	if baseURL == "" {
		return nil, "", fmt.Errorf("no base URL found in spec and none provided via --base-url")
	}

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("no paths found in spec")
	}

	var endpoints []Endpoint
	for path, rawMethods := range paths {
		methods, ok := rawMethods.(map[string]interface{})
		if !ok {
			continue
		}
		for method := range methods {
			url := strings.TrimRight(baseURL, "/") + path
			endpoints = append(endpoints, Endpoint{
				Method: strings.ToUpper(method),
				URL:    url,
				Path:   path,
			})
		}
	}

	return endpoints, baseURL, nil
}
