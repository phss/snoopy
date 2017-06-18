package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func printRequest(proxiedUrl string, r *http.Request, showBody bool) {
	fmt.Println("---")
	color.Blue("Request")
	color.Green("%s %s -> %s", r.Method, r.URL.Path, proxiedUrl)
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	if showBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
	}
	fmt.Println("")
}

func printResponse(resp *http.Response, body string, showBody bool) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(resp.Status))
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	if showBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		fmt.Println(body)
	}
}

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

type Config struct {
	showBody     bool
	proxyConfigs []ProxyConfig
}

type ProxyConfig struct {
	port int
	url  string
}

func parseOptions() Config {
	showBody := flag.Bool("showBody", false, "shows the request and response bodies")
	port := flag.Int("port", 8080, "proxy port")
	url := flag.String("url", "http://www.example.com", "url")
	flag.Parse()

	return Config{
		showBody: *showBody,
		proxyConfigs: []ProxyConfig{
			ProxyConfig{port: *port, url: *url},
		},
	}
}

func main() {
	config := parseOptions()
	var wg sync.WaitGroup
	wg.Add(len(config.proxyConfigs))

	for _, proxyConfig := range config.proxyConfigs {
		go func(pc ProxyConfig) {
			proxy(pc.port, pc.url, config.showBody)
			wg.Done()
		}(proxyConfig)
	}

	wg.Wait()
}
