package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func makeProxyRequest(proxiedUrl string, r *http.Request) (*http.Response, string) {
	client := http.Client{}
	proxyRequest, _ := http.NewRequest(r.Method, proxiedUrl, r.Body)
	resp, _ := client.Do(proxyRequest)
	body, _ := ioutil.ReadAll(resp.Body)
	return resp, string(body)
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

func proxy(port int, proxiedBaseUrl string, showBody bool) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxiedUrl := proxiedBaseUrl + r.URL.Path
		printRequest(proxiedUrl, r, showBody)
		resp, body := makeProxyRequest(proxiedUrl, r)
		printResponse(resp, body, showBody)
		returnProxyResponse(resp, body, w)
	})

	fmt.Printf("Proxying for %s on http://localhost:%d \n", proxiedBaseUrl, port)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	log.Fatal(server.ListenAndServe())
}

func main() {
	config := parseConfigOptions()
	var wg sync.WaitGroup
	wg.Add(len(config.ProxyConfigs))

	for _, proxyConfig := range config.ProxyConfigs {
		go func(pc ProxyConfig) {
			proxy(pc.Port, pc.Url, config.ShowBody)
			wg.Done()
		}(proxyConfig)
	}

	wg.Wait()
}
