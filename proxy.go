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
	client := http.Client{}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", proxyConfig.Port),
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, httpRequest *http.Request) {
			request := NewProxyRequestFrom(httpRequest, proxyConfig.Url)
			printHeader(proxyConfig.Name)
			printRequest(request, config)
			response := request.MakeRequest(client)
			printResponse(response, config)
			response.WriteResponse(writer)
		}),
	}
	log.Fatal(server.ListenAndServe())
}
