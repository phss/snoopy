package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	config := parseConfigOptions()

	var wg sync.WaitGroup
	wg.Add(len(config.ProxyConfigs))

	for _, proxyConfig := range config.ProxyConfigs {
		go func(pc ProxyConfig) {
			proxy(pc, config)
			wg.Done()
		}(proxyConfig)
	}

	wg.Wait()
}

func proxy(proxyConfig ProxyConfig, config Config) {
	fmt.Printf("Proxying for %s on http://localhost:%d \n", proxyConfig.Url, proxyConfig.Port)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", proxyConfig.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyRequest := NewProxyRequestFrom(r, proxyConfig.Url)
			printRequest(proxyRequest, config)
			response := makeProxyRequest(proxyRequest)
			printResponse(response, config)
			returnProxyResponse(response, w)
		}),
	}
	log.Fatal(server.ListenAndServe())
}

func returnProxyResponse(response ProxyResponse, w http.ResponseWriter) {
	for _, header := range response.Headers {
		w.Header().Add(header.Name, header.Value)
	}

	fmt.Fprintf(w, "%s", response.Body)
}

func makeProxyRequest(request ProxyRequest) ProxyResponse {
	client := http.Client{}
	resp, _ := client.Do(request.NewProxiedHttpRequest())
	return NewProxyResponseFrom(resp)
}
