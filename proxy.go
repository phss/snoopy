package main

import (
	"fmt"
	"io/ioutil"
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
			resp, body := makeProxyRequest(proxyRequest)
			printResponse(resp, body, config.ShowBody)
			returnProxyResponse(resp, body, w)
		}),
	}
	log.Fatal(server.ListenAndServe())
}

func returnProxyResponse(resp *http.Response, body string, w http.ResponseWriter) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	defer resp.Body.Close()
	fmt.Fprintf(w, "%s", body)
}

func makeProxyRequest(request ProxyRequest) (*http.Response, string) {
	client := http.Client{}
	resp, _ := client.Do(request.NewProxiedHttpRequest())
	body, _ := ioutil.ReadAll(resp.Body)
	return resp, string(body)
}
